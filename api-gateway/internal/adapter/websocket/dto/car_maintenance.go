package wsdto

type CarMaintenanceEventMessage struct {
	CarID      string `json:"carID"`
	TemplateID string `json:"templateID"`
	RecordID   string `json:"recordID"`
	EventType  string `json:"eventType"`
	OccurredAt string `json:"occurredAt"`
}
