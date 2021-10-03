package api

import (
	"fmt"
	"os"

	"github.com/imroc/req"
)

type Validator struct {
	TotalDeligatedStake     string `json:"totalDelegatedStake"`
	UptimePercentage        string `json:"uptimePercentage"`
	ProposalsMissed         int64  `json:"proposalsMissed"`
	Address                 string `json:"address"`
	InfoURL                 string `json:"infoURL"`
	OwnerDelegation         string `json:"ownerDeligation"`
	Name                    string `json:"name"`
	ValidatorFee            string `json:"validatorFee"`
	Registered              bool   `json:"registered"`
	OwnerAddress            string `json:"ownerAddress"`
	IsExternalStakeAccepted bool   `json:"isExternalStakeAccepted"`
	ProposalsCompleted      int64  `json:"proposalsCompleted"`
}

type GetValidatorSetResult struct {
	Cursor     string      `json:"cursor"`
	Validators []Validator `json:"validators"`
}

type GetValidatorSetPayload struct {
	Result  GetValidatorSetResult `json:"result"`
	Id      int64                 `json:"id"`
	JsonRPC string                `json:"jsonrpc"`
}

type GetValidatorSetParams struct {
	Size   int64  `json:"size"`
	Cursor string `json:"cursor"`
}

type GetValidatorSetRequest struct {
	JsonRPC string                `json:"jsonrpc"`
	Method  string                `json:"method"`
	Params  GetValidatorSetParams `json:"params"`
	Id      int64                 `json:"id"`
}

func GetValidatorSet() (result *GetValidatorSetPayload, err error) {
	response, err := req.Post(fmt.Sprintf("%s/archive", os.Getenv("ARCHIVE_NODE_API_URL")), req.BodyJSON(&GetValidatorSetRequest{
		JsonRPC: "2.0",
		Method:  "validators.get_next_epoch_set",
		Params: GetValidatorSetParams{
			Size:   1000,
			Cursor: "1",
		},
		Id: 1,
	}))
	if err != nil {
		return nil, err
	}

	err = response.ToJSON(&result)
	return
}
