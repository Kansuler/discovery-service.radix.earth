package api

import (
	"fmt"
	"os"

	"github.com/imroc/req"
)

type Channel struct {
	LocalPort int64  `json:"localPort"`
	IP        string `json:"ip"`
	Type      string `json:"type"`
	URI       string `json:"uri"`
}

type PeerEntry struct {
	Address  string    `json:"address"`
	Channels []Channel `json:"channels"`
}

type GetPeersPayload struct {
	Result  []PeerEntry `json:"result"`
	Id      int64       `json:"id"`
	JsonRPC string      `json:"jsonrpc"`
}

type GetPeersRequest struct {
	JsonRPC string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Params  struct{} `json:"params"`
	Id      int64    `json:"id"`
}

func GetPeers() (result *GetPeersPayload, err error) {
	response, err := req.Post(fmt.Sprintf("%s/system", os.Getenv("ARCHIVE_NODE_SYSTEM_URL")), req.BodyJSON(&GetPeersRequest{
		JsonRPC: "2.0",
		Method:  "networking.get_peers",
		Id:      1,
	}))
	if err != nil {
		return nil, err
	}

	err = response.ToJSON(&result)
	return
}
