package license

import (
	"context"
	"encoding/json"
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

var licenseStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "get license status",
	Long:  `get license status`,

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
				fmt.Printf("unable to initiate port forward: %s\n", err.Error())
				return err
			}
			gRPC = "localhost:" + strconv.FormatInt(pf.LocalPort, 10)
		}

		conn, err := grpc.Dial(gRPC, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Printf("unable to dial to the target grpc: %s\n", err.Error())
			return err
		}
		defer conn.Close()

		licenseClient := pb.NewLicenseClient(conn)
		req := &pb.LicenseStatusRequest{}
		resp, err := licenseClient.GetLicenseStatus(context.Background(), req)
		if err != nil {
			fmt.Printf("unable to get license status: %s\n", err.Error())
			return err
		}

		jsonResp, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return err
		}
		fmt.Println(string(jsonResp))

		return nil
	},
}

func init() {
	LicenseCmd.AddCommand(licenseStatusCmd)
}
