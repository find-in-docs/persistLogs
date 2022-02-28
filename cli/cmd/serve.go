package cmd

import (
	"fmt"

	"github.com/samirgadkari/persist/pkg/conn"
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

		sidecar := conn.InitSidecar("persist")

		err := sidecar.Sub("search.v1.*", allTopicsRecvChanSize)
		if err != nil {
			return
		}

		for {
			subTopicRsp, err := sidecar.Recv("search.v1.*")
			if err != nil {
				sidecar.Log("Error receiving from sidecar: %#v\n", err)
				break
			}

			fmt.Printf("Received from sidecar: \n\t%#v\n", subTopicRsp)
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
