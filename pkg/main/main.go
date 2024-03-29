package main

import (
  "sync"
	"context"
	"fmt"
	"os"
	"regexp"
  // "time"

	"github.com/find-in-docs/persistLogs/pkg/config"
	"github.com/find-in-docs/persistLogs/pkg/data"
	"github.com/find-in-docs/sidecar/pkg/client"
	// "github.com/find-in-docs/sidecar/pkg/utils"
	pb "github.com/find-in-docs/sidecar/protos/v1/messages"
	"github.com/spf13/viper"
)

const (
	allTopicsRecvChanSize = 32
)

func formatMsg(msg *string, re *regexp.Regexp) *string {

	result := re.ReplaceAllString(*msg, "")
	return &result
}

func main() {

  var wg sync.WaitGroup
  wg.Add(1)

  fmt.Printf("Loading configuration\n")
	config.Load()

  fmt.Printf("Connecting to DB\n")
	// Setup database
	db, err := data.DBConnect()
	if err != nil {
		return
	}

  fmt.Printf("Creating table\n")
	tableName := "logs"
	err = db.CreateTable(tableName)
	if err != nil {
		return
	}

  fmt.Printf("Initializing sidecar\n")
  fmt.Printf("serviceName: %s\n", viper.GetString("serviceName"))
	sidecar, err := client.InitSidecar(viper.GetString("serviceName"), nil)
	if err != nil {
		fmt.Printf("Error initializing sidecar: %v\n", err)
		os.Exit(-1)
	}

	topic := "search.log.v1"

	msgStrRegex := regexp.MustCompile(`\\+?\"|\\+?n|\\+?t`)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

  fmt.Printf("Setting up processing of messages\n")
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

	/*
		sidecar.Logger.Log("Persist sending log message test: %s\n", "search.log.v1")
		time.Sleep(3 * time.Second)

		var retryNum uint32 = 1
		retryDelayDuration, err := time.ParseDuration("200ms")
		if err != nil {
			fmt.Printf("Error creating Golang time duration.\nerr: %v\n", err)
			os.Exit(-1)
		}
		retryDelay := durationpb.New(retryDelayDuration)

		err = sidecar.Pub(ctx, "search.data.v1", []byte("test pub message"),
			&pb.RetryBehavior{
				RetryNum:   &retryNum,
				RetryDelay: retryDelay,
			},
		)
		if err != nil {
			fmt.Printf("Error publishing message.\n\terr: %v\n", err)
		}
	*/

	/* Message with no retry params
	err = sidecar.Pub(ctx, "search.data.v1", []byte("test pub message"), nil)
	*/

  /*
  This section was for testing if the GoRoutines are all
  finished before exiting. This was only meant as a debug
  mechanism. Since we're now running in minikube, this mechanism
  will not work, so it is commented out for now.
  Not sure what to replace it with, at this time.

	fmt.Println("Press the Enter key to stop")
	fmt.Scanln()
	fmt.Println("User pressed Enter key")

	// Signal that we want the process subscription goroutines to end.
	// This cancellation causes the goroutines to unsubscribe from the topic
	// before they end themselves.
	cancel()

	sleepDur, _ := time.ParseDuration("3s")
	fmt.Printf("Sleeping for %s seconds\n", sleepDur)
	time.Sleep(sleepDur)

	utils.ListGoroutinesRunning()
  */
	wg.Wait()   // wait forever
}
