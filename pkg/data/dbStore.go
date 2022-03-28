package data

import (
	"context"
	"fmt"
	"os"
	"regexp"

	"github.com/jackc/pgx/v4"
	pb "github.com/samirgadkari/sidecar/protos/v1/messages"
	"github.com/spf13/viper"
)

type DB struct {
	conn          *pgx.Conn
	tableName     string
	persistSchema string
	msgStrRegex   *regexp.Regexp
}

type WordInt uint64
type DocumentId WordInt

type Doc struct {
	Header pb.Header
	Bytes  []byte
}

func DBConnect() (*DB, error) {

	conn, err := pgx.Connect(context.Background(), viper.GetString("output.connection"))
	if err != nil {
		fmt.Printf("Error connecting to postgres database using: %s\n",
			viper.GetString("output.location"))
		fmt.Printf("err: %v\n", err)
		return nil, err
	}

	persistSchema := `(msgType integer,
			srcServType varchar(16),
			dstServType varchar(16),
			servId varchar(36),
			msgId integer,
			msg text)`

	db := DB{
		conn:          conn,
		persistSchema: persistSchema,
		msgStrRegex:   regexp.MustCompile(`\\+?`),
	}

	return &db, nil
}

func (db *DB) createTable(tableName string) {

	checkIfExists := `select 'public.` + tableName + `'::regclass;`
	if _, err := db.conn.Exec(context.Background(), checkIfExists); err != nil {
		fmt.Printf("Table %s does not exist, so create it.\n", tableName)

		createString := `create table ` + tableName + ` ` + db.persistSchema + `;`
		if _, err := db.conn.Exec(context.Background(), createString); err != nil {
			fmt.Printf("Failed to create the persistSchema. err: %v\n", err)
			os.Exit(-1)
		}
	}

	db.tableName = tableName
}

func (db *DB) CreateTable(tableName string) error {

	if db.conn == nil {
		fmt.Printf("Create db connection before creating schema\n")
		return fmt.Errorf("Create db connection before creating schema\n")
	}

	checkIfExists := `select 'public.` + tableName + `'::regclass;`
	if _, err := db.conn.Exec(context.Background(), checkIfExists); err != nil {
		fmt.Printf("Table %s does not exist, so create it.\n", tableName)

		createString := `create table ` + tableName + ` ` + db.persistSchema + `;`
		if _, err := db.conn.Exec(context.Background(), createString); err != nil {
			fmt.Printf("Failed to create the db.persistSchema. err: %v\n", err)
			return err
		}
	}

	return nil
}

func (db *DB) formatMsg(msg *string) *string {

	result := db.msgStrRegex.ReplaceAllString(*msg, "")
	return &result
}

func (db *DB) StoreData(header *pb.Header, msg *string, tableName string) error {

	msg2 := db.formatMsg(msg)
	fmt.Printf("msg2: %s\n\n\n", msg2)

	createDocString := `(msgType, srcServType, dstServType, servId, msgId, msg)`
	insertStatement := `insert into ` + tableName + ` ` + createDocString +
		` values ($1, $2, $3, $4, $5, $6);`

	if _, err := db.conn.Exec(context.Background(), insertStatement,
		header.MsgType, header.SrcServType, header.DstServType,
		header.ServId, header.MsgId, msg2); err != nil {
		fmt.Printf("Store data failed. err: %v\n", err)
		return err
	}

	return nil
}

func (db *DB) ReadData() *Doc {

	return nil
}

func (db *DB) DBDisconnect() error {

	if db.conn == nil {
		fmt.Printf("conn is nil\n")
		os.Exit(-1)
	}

	err := db.conn.Close(context.Background())
	if err != nil {
		fmt.Printf("Error closing DB connection: %v\n", err)
		return err
	}

	return nil
}
