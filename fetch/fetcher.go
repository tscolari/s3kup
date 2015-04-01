package fetch

import (
	"errors"
	"fmt"

	"github.com/tscolari/s3up/s3"
)

type Fetcher struct {
	s3 s3.S3Client
}

func New(client s3.S3Client) Fetcher {
	return Fetcher{
		s3: client,
	}
}

func (f Fetcher) Fetch(backupName string) ([]byte, error) {
	versions, err := f.s3.List(backupName)
	if err != nil {
		return nil, err
	}

	if len(versions) == 0 {
		message := fmt.Sprintf("There's no backup named '%s' on this bucket", backupName)
		return nil, errors.New(message)
	}

	lastVersion := versions[len(versions)-1]
	return f.s3.Get(lastVersion.Path)
}
