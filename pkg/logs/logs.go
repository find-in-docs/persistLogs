package logs

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/samirgadkari/persist/pkg/config"
	"github.com/samirgadkari/persist/pkg/data"
	"github.com/samirgadkari/sidecar/pkg/client"
	pb "github.com/samirgadkari/sidecar/protos/v1/messages"
	"google.golang.org/protobuf/types/known/durationpb"
)

const (
	allTopicsRecvChanSize = 32
)

func formatMsg(msg *string, re *regexp.Regexp) *string {

	result := re.ReplaceAllString(*msg, "")
	return &result
}

func PersistLogs() {

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

	sidecar := client.InitSidecar(tableName, nil)
	if sidecar == nil {
		fmt.Printf("Error initializing sidecar - Exiting\n")
		os.Exit(-1)
	}

	topic := "search.*.v1"

	msgStrRegex := regexp.MustCompile(`\\+?\"|\\+?n|\\+?t`)

	ctx, cancel := context.WithCancel(context.Background())

	err = sidecar.ProcessSubMsgs(ctx, topic,
		allTopicsRecvChanSize, func(m *pb.SubTopicResponse) {

			msg := fmt.Sprintf("Received from sidecar:\n\t%s", m.String())
			msg2 := formatMsg(&msg, msgStrRegex)
			fmt.Printf("%s\n", *msg2)

			db.StoreData(m.Header, msg2, tableName)
		})
	if err != nil {
		fmt.Printf("Error processing subscription messages:\n\ttopic: %s\n\terr: %v\n",
			topic, err)
	}

	sidecar.Logger.Log("Persist sending log message test: %s\n", "search.log.v1")
	time.Sleep(3 * time.Second)

	var retryNum uint32 = 4
	retryDelayDuration, err := time.ParseDuration("2s")
	if err != nil {
		fmt.Printf("Error creating Golang time duration.\nerr: %v\n", err)
		os.Exit(-1)
	}
	retryDelay := durationpb.New(retryDelayDuration)

	err = sidecar.Pub("search.data.v1", []byte("test pub message"),
		&pb.RetryBehavior{
			RetryNum:   &retryNum,
			RetryDelay: retryDelay,
		},
	)
	if err != nil {
		fmt.Printf("Error publishing message.\n\terr: %v\n", err)
		os.Exit(-1)
	}

	sidecar.Unsub(topic)
	cancel()  // Signal that we want the process subscription goroutines to end
	select {} // This will wait forever
}
