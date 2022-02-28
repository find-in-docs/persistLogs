package conn

import (
	"fmt"
	"os"

	"github.com/samirgadkari/persist/pkg/config"
	"github.com/samirgadkari/sidecar/pkg/client"
	"github.com/spf13/viper"
)

func InitSidecar(serviceName string) *client.SC {

	config.LoadConfig()

	sidecarServiceAddr := viper.GetString("sidecarServiceAddr")
	_, sidecar, err := client.Connect(serviceName, sidecarServiceAddr)
	if err != nil {
		fmt.Printf("Error connecting to client: %v\n", err)
		os.Exit(-1)
	}

	return sidecar
}
