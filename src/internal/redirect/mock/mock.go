package mock

//go:generate mockgen -destination=./api.go -package=mock github.com/aws/aws-sdk-go/service/s3/s3iface S3API
