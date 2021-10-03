package api

import (
	"fmt"
	"os"

	"github.com/imroc/req"
)

type GetThroughputPayload struct {
	Result struct {
		TPS int64 `json:"tps"`
	} `json:"result"`
	Id      int64  `json:"id"`
	JsonRPC string `json:"jsonrpc"`
}

type GetThroughputRequest struct {
	JsonRPC string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Params  struct{} `json:"params"`
	Id      int64    `json:"id"`
}

func GetThroughput() (result *GetThroughputPayload, err error) {
	response, err := req.Post(fmt.Sprintf("%s/archive", os.Getenv("ARCHIVE_NODE_API_URL")), req.BodyJSON(&GetThroughputRequest{
		JsonRPC: "2.0",
		Method:  "network.get_throughput",
		Id:      1,
	}))
	if err != nil {
		return nil, err
	}

	err = response.ToJSON(&result)
	return
}
