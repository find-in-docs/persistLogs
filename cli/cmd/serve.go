/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"time"

	"github.com/samirgadkari/persist/pkg/config"
	"github.com/samirgadkari/sidecar/pkg/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the persistence service.",
	Long: `Start the persistence service. This service will get all messages from
the message queue and write them into a database.`,
	Run: func(cmd *cobra.Command, args []string) {

		config.LoadConfig()

		sidecarServiceAddr := viper.GetString("sidecarServiceAddr")
		_, sidecar, err := client.Connect(sidecarServiceAddr)
		if err != nil {
			return
		}

		logMsgTest := "Persist sending test log message."
		err = sidecar.Log(&logMsgTest)
		if err != nil {
			return
		}

		err = sidecar.Sub("search.v1.*")
		if err != nil {
			return
		}

		go func() {
			for {
				subTopicRsp, err := sidecar.Recv()
				if err != nil {
					fmt.Printf("Error receiving from sidecar: %#v\n", err)
					break
				}

				fmt.Printf("Received from sidecar: \n\t%#v\n", subTopicRsp)
			}
		}()

		pubMsgTest := []byte("Persist sending test pub message.")
		err = sidecar.Pub("search.v1.test", pubMsgTest)
		if err != nil {
			return
		}

		time.Sleep(3 * time.Second)
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
