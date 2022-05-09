package cmd

import (
	"github.com/bpalermo/memento/internal/logger"
	"github.com/bpalermo/memento/pkg/register"
	"github.com/spf13/cobra"
	clientV3 "go.etcd.io/etcd/client/v3"
	"time"
)

const (
	basePathKey    = "base-path"
	serviceNameKey = "service-name"
	servicePortKey = "service-port"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Service registration",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: registerToEtcd,
}

func init() {
	rootCmd.AddCommand(registerCmd)

	registerCmd.Flags().StringP(basePathKey, "b", "/discovery", "Etcd discovery base path")
	registerCmd.Flags().StringP(serviceNameKey, "s", "", "Service name")
	registerCmd.Flags().Uint16P(servicePortKey, "p", 18080, "Service port")

	_ = registerCmd.MarkFlagRequired(serviceNameKey)
}

func registerToEtcd(cmd *cobra.Command, _ []string) {
	log := logger.New()

	basePath, _ := cmd.Flags().GetString(basePathKey)
	serviceName, _ := cmd.Flags().GetString(serviceNameKey)
	servicePort, _ := cmd.Flags().GetUint16(servicePortKey)

	ser, err := register.NewEtcdRegister(clientV3.Config{
		Endpoints:   []string{""},
		DialTimeout: 5 * time.Second,
	}, log, basePath, serviceName, servicePort, 5, 10)
	if err != nil {
		log.WithError(err).Fatal("could start register")
	}

	// listen to keep alive chan
	go ser.Listen()

	select {}
}
