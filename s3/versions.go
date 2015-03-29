package s3

type Versions []Version

func (v Versions) Less(i, j int) bool {
	return v[i].Version < v[j].Version
}

func (v Versions) Len() int {
	return len(v)
}

func (v Versions) Swap(i, j int) {
	temp := v[i]
	v[i] = v[j]
	v[j] = temp
}
