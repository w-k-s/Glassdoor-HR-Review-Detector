package objectstore

import (
	"fmt"
	"io"

	"com.github/w-k-s/glassdoor-hr-review-detector/pkg/services"
)

type S3 struct{}

func MustS3() services.ObjectStoreService {
	return &S3{}
}

func (s *S3) Get(bucket string, object string) (io.Reader, error) {
	return nil, fmt.Errorf("TODO")
}

func (s *S3) Put(bucket string, object string, content string) error {
	return fmt.Errorf("TODO")
}
