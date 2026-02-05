package plugincontract

// EventPayload represents a system event delivered to plugins via SQS.
type EventPayload struct {
	EventType  string `json:"eventType"`
	OccurredAt string `json:"occurredAt"`
	AccountID  string `json:"accountId"`
	Data       Args   `json:"data,omitempty"`
}
