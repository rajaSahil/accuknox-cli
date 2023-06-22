package cmd

import (
	"github.com/accuknox/accuknox-cli/report"
	"github.com/spf13/cobra"
)

var reportOptions report.Options

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Report from discovery engine",
	Long:  `Discovery engine keeps the telemetry information from the policy enforcement engines and the karmor connects to it to provide this as observability data`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := reportOptions.Report(client); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)

	reportCmd.Flags().StringVar(&reportOptions.GRPC, "gRPC", "", "gRPC server information")
	reportCmd.Flags().StringArrayVarP(&reportOptions.Clusters, "clusters", "c", []string{""}, "Clusters names")
	reportCmd.Flags().StringArrayVarP(&reportOptions.Namespaces, "namespaces", "n", []string{""}, "Namespaces names")
	reportCmd.Flags().StringArrayVarP(&reportOptions.ResourceType, "resource-types", "r", []string{""}, "Resource types")
	reportCmd.Flags().StringArrayVarP(&reportOptions.ResourceName, "resource-names", "w", []string{""}, "Resource names")
	reportCmd.Flags().StringVarP(&reportOptions.Labels, "labels", "l", "", "Labels")
	reportCmd.Flags().StringVar(&reportOptions.ContainerName, "container", "", "Container name")
	reportCmd.Flags().StringArrayVarP(&reportOptions.Source, "source", "s", []string{""}, "Source path")
	reportCmd.Flags().StringArrayVarP(&reportOptions.Destination, "destination", "d", []string{""}, "Destination path")
	reportCmd.Flags().StringVar(&reportOptions.Operation, "o", "", "Operation type")

	reportCmd.Flags().StringVarP(&reportOptions.BaselineJsonPath, "baseline-json-filepath", "b", "baseline/report.json", "Operation type")
}
