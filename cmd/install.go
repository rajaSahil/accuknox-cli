// SPDX-License-Identifier: Apache-2.0
// Copyright 2021 Authors of KubeArmor

package cmd

import (
	"fmt"
	"os"
	"os/exec"

	de "github.com/accuknox/accuknox-cli/discoveryengine"
	"github.com/accuknox/accuknox-cli/install"
	"github.com/spf13/cobra"
)

var (
	installOptions install.Options
	key            string
	user           string
)

// installCmd represents the get command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install KubeArmor in a Kubernetes Cluster",
	Long:  `Install KubeArmor in a Kubernetes Clusters`,
	RunE: func(cmd *cobra.Command, args []string) error {
		//kubearmor
		installOptions.Animation = true
		if err := installOptions.Env.CheckAndSetValidEnvironmentOption(cmd.Flag("env").Value.String()); err != nil {
			return fmt.Errorf("error in checking environment option: %v", err)
		}
		if err := install.K8sInstaller(client, installOptions); err != nil {
			return err
		}

		//Discovery Engine
		_, err := exec.LookPath("kubectl")
		if err != nil {
			fmt.Println("kubectl is not installed. Follow link to install kubectl : https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/")
			return err
		}

		command := "kubectl"
		arg := []string{"apply", "-f", "https://raw.githubusercontent.com/accuknox/discovery-engine/feature-report/deployments/k8s/deployment.yaml"}

		newCmd := exec.Command(command, arg...)

		newCmd.Stdout = os.Stdout
		newCmd.Stderr = os.Stderr

		err = newCmd.Run()
		if err != nil {
			fmt.Printf("Failed to execute command: %v\n", err)
			return err
		}
		fmt.Println("🥳  Done Installing Discovery Engine")

		de.CheckPods(client)

		if user == "" || key == "" {
			return nil
		}

		err = de.InstallLicense(client, key, user)
		return err
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	//license
	installCmd.Flags().StringVar(&key, "key", "", "license key for installing license (required)")
	installCmd.Flags().StringVar(&user, "user", "", "user id for installing license")

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
