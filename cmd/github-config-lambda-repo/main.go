package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	lmb "github.com/suzuki-shunsuke/github-config/pkg/lambda"
)

func main() {
	handler := lmb.Handler{}
	lambda.Start(handler.StartRepo)
}
