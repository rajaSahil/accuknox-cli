package report

import (
	"fmt"
	rpb "github.com/accuknox/accuknox-cli/rpb"
	"github.com/accuknox/accuknox-cli/summary"
	opb "github.com/accuknox/auto-policy-discovery/src/protobuf/v1/observability"
	"github.com/clarketm/json"
	"github.com/lensesio/tableprinter"
	"io/ioutil"
	"os"
	"strings"
)

func getTableOutput(diffReport *rpb.ReportResponse) {

	for ck, cv := range diffReport.Clusters {
		for nk, nv := range cv.Namespaces {
			for _, rtv := range nv.ResourceTypes {
				for rsdk, rsdv := range rtv.Resources {

					resp := &opb.Response{
						DeploymentName:    rsdk,
						PodName:           "",
						ClusterName:       ck,
						Namespace:         nk,
						Label:             rsdv.GetMetaData().Label,
						ContainerName:     rsdv.GetMetaData().ContainerName,
						ProcessData:       rsdv.GetSummaryData().GetProcessData(),
						FileData:          rsdv.GetSummaryData().GetFileData(),
						IngressConnection: rsdv.GetSummaryData().GetIngressConnection(),
						EgressConnection:  rsdv.GetSummaryData().GetEgressConnection(),
						BindConnection:    rsdv.GetSummaryData().GetBindConnection(),
					}
					if len(rsdv.GetSummaryData().GetProcessData()) > 0 {
						summary.DisplaySummaryOutput(resp, false, "process")
						tableprinter.Print(os.Stdout, rsdv.GetSummaryData().GetProcessData())
					}
					if len(rsdv.GetSummaryData().GetFileData()) > 0 {
						summary.DisplaySummaryOutput(resp, false, "file")
						tableprinter.Print(os.Stdout, rsdv.GetSummaryData().GetFileData())
					}
					if len(rsdv.GetSummaryData().GetIngressConnection()) > 0 || len(rsdv.GetSummaryData().GetEgressConnection()) > 0 || len(rsdv.GetSummaryData().GetBindConnection()) > 0 {
						summary.DisplaySummaryOutput(resp, false, "network")
						tableprinter.Print(os.Stdout, rsdv.GetSummaryData().GetEgressConnection())
					}

				}
			}
		}
	}

}

