package s3_test

import (
	"sort"

	. "github.com/tscolari/s3up/s3"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Versions", func() {
	It("responds to the sort interface, sorting by version number", func() {
		version1 := Version{Version: "1"}
		version10 := Version{Version: "10"}
		version50 := Version{Version: "50"}
		version100 := Version{Version: "100"}

		versions := Versions{
			version50,
			version100,
			version1,
			version10,
		}

		sort.Sort(versions)

		Expect(versions).To(Equal(Versions{
			version1,
			version10,
			version50,
			version100,
		}))
	})
})
