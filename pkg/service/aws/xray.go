package aws

import (
	"log"

	"github.com/aws/aws-xray-sdk-go/strategy/sampling"
	"github.com/aws/aws-xray-sdk-go/xray"
)

func InitXrayConfig() {
	log.Println("init x-ray configudaration")
	s, _ := sampling.NewCentralizedStrategyWithFilePath("aws/rules.json")
	xray.Configure(xray.Config{SamplingStrategy: s})
}
