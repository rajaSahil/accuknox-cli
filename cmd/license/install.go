package license

import (
	"context"
	"fmt"
	"github.com/accuknox/accuknox-cli/utils"
	"os"
	"strconv"

	"github.com/accuknox/accuknox-cli/k8s"
	pb "github.com/accuknox/auto-policy-discovery/src/protobuf/v1/license"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	key  string
	user string
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install License",
	Long:  `Install license for discovery engine`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := k8s.ConnectK8sClient()
		if err != nil {
			fmt.Printf("unable to create Kubernetes clients: %s\n", err.Error())
			return err
		}

		gRPC := ""
		targetSvc := "discovery-engine"

		if val, ok := os.LookupEnv("DISCOVERY_SERVICE"); ok {
			gRPC = val
		} else {
			pf, err := utils.InitiatePortForward(client, port, port, matchLabels, targetSvc)
			if err != nil {
				return err
			}
			gRPC = "localhost:" + strconv.FormatInt(pf.LocalPort, 10)
		}

		conn, err := grpc.Dial(gRPC, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return err
		}
		defer conn.Close()

		licenseClient := pb.NewLicenseClient(conn)

		req := &pb.LicenseInstallRequest{
			Key:    key,
			UserId: user,
		}
		_, err = licenseClient.InstallLicense(context.Background(), req)
		if err != nil {
			return err
		}
		fmt.Printf("License installed successfully for discovery engine.\n")

		return nil
	},
}

func init() {
	LicenseCmd.AddCommand(installCmd)

	installCmd.Flags().StringVar(&key, "key", "", "license key for installing license (required)")
	installCmd.Flags().StringVar(&user, "user", "", "user id for installing license")

	err := installCmd.MarkFlagRequired("key")
	if err != nil {
		fmt.Printf("key flag is required : %s\n", err)
	}
	err = installCmd.MarkFlagRequired("user")
	if err != nil {
		fmt.Printf("user flag is required : %s\n", err)
	}
}
