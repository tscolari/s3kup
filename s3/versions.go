package s3

import (
	"log"
	"strconv"
)

type Versions []Version

func (v Versions) Less(i, j int) bool {
	version1, err := strconv.ParseInt(v[i].Version, 10, 64)
	if err != nil {
		log.Fatal("Remote version '", version1, "' can't be parsed.")
	}
	version2, err := strconv.ParseInt(v[j].Version, 10, 64)
	if err != nil {
		log.Fatal("Remote version '", version2, "' can't be parsed.")
	}

	return version1 < version2
}

func (v Versions) Len() int {
	return len(v)
}

func (v Versions) Swap(i, j int) {
	temp := v[i]
	v[i] = v[j]
	v[j] = temp
}