func getDiff(baselineReport, report *rpb.ReportResponse, ignorePath []string) error {
	diffReport := &rpb.ReportResponse{}
	diffReport.Clusters = map[string]*rpb.ClusterData{}

	for ck, cv := range report.GetClusters() {
		if _, ok := baselineReport.Clusters[ck]; !ok {
			diffReport.Clusters[ck] = cv
			continue
		}

		// TODO: can be escaped, only create object if all level check passed i.e until resource_name
		diffReport.Clusters[ck] = &rpb.ClusterData{
			ClusterName: ck,
			Namespaces:  map[string]*rpb.NamespaceData{},
		}

		for nk, nv := range cv.GetNamespaces() {
			if _, ok := baselineReport.Clusters[ck].Namespaces[nk]; !ok {
				diffReport.Clusters[ck].Namespaces[nk] = nv
				continue
			}

			// TODO: can be escaped, only create object if all level check passed i.e until resource_name
			diffReport.Clusters[ck].Namespaces[nk] = &rpb.NamespaceData{
				NamespaceName: nk,
				ResourceTypes: map[string]*rpb.ResourceTypeData{},
			}

			for rtk, rtv := range nv.GetResourceTypes() {
				if _, ok := baselineReport.Clusters[ck].Namespaces[nk].ResourceTypes[rtk]; !ok {
					diffReport.Clusters[cv.ClusterName].Namespaces[nv.NamespaceName].ResourceTypes[rtk] = rtv
					continue
				}

				// TODO: can be escaped, only create object if all level check passed i.e until resource_name
				diffReport.Clusters[cv.ClusterName].Namespaces[nv.NamespaceName].ResourceTypes[rtk] = &rpb.ResourceTypeData{
					ResourceType: rtk,
					Resources:    map[string]*rpb.ResourceData{},
				}

				for rsdk, rsdv := range rtv.GetResources() {
					if _, ok := baselineReport.Clusters[ck].Namespaces[nk].ResourceTypes[rtk].Resources[rsdk]; !ok {
						diffReport.Clusters[cv.ClusterName].Namespaces[nv.NamespaceName].ResourceTypes[rtk].Resources[rsdk] = rsdv
						continue
					}

					// TODO: can be escaped, only create object if all level check passed i.e until resource_name
					diffReport.Clusters[cv.ClusterName].Namespaces[nv.NamespaceName].ResourceTypes[rtk].Resources[rsdk] = &rpb.ResourceData{
						ResourceType: rsdv.ResourceType,
						ResourceName: rsdv.ResourceName,
						MetaData: &rpb.MetaData{
							Label:         rsdv.GetMetaData().Label,
							ContainerName: rsdv.GetMetaData().ContainerName,
						},
						SummaryData: &rpb.SummaryData{
							ProcessData:       []*opb.SysProcFileSummaryData{},
							FileData:          []*opb.SysProcFileSummaryData{},
							IngressConnection: []*opb.SysNwSummaryData{},
							EgressConnection:  []*opb.SysNwSummaryData{},
							BindConnection:    []*opb.SysNwSummaryData{},
						},
					}

					for _, rsd := range rsdv.GetSummaryData().GetProcessData() {

						if ignoreComparison(rsd.Source, rsd.Destination, ignorePath) {
							continue
						}
						addInDiff := true

						for _, bsd := range baselineReport.Clusters[ck].Namespaces[nk].ResourceTypes[rtk].Resources[rsdk].GetSummaryData().GetProcessData() {
							d := compareProcessAndFileDataForEquality(rsd, bsd)
							if d {
								addInDiff = false
								break
							}

						}

						if addInDiff {
							diffReport.Clusters[cv.ClusterName].Namespaces[nv.NamespaceName].ResourceTypes[rtk].Resources[rsdk].SummaryData.ProcessData = append(diffReport.Clusters[cv.ClusterName].Namespaces[nv.NamespaceName].ResourceTypes[rtk].Resources[rsdk].SummaryData.ProcessData, rsd)
						}

					}

					for _, rsd := range rsdv.GetSummaryData().GetFileData() {

						if ignoreComparison(rsd.Source, rsd.Destination, ignorePath) {
							continue
						}
						addInDiff := true

						for _, bsd := range baselineReport.Clusters[ck].Namespaces[nk].ResourceTypes[rtk].Resources[rsdk].GetSummaryData().GetFileData() {
							d := compareProcessAndFileDataForEquality(rsd, bsd)
							if d {
								addInDiff = false
								break
							}

						}

						if addInDiff {
							diffReport.Clusters[cv.ClusterName].Namespaces[nv.NamespaceName].ResourceTypes[rtk].Resources[rsdk].SummaryData.FileData = append(diffReport.Clusters[cv.ClusterName].Namespaces[nv.NamespaceName].ResourceTypes[rtk].Resources[rsdk].SummaryData.FileData, rsd)
						}

					}

					for _, rsd := range rsdv.GetSummaryData().GetIngressConnection() {

						// TODO: check for allowed hosts

						addInDiff := true
						for _, bsd := range baselineReport.Clusters[ck].Namespaces[nk].ResourceTypes[rtk].Resources[rsdk].GetSummaryData().GetIngressConnection() {
							d := compareNetworkDataForEquality("ingress", rsd, bsd)
							if d {
								addInDiff = false
								break
							}

						}
						if addInDiff {
							diffReport.Clusters[cv.ClusterName].Namespaces[nv.NamespaceName].ResourceTypes[rtk].Resources[rsdk].SummaryData.IngressConnection = append(diffReport.Clusters[cv.ClusterName].Namespaces[nv.NamespaceName].ResourceTypes[rtk].Resources[rsdk].SummaryData.IngressConnection, rsd)
						}

					}

					for _, rsd := range rsdv.GetSummaryData().GetEgressConnection() {

						// TODO: check for allowed hosts
						addInDiff := true

						for _, bsd := range baselineReport.Clusters[ck].Namespaces[nk].ResourceTypes[rtk].Resources[rsdk].GetSummaryData().GetEgressConnection() {
							d := compareNetworkDataForEquality("egress", rsd, bsd)
							if d {
								addInDiff = false
								break
							}

						}

						if addInDiff {
							diffReport.Clusters[cv.ClusterName].Namespaces[nv.NamespaceName].ResourceTypes[rtk].Resources[rsdk].SummaryData.EgressConnection = append(diffReport.Clusters[cv.ClusterName].Namespaces[nv.NamespaceName].ResourceTypes[rtk].Resources[rsdk].SummaryData.EgressConnection, rsd)
						}

					}

					for _, rsd := range rsdv.GetSummaryData().GetBindConnection() {

						// TODO: check for allowed hosts
						addInDiff := true

						for _, bsd := range baselineReport.Clusters[ck].Namespaces[nk].ResourceTypes[rtk].Resources[rsdk].GetSummaryData().GetBindConnection() {
							d := compareNetworkDataForEquality("bind", rsd, bsd)
							if d {
								addInDiff = false
								break
							}

						}

						if addInDiff {
							diffReport.Clusters[cv.ClusterName].Namespaces[nv.NamespaceName].ResourceTypes[rtk].Resources[rsdk].SummaryData.BindConnection = append(diffReport.Clusters[cv.ClusterName].Namespaces[nv.NamespaceName].ResourceTypes[rtk].Resources[rsdk].SummaryData.BindConnection, rsd)
						}

					}

				}
			}
		}
	}

	diffReportJson, err := json.MarshalIndent(diffReport, "", "  ")
	if err != nil {
		return err
	}
	//fmt.Printf("%s", diffReportJson)

	// Write in a temp file
	tempF, err := os.CreateTemp("/tmp/", "diff-report")
	if err != nil {
		return err
	}

	if _, err := tempF.Write(diffReportJson); err != nil {
		return err
	}

	// TODO: Definitely needs to change this, very bad code and not want to be in this format
	getTableOutput(diffReport)
	return nil
}

