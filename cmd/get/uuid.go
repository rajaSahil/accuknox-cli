package get

import (
	"fmt"

	"github.com/accuknox/accuknox-cli/k8s"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var uuidCmd = &cobra.Command{
	Use:   "platform-uuid",
	Short: "get platform uuid.",
	Long:  `get platform uuid to generate license for discovery-engine`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := k8s.ConnectK8sClient()
		if err != nil {
			fmt.Printf("unable to create Kubernetes clients: %s\n", err.Error())
			return err
		}

		kubeSystem, err := client.K8sClientset.CoreV1().Namespaces().Get(context.TODO(), "kube-system", metav1.GetOptions{})
		if err != nil {
			fmt.Printf("unable to get platform-uuid: %s\n", err.Error())
			return err
		}
		// convert the namespace UID to a string and print it
		fmt.Printf("Platform UUID: %s\n", string(kubeSystem.UID))
		return nil
	},
}

func init() {
	GetCmd.AddCommand(uuidCmd)

}
