package entities

// Sequence entity for sequence objects
type Sequence struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	Timestamp int64  `json:"timestamp"`
}
