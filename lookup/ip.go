package lookup

import (
	"fmt"
	"os"

	"github.com/imroc/req"
)

type Location struct {
	Latitude     float64 `json:"lat"`
	Longitude    float64 `json:"lon"`
	Country      string  `json:"country"`
	City         string  `json:"city"`
	Region       string  `json:"region"`
	Query        string  `json:"query"`
	ISP          string  `json:"isp"`
	Organisation string  `json:"org"`
}

func IP(ipNumbers []string) (result *[]Location, err error) {
	param := req.Param{
		"fields": "lat,lon,country,city,region,query,isp,org",
	}
	response, err := req.Post(fmt.Sprintf("%s/batch", os.Getenv("IP_API_URL")), req.BodyJSON(&ipNumbers), param)
	if err != nil {
		return nil, err
	}

	err = response.ToJSON(&result)
	return
}
