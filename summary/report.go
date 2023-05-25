package summary

type ProcessData struct {
	Source      string `json:"Source"`
	Destination string `json:"Destination"`
	Status      string `json:"Status"`
}

type SummaryData struct {
	PodName       string        `json:"PodName,omitempty"`
	ClusterName   string        `json:"ClusterName,omitempty"`
	Label         string        `json:"Label,omitempty"`
	ContainerName string        `json:"ContainerName,omitempty"`
	ProcessData   []ProcessData `json:"ProcessData,omitempty"`
}

func GetSummaryData() {

}
