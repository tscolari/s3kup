package s3

import (
	"regexp"

	goamzs3 "github.com/mitchellh/goamz/s3"
)

type Version struct {
	Path         string
	BackupName   string
	Version      string
	LastModified string
	Size         int64
}

func NewVersion(key goamzs3.Key) Version {
	backupNameRegexp, _ := regexp.Compile("^(.*)/\\d+$")
	backupName := backupNameRegexp.FindStringSubmatch(key.Key)
	versionRegexp, _ := regexp.Compile("^.*/(\\d+)$")
	version := versionRegexp.FindStringSubmatch(key.Key)

	return Version{
		Path:         key.Key,
		BackupName:   backupName[1],
		Version:      version[1],
		LastModified: key.LastModified,
		Size:         key.Size,
	}
}
