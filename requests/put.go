package requests

// Request request payload format for adding records
type Request struct {
	Key       string `json:"key,omitempty"`
	Value     string `json:"value,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}
