package report

import (
	"fmt"
	"github.com/accuknox/accuknox-cli/summary"
	opb "github.com/accuknox/auto-policy-discovery/src/protobuf/v1/observability"
	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"github.com/mgutz/ansi"
	"github.com/olekukonko/tablewriter"
	"os"
	"sort"
)

var (
	// SysProcHeader variable contains source process, destination process path, count, timestamp and status
	SysProcHeader = []string{"Src Process", "Destination Process Path", "Status"}
	// SysFileHeader variable contains source process, destination file path, count, timestamp and status
	SysFileHeader = []string{"Src Process", "Destination File Path", "Status"}
	// SysNwHeader variable contains protocol, command, POD/SVC/IP, Port, Namespace, and Labels
	SysNwHeader = []string{"Protocol", "Command", "POD/SVC/IP", "Port", "Namespace", "Labels"}
	// SysBindNwHeader variable contains protocol, command, Bind Port, Bind Address, count and timestamp
	SysBindNwHeader = []string{"Protocol", "Command", "Bind Port", "Bind Address"}
)

// DisplayReportOutput function
func DisplayReportOutput(resp *opb.Response, revDNSLookup bool, resourceType, resourceName string) {

	if len(resp.ProcessData) <= 0 && len(resp.FileData) <= 0 && len(resp.IngressConnection) <= 0 && len(resp.EgressConnection) <= 0 {
		return
	}

	writeResourceInfoToTable(resourceType, resourceName, resp.Namespace, resp.ClusterName, resp.ContainerName, resp.Label)
	diffReportMD += "-------\n"

	writeClusterInfoToMdVar(resourceType, resourceName, resp.Namespace, resp.ClusterName, resp.ContainerName, resp.Label)

	// Colored Status for Allow and Deny
	agc := ansi.ColorFunc("green")
	arc := ansi.ColorFunc("red")
	ayc := ansi.ColorFunc("yellow")

	if len(resp.ProcessData) > 0 || len(resp.FileData) > 0 {
		diffReportMD += "## System access behavior Summary \n "

		if len(resp.ProcessData) > 0 {
			procRowData := [][]string{}
			diffReportMD += "### Process Data \n \n "

			// Display process data
			for _, procData := range resp.ProcessData {
				procStrSlice := []string{}
				procStrSlice = append(procStrSlice, procData.Source)
				procStrSlice = append(procStrSlice, procData.Destination)
				if procData.Status == "Allow" {
					procStrSlice = append(procStrSlice, agc(procData.Status))
				} else if procData.Status == "Audit" {
					procStrSlice = append(procStrSlice, ayc(procData.Status))
				} else {
					procStrSlice = append(procStrSlice, arc(procData.Status))
				}
				procRowData = append(procRowData, procStrSlice)
			}
			sort.Slice(procRowData[:], func(i, j int) bool {
				for x := range procRowData[i] {
					if procRowData[i][x] == procRowData[j][x] {
						continue
					}
					return procRowData[i][x] < procRowData[j][x]
				}
				return false
			})
			//WriteTable(SysProcHeader, procRowData)
			tableOuput(SysProcHeader, procRowData, "Process Data\n")
			fmt.Printf("\n")
		}

		if len(resp.FileData) > 0 {
			// Display file data
			diffReportMD += "### File Access Data \n \n "

			fileRowData := [][]string{}
			for _, fileData := range resp.FileData {
				fileStrSlice := []string{}
				fileStrSlice = append(fileStrSlice, fileData.Source)
				fileStrSlice = append(fileStrSlice, fileData.Destination)
				if fileData.Status == "Allow" {
					fileStrSlice = append(fileStrSlice, agc(fileData.Status))
				} else if fileData.Status == "Audit" {
					fileStrSlice = append(fileStrSlice, ayc(fileData.Status))
				} else {
					fileStrSlice = append(fileStrSlice, arc(fileData.Status))
				}
				fileRowData = append(fileRowData, fileStrSlice)
			}
			sort.Slice(fileRowData[:], func(i, j int) bool {
				for x := range fileRowData[i] {
					if fileRowData[i][x] == fileRowData[j][x] {
						continue
					}
					return fileRowData[i][x] < fileRowData[j][x]
				}
				return false
			})
			//WriteTable(SysFileHeader, fileRowData)
			tableOuput(SysFileHeader, fileRowData, "File Access Data\n")

			fmt.Printf("\n")
		}
	}
	if len(resp.IngressConnection) > 0 || len(resp.EgressConnection) > 0 || len(resp.BindConnection) > 0 {
		diffReportMD += "## Network Behavior Summary\n "
		if len(resp.IngressConnection) > 0 {
			// Display server conn data
			diffReportMD += "### Ingress Connections \n \n "
			inNwRowData := [][]string{}
			for _, ingressConnection := range resp.IngressConnection {
				inNwStrSlice := []string{}
				domainName := summary.DnsLookup(ingressConnection.IP, revDNSLookup)
				inNwStrSlice = append(inNwStrSlice, ingressConnection.Protocol)
				inNwStrSlice = append(inNwStrSlice, ingressConnection.Command)
				inNwStrSlice = append(inNwStrSlice, domainName)
				inNwStrSlice = append(inNwStrSlice, ingressConnection.Port)
				inNwStrSlice = append(inNwStrSlice, ingressConnection.Namespace)
				inNwStrSlice = append(inNwStrSlice, ingressConnection.Labels)
				inNwRowData = append(inNwRowData, inNwStrSlice)
			}
			//WriteTable(SysNwHeader, inNwRowData)
			tableOuput(SysNwHeader, inNwRowData, "Ingress Connections\n")

			fmt.Printf("\n")
		}

		if len(resp.EgressConnection) > 0 {
			// Display server conn data
			diffReportMD += "### Egress Connections \n \n "
			outNwRowData := [][]string{}
			for _, egressConnection := range resp.EgressConnection {
				outNwStrSlice := []string{}
				domainName := summary.DnsLookup(egressConnection.IP, revDNSLookup)
				outNwStrSlice = append(outNwStrSlice, egressConnection.Protocol)
				outNwStrSlice = append(outNwStrSlice, egressConnection.Command)
				outNwStrSlice = append(outNwStrSlice, domainName)
				outNwStrSlice = append(outNwStrSlice, egressConnection.Port)
				outNwStrSlice = append(outNwStrSlice, egressConnection.Namespace)
				outNwStrSlice = append(outNwStrSlice, egressConnection.Labels)
				outNwRowData = append(outNwRowData, outNwStrSlice)
			}
			//WriteTable(SysNwHeader, outNwRowData)
			tableOuput(SysNwHeader, outNwRowData, "Egress Connections\n")

			fmt.Printf("\n")
		}

		if len(resp.BindConnection) > 0 {
			// Display bind connections details
			diffReportMD += "### Binds \n \n"

			bindNwRowData := [][]string{}
			for _, bindConnection := range resp.BindConnection {
				bindNwStrSlice := []string{}
				bindNwStrSlice = append(bindNwStrSlice, bindConnection.Protocol)
				bindNwStrSlice = append(bindNwStrSlice, bindConnection.Command)
				bindNwStrSlice = append(bindNwStrSlice, bindConnection.BindPort)
				bindNwStrSlice = append(bindNwStrSlice, bindConnection.BindAddress)
				bindNwRowData = append(bindNwRowData, bindNwStrSlice)
			}
			//WriteTable(SysBindNwHeader, bindNwRowData)
			tableOuput(SysBindNwHeader, bindNwRowData, "Binds\n")

			fmt.Printf("\n")
		}
	}
}

