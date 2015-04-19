package fetch

import (
	"errors"
	"fmt"

	"github.com/tscolari/s3up/s3"
)

type Fetcher struct {
	s3 S3Client
}

type S3Client interface {
	List(path string) (versions s3.Versions, err error)
	Get(path string) ([]byte, error)
}

func New(client S3Client) Fetcher {
	return Fetcher{
		s3: client,
	}
}

func (f Fetcher) FetchLatest(backupName string) ([]byte, error) {
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

func (f Fetcher) FetchVersion(backupName string, version int64) ([]byte, error) {
	versionPath := fmt.Sprintf("%s/%d", backupName, version)

	content, err := f.s3.Get(versionPath)
	if err != nil && err.Error() == "The specified key does not exist." {
		message := fmt.Sprintf("Could not find version '%d'", version)
		err = errors.New(message)
	}

	return content, err
}
