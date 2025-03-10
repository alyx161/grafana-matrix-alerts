package actions

import "time"

type AlertsData struct {
	Receiver    string        `json:"receiver"`
	Status      string        `json:"status"`
	Alerts      []SingleAlert `json:"alerts"`
	GroupLabels struct {
		Alertname     string `json:"alertname"`
		GrafanaFolder string `json:"grafana_folder"`
	} `json:"groupLabels"`
	CommonLabels struct {
		Alertname     string `json:"alertname"`
		GrafanaFolder string `json:"grafana_folder"`
	} `json:"commonLabels"`
	CommonAnnotations struct {
		Summary string `json:"summary"`
	} `json:"commonAnnotations"`
	ExternalURL     string `json:"externalURL"`
	Version         string `json:"version"`
	GroupKey        string `json:"groupKey"`
	TruncatedAlerts int    `json:"truncatedAlerts"`
	OrgID           int    `json:"orgId"`
	Title           string `json:"title"`
	State           string `json:"state"`
	Message         string `json:"message"`
}

type SingleAlert struct {
	Status string `json:"status"`
	Labels struct {
		Alertname     string `json:"alertname"`
		GrafanaFolder string `json:"grafana_folder"`
	} `json:"labels"`
	Annotations struct {
		Summary string `json:"summary"`
	} `json:"annotations"`
	StartsAt     time.Time          `json:"startsAt"`
	EndsAt       time.Time          `json:"endsAt"`
	GeneratorURL string             `json:"generatorURL"`
	Fingerprint  string             `json:"fingerprint"`
	SilenceURL   string             `json:"silenceURL"`
	DashboardURL string             `json:"dashboardURL"`
	PanelURL     string             `json:"panelURL"`
	Values       map[string]float64 `json:"values"`
	ValueString  string             `json:"valueString"`
}
