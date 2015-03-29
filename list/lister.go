package list

import "github.com/tscolari/s3up/s3"

type Lister struct {
	s3 s3.S3Client
}

func New(client s3.S3Client) Lister {
	return Lister{
		s3: client,
	}
}

func (l Lister) List(path string) (s3.Versions, error) {
	versions, err := l.s3.List(path)
	return versions, err
}
