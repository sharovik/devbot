package dto

//SlackRequestAcknowledge the request for acknowledging of the event message
type SlackRequestAcknowledge struct {
	EnvelopeID string `json:"envelope_id"`
	Payload    string `json:"payload"`
}
