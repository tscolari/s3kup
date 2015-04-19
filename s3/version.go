package s3

import (
	"errors"
	"regexp"
	"strconv"
	"time"

	goamzs3 "github.com/mitchellh/goamz/s3"
)

type Version struct {
	Path         string
	BackupName   string
	Version      int64
	LastModified time.Time
	Size         uint64
}

func NewVersion(key goamzs3.Key) (Version, error) {
	backupNameRegexp, _ := regexp.Compile("^(.*)/\\d+$")
	backupName := backupNameRegexp.FindStringSubmatch(key.Key)

	versionRegexp, _ := regexp.Compile("^.*/(\\d+)$")
	versionStr := versionRegexp.FindStringSubmatch(key.Key)
	if len(versionStr) < 2 {
		return Version{}, errors.New("Remote version '" + key.Key + "' can't be parsed")
	}

	versionInt, err := strconv.ParseInt(versionStr[1], 10, 64)
	if err != nil {
		return Version{}, errors.New("Remote version '" + versionStr[1] + "' can't be parsed")
	}

	lastModified, err := time.Parse(time.RFC3339, key.LastModified)
	if err != nil {
		return Version{}, errors.New("Failed to parse the version timestamp. '" + key.LastModified + "' was not recognized")
	}

	return Version{
		Path:         key.Key,
		BackupName:   backupName[1],
		Version:      versionInt,
		LastModified: lastModified,
		Size:         uint64(key.Size),
	}, nil
}
