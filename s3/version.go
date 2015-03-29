package s3

import (
	"log"
	"regexp"
	"strconv"

	goamzs3 "github.com/mitchellh/goamz/s3"
)

type Version struct {
	Path         string
	BackupName   string
	Version      int64
	LastModified string
	Size         int64
}

func NewVersion(key goamzs3.Key) Version {
	backupNameRegexp, _ := regexp.Compile("^(.*)/\\d+$")
	backupName := backupNameRegexp.FindStringSubmatch(key.Key)

	versionRegexp, _ := regexp.Compile("^.*/(\\d+)$")
	versionStr := versionRegexp.FindStringSubmatch(key.Key)
	versionInt, err := strconv.ParseInt(versionStr[1], 10, 64)
	if err != nil {
		log.Fatal("Remote version '", versionStr, "' can't be parsed.")
	}

	return Version{
		Path:         key.Key,
		BackupName:   backupName[1],
		Version:      versionInt,
		LastModified: key.LastModified,
		Size:         key.Size,
	}
}
