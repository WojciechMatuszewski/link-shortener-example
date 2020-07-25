package main

import (
	"context"
	"encoding/json"
	"net/http"

	"link-shortener/src/internal/platform/link"

	"github.com/pkg/errors"

	"github.com/aws/aws-lambda-go/events"
)

type Handler func(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error)

type Body struct {
	Origin string `json:"origin"`
}

type Response struct {
	Path string `json:"url"`
}

type Redirector interface {
	Exists(ctx context.Context, l link.Link) (bool, error)
	Create(ctx context.Context, l link.Link) error
}

type LinkCreator interface {
	New(origin string) (link.Link, error)
}

func NewHandler(redirector Redirector, linkCreator LinkCreator) Handler {
	return func(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		var body Body
		err := json.Unmarshal([]byte(request.Body), &body)
		if err != nil {
			return newResponse(http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError),
				errors.Wrapf(err, "marshalling body %v", request.Body)
		}

		l, err := linkCreator.New(body.Origin)
		if err != nil {
			if errors.Cause(err) == link.ErrOriginNotValid {
				return newResponse(http.StatusText(http.StatusBadRequest), http.StatusBadRequest),
					nil
			}

			return newResponse(http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError),
				errors.Wrapf(err, "creating link %v", request.Body)
		}

		objExists, err := redirector.Exists(ctx, l)
		if err != nil {
			return newResponse(http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError),
				errors.Wrapf(err, "checking if object exists %v", l)
		}
		if objExists {
			return newResponse(http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError),
				errors.Wrapf(err, "object already exists %v", l)
		}

		err = redirector.Create(ctx, l)
		if err != nil {
			return newResponse(http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError),
				errors.Wrapf(err, "creating link %v", request.Body)
		}

		resp := Response{Path: l.Path}
		buf, err := json.Marshal(resp)
		if err != nil {
			return newResponse(http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError),
				errors.Wrap(err, "marshaling the response")
		}

		return newResponse(string(buf), http.StatusOK), nil
	}
}

func newResponse(body string, statusCode int) events.APIGatewayV2HTTPResponse {
	return events.APIGatewayV2HTTPResponse{Body: body, StatusCode: statusCode}
}
