package redirect_test

import (
	"context"
	"errors"
	"testing"

	"link-shortener/src/internal/platform/link"
	"link-shortener/src/internal/redirect"
	"link-shortener/src/internal/redirect/mock"

	"github.com/stretchr/testify/require"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/golang/mock/gomock"
)

func TestService_Create(t *testing.T) {
	ctx := context.Background()

	t.Run("failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		api := mock.NewMockS3API(ctrl)
		redirectService := redirect.NewService(api, "bucketName")
		linkService := link.NewServie("https://wojtek.tk")

		l, _ := linkService.New("https://google.com")
		api.EXPECT().PutObjectWithContext(ctx, &s3.PutObjectInput{
			Body:                    nil,
			Bucket:                  aws.String("bucketName"),
			Key:                     aws.String(l.ID),
			WebsiteRedirectLocation: aws.String(l.Origin),
		}).Return(nil, errors.New("boom"))

		err := redirectService.Create(ctx, l)
		require.Error(t, err)
	})

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		api := mock.NewMockS3API(ctrl)
		redirectService := redirect.NewService(api, "bucketName")
		linkService := link.NewServie("https://wojtek.tk")

		l, _ := linkService.New("https://google.com")
		api.EXPECT().PutObjectWithContext(ctx, &s3.PutObjectInput{
			Body:                    nil,
			Bucket:                  aws.String("bucketName"),
			Key:                     aws.String(l.ID),
			WebsiteRedirectLocation: aws.String(l.Origin),
		}).Return(nil, nil)

		err := redirectService.Create(ctx, l)
		require.NoError(t, err)
	})
}

func TestService_Exists(t *testing.T) {
	ctx := context.Background()
	t.Run("does not exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		api := mock.NewMockS3API(ctrl)
		redirectService := redirect.NewService(api, "bucketName")
		linkService := link.NewServie("https://wojtek.tk")

		l, _ := linkService.New("https://google.com")
		api.EXPECT().HeadObjectWithContext(ctx, &s3.HeadObjectInput{
			Bucket: aws.String("bucketName"),
			Key:    aws.String(l.ID),
		}).Return(nil, errors.New(s3.ErrCodeNoSuchKey))

		exists, err := redirectService.Exists(ctx, l)
		require.NoError(t, err)
		require.False(t, exists)
	})

	t.Run("failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		api := mock.NewMockS3API(ctrl)
		redirectService := redirect.NewService(api, "bucketName")
		linkService := link.NewServie("https://wojtek.tk")

		l, _ := linkService.New("https://google.com")
		api.EXPECT().HeadObjectWithContext(ctx, &s3.HeadObjectInput{
			Bucket: aws.String("bucketName"),
			Key:    aws.String(l.ID),
		}).Return(nil, errors.New("boom"))

		exists, err := redirectService.Exists(ctx, l)
		require.Error(t, err)
		require.False(t, exists)
	})

	t.Run("already exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		api := mock.NewMockS3API(ctrl)
		redirectService := redirect.NewService(api, "bucketName")
		linkService := link.NewServie("https://wojtek.tk")

		l, _ := linkService.New("https://google.com")
		api.EXPECT().HeadObjectWithContext(ctx, &s3.HeadObjectInput{
			Bucket: aws.String("bucketName"),
			Key:    aws.String(l.ID),
		}).Return(nil, nil)

		exists, err := redirectService.Exists(ctx, l)
		require.NoError(t, err)
		require.True(t, exists)
	})
}
