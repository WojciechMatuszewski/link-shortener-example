package env

import (
	"fmt"
	"os"
)

const (
	BUCKET_NAME   = "BUCKET_NAME"
	BUCKET_DOMAIN = "BUCKET_DOMAIN"
)

func Get(key string) string {
	v, found := os.LookupEnv(key)
	if !found {
		panic(fmt.Errorf("environment value for key %q not found", key))
	}
	return v
}
