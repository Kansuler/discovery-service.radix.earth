package api

import (
	"fmt"
	"os"

	"github.com/imroc/req"
)

type GetDemandPayload struct {
	Result struct {
		TPS int64 `json:"tps"`
	} `json:"result"`
	Id      int64  `json:"id"`
	JsonRPC string `json:"jsonrpc"`
}

type GetDemandRequest struct {
	JsonRPC string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Params  struct{} `json:"params"`
	Id      int64    `json:"id"`
}

func GetDemand() (result *GetDemandPayload, err error) {
	response, err := req.Post(fmt.Sprintf("%s/archive", os.Getenv("ARCHIVE_NODE_API_URL")), req.BodyJSON(&GetDemandRequest{
		JsonRPC: "2.0",
		Method:  "network.get_demand",
		Id:      1,
	}))
	if err != nil {
		return nil, err
	}

	err = response.ToJSON(&result)
	return
}
