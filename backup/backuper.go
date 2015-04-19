package backup

import (
	"fmt"
	"sort"
	"time"

	"github.com/tscolari/s3up/log"
	"github.com/tscolari/s3up/s3"
)

type Backuper struct {
	s3Client       S3Client
	versionsToKeep int
}

type S3Client interface {
	Store(path string, content []byte) error
	List(path string) (versions s3.Versions, err error)
	Delete(path string) error
}

func New(s3Client S3Client, versionsToKeep int) Backuper {
	return Backuper{
		s3Client:       s3Client,
		versionsToKeep: versionsToKeep,
	}
}

func (b Backuper) Backup(fileName string, fileContent []byte) error {
	log.Info("Started backup of", fileName)
	err := b.putFile(fileName, fileContent)
	if err != nil {
		return err
	}

	return b.cleanUpOldVersions(fileName)
}

func (b Backuper) putFile(fileName string, fileContent []byte) error {
	timestamp := time.Now().UnixNano()
	fileName = fmt.Sprintf("%s/%d", fileName, timestamp)
	log.Info(" -- File version:", timestamp)
	return b.s3Client.Store(fileName, fileContent)
}

func (b Backuper) cleanUpOldVersions(fileName string) error {
	log.Info(" -- Looking for old versions to delete. keeping", b.versionsToKeep)
	storedVersions, err := b.s3Client.List(fileName)
	if err != nil {
		return err
	}

	sort.Sort(storedVersions)
	if len(storedVersions) >= b.versionsToKeep {
		extraVersions := len(storedVersions) - b.versionsToKeep
		log.Info(" --", extraVersions, "old versions will be deleted")
		for _, version := range storedVersions[:extraVersions] {
			err = b.s3Client.Delete(version.Path)
			log.Info(" -- deleted:", version.Version)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
