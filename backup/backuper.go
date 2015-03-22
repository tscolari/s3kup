package backup

import (
	"fmt"
	"sort"
	"time"
)

type S3Client interface {
	Store(path string, content []byte) error
	List(path string) (files []string, err error)
	Delete(path string) error
}

type Backuper struct {
	s3Client       S3Client
	versionsToKeep int
}

func New(s3Client S3Client, versionsToKeep int) Backuper {
	return Backuper{
		s3Client:       s3Client,
		versionsToKeep: versionsToKeep,
	}
}

func (b Backuper) Backup(fileName string, fileContent []byte) error {
	err := b.putFile(fileName, fileContent)
	if err != nil {
		return err
	}

	return b.cleanUpOldVersions(fileName)
}

func (b Backuper) putFile(fileName string, fileContent []byte) error {
	timestamp := time.Now().UnixNano()
	fileName = fmt.Sprintf("%s/%d", fileName, timestamp)
	return b.s3Client.Store(fileName, fileContent)
}

func (b Backuper) cleanUpOldVersions(fileName string) error {
	storedVersions, err := b.s3Client.List(fileName)
	if err != nil {
		return err
	}

	sort.Strings(storedVersions)
	if len(storedVersions) >= b.versionsToKeep {
		extraVersions := len(storedVersions) - b.versionsToKeep
		for _, version := range storedVersions[:extraVersions] {
			err = b.s3Client.Delete(version)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
