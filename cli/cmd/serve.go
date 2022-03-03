package cmd

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/samirgadkari/persist/pkg/config"
	"github.com/samirgadkari/persist/pkg/data"
	"github.com/samirgadkari/sidecar/pkg/client"
	scconn "github.com/samirgadkari/sidecar/pkg/conn"
	pb "github.com/samirgadkari/sidecar/protos/v1/messages"
	"github.com/spf13/cobra"
)

const (
	allTopicsRecvChanSize = 32
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the persistence service.",
	Long: `Start the persistence service. This service will get all messages from
the message queue and write them into a database.`,
	Run: func(cmd *cobra.Command, args []string) {

		config.Load()

		// Setup database
		db, err := data.DBConnect()
		if err != nil {
			return
		}

		tableName := "persist"
		err = db.CreateTable(tableName)
		if err != nil {
			return
		}

		sidecar := client.InitSidecar(tableName)

		topic := "search.v1.*"
		buf := &bytes.Buffer{}
		enc := gob.NewEncoder(buf)

		go func() {
			if err = sidecar.ProcessSubMsgs(topic, allTopicsRecvChanSize, func(m *pb.SubTopicResponse) {

				scconn.PrintSubTopicRsp("Received from sidecar:", m)

				err = enc.Encode(*m)
				if err != nil {
					fmt.Printf("Error converting message to Big Endian: \n\terror: %v\n", err)
				}

				db.StoreData(m.Header, buf, tableName)
			}); err != nil {
				fmt.Printf("Error processing subscription messages:\n\ttopic: %s\n\terr: %v\n",
					topic, err)
			}
		}()

		sidecar.Log("Persist sending log message test: %s", "search.v1.log")
		time.Sleep(3 * time.Second)
		sidecar.Unsub(topic)
		select {} // This will wait forever
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
