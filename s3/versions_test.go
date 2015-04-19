package s3_test

import (
	"sort"
	"time"

	"github.com/tscolari/s3kup/s3"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Versions", func() {

	createVersion := func(number int, date string) s3.Version {
		lastModified, err := time.Parse(time.RFC3339, date)
		Expect(err).ToNot(HaveOccurred())

		return s3.Version{
			Version:      int64(number),
			LastModified: lastModified,
		}
	}

	It("responds to the sort interface, sorting by LastModified field", func() {
		version1 := createVersion(1, "2019-03-29T11:54:42.819+01:00")
		version10 := createVersion(10, "2020-03-29T11:54:42.819+01:00")
		version50 := createVersion(50, "2018-03-29T11:54:42.819+01:00")
		version100 := createVersion(100, "2017-03-29T11:54:42.819+01:00")

		versions := s3.Versions{
			version50,
			version100,
			version1,
			version10,
		}

		sort.Sort(versions)

		Expect(versions).To(Equal(s3.Versions{
			version100,
			version50,
			version1,
			version10,
		}))
	})
})
