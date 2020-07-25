package main

import (
	"link-shortener/src/internal/platform/env"
	"link-shortener/src/internal/platform/link"
	"link-shortener/src/internal/redirect"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	sess := session.Must(session.NewSession())
	s3API := s3.New(sess)

	redirectService := redirect.NewService(s3API, env.Get(env.BUCKET_NAME))
	linkService := link.NewServie(env.Get(env.BUCKET_DOMAIN))

	handler := NewHandler(redirectService, linkService)
	lambda.Start(handler)
}