func writeClusterInfoToMdVar(resourceType, resourceName, namespace, clustername, containername, labels string) {

	//twInner := table.NewWriter()
	//twInner.AppendHeader(table.Row{"", ""})
	//twInner.AppendRow(table.Row{"Cluster Name", clustername})
	//twInner.AppendRow(table.Row{"Namespace Name", namespace})
	//twInner.AppendRow(table.Row{"Resource Type", resourceType})
	//twInner.AppendRow(table.Row{"Resource Name", resourceName})
	//twInner.AppendRow(table.Row{"Container Name", containername})
	//twInner.AppendRow(table.Row{"Labels", labels})
	////twInner.Style().Options = table.OptionsNoBorders
	//twInner.SetStyle(table.StyleLight)

	twOuter := table.NewWriter()
	twOuter.AppendHeader(table.Row{"Resource Information"})
	//twOuter.AppendRow(table.Row{twInner.RenderMarkdown()})
	//twOuter.SetAlignHeader([]text.Align{text.AlignCenter})
	//twOuter.SetStyle(table.StyleDouble)
	var s string
	s += "\n**Cluster Name** " + "\t\t" + "`" + clustername + "`"
	s += "\n**Namespace Name** " + "\t\t" + "`" + namespace + "`"
	s += "\n**Resource Type**" + "\t\t" + "`" + resourceType + "`"
	s += "\n**Resource Name**" + "\t\t" + "`" + resourceName + "`"
	s += "\n**Container Name**" + "\t\t" + "`" + containername + "`"
	s += "\n**Labels**" + "\t\t" + "`" + labels + "`"
	s += "\n"
	twOuter.AppendRow(table.Row{s})

	diffReportMD += twOuter.RenderMarkdown() + "\n\n"

}

// WriteTable function
func WriteTable(header []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}

func writeResourceInfoToTable(resourceType, resourceName, namespace, clustername, containername, labels string) {

	fmt.Printf("\n")

	podinfo := [][]string{
		{"Cluster Name", clustername},
		{"Namespace Name", namespace},
		{"Resource Type", resourceType},
		{"Resource Name", resourceName},
		{"Container Name", containername},
		{"Labels", labels},
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	for _, v := range podinfo {
		table.Append(v)
	}
	table.Render()
}

var diffReportMD string

func tableOuput(header []string, data [][]string, title string) {
	tr := table.NewWriter()
	var row table.Row
	for _, d := range header {
		row = append(row, d)
	}
	tr.AppendHeader(row)
	tr.SetAlignHeader([]text.Align{text.AlignCenter})
	for _, v := range data {
		row = table.Row{}
		for _, d := range v {
			row = append(row, d)
		}
		tr.AppendRow(row)
	}
	tr.SetStyle(table.StyleLight)
	tr.SetTitle(title)
	fmt.Printf("%s", tr.Render())
	//fmt.Printf("%s", tr.RenderMarkdown())
	tr.SetTitle("")

	diffReportMD += tr.RenderMarkdown() + "\n \n \n"

}
