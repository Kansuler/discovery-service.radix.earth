package jobs

type RequestPayload struct {
	JsonRPC string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Params  struct{} `json:"params"`
	Id      int64    `json:"id"`
}