// TODO: can return summary data i.e rsd but now returning bool
func compareProcessAndFileDataForEquality(rsd, bsd *opb.SysProcFileSummaryData) bool {
	if rsd.Source == bsd.Source && rsd.Destination == bsd.Destination && rsd.Status == bsd.Status {
		return true
	}
	return false
}

// TODO: can return summary data i.e rsd but now returning bool
func compareNetworkDataForEquality(nwType string, rsd, bsd *opb.SysNwSummaryData) bool {
	if (nwType == "egress" || nwType == "ingress") && (rsd.Protocol == bsd.Protocol && rsd.IP == bsd.IP && rsd.Command == bsd.Command && rsd.Port == bsd.Port && rsd.Labels != bsd.Labels && rsd.Namespace == bsd.Namespace) {
		return true
	}

	if nwType == "bind" && (rsd.Protocol == bsd.Protocol && rsd.IP == bsd.IP && rsd.Command == bsd.Command && rsd.Port == bsd.Port && rsd.BindPort == bsd.BindPort && rsd.BindAddress == bsd.BindAddress) {
		return true
	}

	return false
}

func ignoreComparison(sPath, dPath string, ignorePath []string) bool {
	ignoreCmp := false

	for _, iPath := range ignorePath {
		if strings.Contains(sPath, iPath) || strings.Contains(dPath, iPath) {
			ignoreCmp = true
			break
		}
	}
	return ignoreCmp
}

func readBaselineReportJson(filePath string) *rpb.ReportResponse {

	jsonFile, err := os.Open(filePath)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {

		}
	}(jsonFile)

	byteValue, _ := ioutil.ReadAll(jsonFile)
	t := &rpb.ReportResponse{}
	_ = json.Unmarshal(byteValue, t)
	return t
}
