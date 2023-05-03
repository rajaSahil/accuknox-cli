// SPDX-License-Identifier: Apache-2.0
// Copyright 2021 Authors of KubeArmor

package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/accuknox/accuknox-cli/install"
	"github.com/accuknox/accuknox-cli/utils"
	pb "github.com/accuknox/auto-policy-discovery/src/protobuf/v1/license"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	installOptions install.Options
	key            string
	user           string
)
var matchLabels = map[string]string{"app": "discovery-engine"}
var port int64 = 9089

// installCmd represents the get command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install KubeArmor in a Kubernetes Cluster",
	Long:  `Install KubeArmor in a Kubernetes Clusters`,
	RunE: func(cmd *cobra.Command, args []string) error {
		installOptions.Animation = true
		if err := installOptions.Env.CheckAndSetValidEnvironmentOption(cmd.Flag("env").Value.String()); err != nil {
			return fmt.Errorf("error in checking environment option: %v", err)
		}
		if err := install.K8sInstaller(client, installOptions); err != nil {
			return err
		}
		return nil
	},
}
var licenseCmd = &cobra.Command{
	Use:   "license",
	Short: "install license",
	Long:  `install license with flags key and user`,
	RunE: func(cmd *cobra.Command, args []string) error {
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

		client := pb.NewLicenseClient(conn)

		req := &pb.LicenseInstallRequest{
			Key:    key,
			UserId: user,
		}
		_, err = client.InstallLicense(context.Background(), req)
		if err != nil {
			return err
		}
		fmt.Printf("License installed successfully for discovery engine.\n")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.AddCommand(licenseCmd)

	//license
	licenseCmd.Flags().StringVar(&key, "key", "", "license key for installing license (required)")
	licenseCmd.Flags().StringVar(&user, "user", "", "user id for installing license")
	err := licenseCmd.MarkFlagRequired("key")
	if err != nil {
		fmt.Printf("Required flag empty : %s", err)
	}

	installCmd.Flags().StringVarP(&installOptions.Namespace, "namespace", "n", "kube-system", "Namespace for resources")
	installCmd.Flags().StringVarP(&installOptions.KubearmorImage, "image", "i", "kubearmor/kubearmor:stable", "Kubearmor daemonset image to use")
	installCmd.Flags().StringVarP(&installOptions.InitImage, "init-image", "", "kubearmor/kubearmor-init:stable", "Kubearmor daemonset init container image to use")
	installCmd.Flags().StringVarP(&installOptions.Tag, "tag", "t", "", "Change image tag/version for default kubearmor images (This will overwrite the tags provided in --image/--init-image)")
	installCmd.Flags().StringVarP(&installOptions.Audit, "audit", "a", "", "Kubearmor Audit Posture Context [all,file,network,capabilities]")
	installCmd.Flags().StringVarP(&installOptions.Block, "block", "b", "", "Kubearmor Block Posture Context [all,file,network,capabilities]")
	installCmd.Flags().BoolVar(&installOptions.Save, "save", false, "Save KubeArmor Manifest ")
	installCmd.Flags().BoolVar(&installOptions.Local, "local", false, "Use Local KubeArmor Images (sets ImagePullPolicy to 'IfNotPresent') ")
	installCmd.Flags().StringVarP(&installOptions.Env.Environment, "env", "e", "", "Supported KubeArmor Environment [k3s,microK8s,minikube,gke,bottlerocket,eks,docker,oke,generic]")

}
