package report

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/accuknox/accuknox-cli/k8s"
	rpb "github.com/accuknox/accuknox-cli/rpb"
	"github.com/accuknox/accuknox-cli/utils"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"strconv"
)

// DefaultReqType : default option for request type
var DefaultReqType = "process,file,network"
var matchLabels = map[string]string{"app": "discovery-engine"}
var port int64 = 9089

type Options struct {
	GRPC             string
	Clusters         []string
	Namespaces       []string
	ResourceType     []string
	ResourceName     []string
	Labels           string
	ContainerName    string
	PodName          string
	Source           []string
	Destination      []string
	Operation        string
	IgnorePaths      []string
	BaselineJsonPath string
}

func GetReport(c *k8s.Client, o *Options) (*rpb.ReportResponse, error) {
	gRPC := ""
	targetSvc := "discovery-engine"
	if o.GRPC != "" {
		gRPC = o.GRPC
	} else {
		if val, ok := os.LookupEnv("DISCOVERY_SERVICE"); ok {
			gRPC = val
		} else {
			pf, err := utils.InitiatePortForward(c, port, port, matchLabels, targetSvc)
			if err != nil {
				return nil, err
			}
			gRPC = "localhost:" + strconv.FormatInt(pf.LocalPort, 10)
		}
	}

	req := &rpb.ReportRequest{
		Clusters:     o.Clusters,
		Namespaces:   o.Namespaces,
		ResourceType: o.ResourceType,
		ResourceName: o.ResourceName,
		PodName:      o.PodName,
		MetaData: &rpb.MetaData{
			Label:         o.Labels,
			ContainerName: o.ContainerName,
		},
		Operation:   o.Operation,
		Source:      o.Source,
		Destination: o.Destination,
	}

	// create a client
	conn, err := grpc.Dial(gRPC, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.New("could not connect to the server. Possible troubleshooting:\n- Check if discovery engine is running\n- kubectl get po -n accuknox-agents")
	}
	defer conn.Close()

	client := rpb.NewReportClient(conn)

	res, err := client.GetReport(context.Background(), req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (o *Options) Report(c *k8s.Client) error {

	report, err := GetReport(c, o)
	if err != nil {
		log.Error().Msgf("error while getting report, error: %s", err.Error())
		return err
	}
	reportJson, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	//fmt.Printf("%s \n", reportJson)

	// Write in a temp file
	tempF, err := os.CreateTemp("/tmp/", "report-*.json")
	if err != nil {
		return err
	}

	if _, err := tempF.Write(reportJson); err != nil {
		return err
	}

	baselineReport := readBaselineReportJson(o.BaselineJsonPath)

	err = getDiff(baselineReport, report, o.IgnorePaths)

	if err != nil {
		return err
	}

	return nil
}
