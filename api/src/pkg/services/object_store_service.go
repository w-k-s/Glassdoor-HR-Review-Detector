package services

import (
	"io"
)

type ObjectStoreService interface {
	Get(bucket string, file string) (io.Reader, error)
	Put(bucket string, file string, content string) error
}
