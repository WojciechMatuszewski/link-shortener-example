package redirect

import (
	"context"
	"strings"

	"link-shortener/src/internal/platform/link"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type Service struct {
	api        s3iface.S3API
	bucketName string
}

func NewService(api s3iface.S3API, bucketName string) *Service {
	return &Service{api: api, bucketName: bucketName}
}

func (s *Service) Exists(ctx context.Context, l link.Link) (bool, error) {
	_, err := s.api.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(l.ID),
	})
	if err != nil {
		if strings.Contains(err.Error(), s3.ErrCodeNoSuchKey) {
			return false, nil
		}
		return false, errors.Wrapf(err, "checking for redirect existence based on link: %q", l)
	}

	return true, nil
}

func (s *Service) Create(ctx context.Context, l link.Link) error {
	_, err := s.api.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Body:                    nil,
		Bucket:                  aws.String(s.bucketName),
		Key:                     aws.String(l.ID),
		WebsiteRedirectLocation: aws.String(l.Origin),
	})
	if err != nil {
		return errors.Wrapf(err, " creating redirect for link: %q", l)
	}
	return nil
}
