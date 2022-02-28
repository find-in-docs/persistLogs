package cmd

import (
	"fmt"

	"github.com/samirgadkari/sidecar/pkg/client"
	"github.com/samirgadkari/sidecar/protos/v1/messages"
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

		sidecar := client.InitSidecar("persist")

		topic := "search.v1.*"
		if err := sidecar.ProcessSubMsgs(topic, allTopicsRecvChanSize, func(m *messages.SubTopicResponse) {

			fmt.Printf("Received from sidecar: \n\t%v\n", m)
		}); err != nil {
			fmt.Printf("Error processing subscription messages:\n\ttopic: %s\n\terr: %v\n",
				topic, err)
		}
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
