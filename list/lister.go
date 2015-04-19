package list

import "github.com/tscolari/s3kup/s3"

type Lister struct {
	s3 S3Client
}

type S3Client interface {
	List(path string) (versions s3.Versions, err error)
}

func New(client S3Client) Lister {
	return Lister{
		s3: client,
	}
}

func (l Lister) List(path string) (s3.Versions, error) {
	versions, err := l.s3.List(path)
	return versions, err
}
