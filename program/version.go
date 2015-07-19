package program

import (
	"fmt"
)

// Versioning: [major].[minor].[revision]

type VersionNo struct {
	Major    int
	Minor    int
	Revision int
}

func ParseVersionNumber(verstr string) VersionNo {
	vn := VersionNo{}
	fmt.Sscanf(verstr, "%d.%d.%d", &vn.Major, &vn.Minor, &vn.Revision)
	if vn.Major < 0 {
		vn.Major = 0
	}
	if vn.Minor < 0 {
		vn.Minor = 0
	}
	if vn.Revision < 0 {
		vn.Revision = 0
	}
	return vn
}

func (vn VersionNo) Compare(anotherVersion VersionNo) int {
	if vn.Major < anotherVersion.Major {
		return -1
	}
	if vn.Major > anotherVersion.Major {
		return 1
	}
	if vn.Minor < anotherVersion.Minor {
		return -1
	}
	if vn.Minor > anotherVersion.Minor {
		return 1
	}
	if vn.Revision < anotherVersion.Revision {
		return -1
	}
	if vn.Revision > anotherVersion.Revision {
		return 1
	}
	return 0
}

func (vn VersionNo) String() string {
	if vn.Revision == 0 {
		return fmt.Sprintf("%d.%d", vn.Major, vn.Minor)
	}
	return fmt.Sprintf("%d.%d.%d", vn.Major, vn.Minor, vn.Revision)
}
