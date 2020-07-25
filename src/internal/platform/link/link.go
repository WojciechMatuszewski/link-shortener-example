package link

import (
	"fmt"
	"net/url"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// ErrOriginNotValid is returned when origin is not a valid URI
var ErrOriginNotValid = errors.New("origin is not valid.")

type Link struct {
	ID     string
	Path   string
	Origin string
}

type Service struct {
	domain string
}

func NewServie(domain string) *Service {
	return &Service{domain: domain}
}

func (s *Service) New(origin string) (Link, error) {
	if !IsOriginValid(origin) {
		return Link{}, ErrOriginNotValid
	}

	id := newID()
	return Link{
		ID:     id,
		Path:   newPath(s.domain, id),
		Origin: origin,
	}, nil
}

func (l Link) String() string {
	return l.Path
}

// IsOriginValid checks if the link is valid.
func IsOriginValid(origin string) bool {
	_, err := url.ParseRequestURI(origin)
	return err == nil
}

func newPath(domain, id string) string {
	return fmt.Sprintf("%v/%v", domain, id)
}

func newID() string {
	return uuid.Must(uuid.NewRandom()).String()[:7]
}
