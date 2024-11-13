package version

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// NodeVersion represents the current version of the node software.
// It should be updated with each new release, adjusting the Major, Minor, or Patch version numbers as necessary.
// When a major release occurs, the Meta field should be cleared (set to an empty string).
// For release candidates, set the Meta field to "rc1", "rc2", and so on.
// During development, set the Meta field to "beta".
var NodeVersion = Version{
	Major: 1,
	Minor: 6,
	Patch: 0,
	Meta:  "",
	Alias: "Mumbai",
}

// Version defines the version of Pactus software.
// It follows the semantic versioning 2.0.0 spec (http://semver.org/).
type Version struct {
	Major uint   // Major version number
	Minor uint   // Minor version number
	Patch uint   // Patch version number
	Meta  string // Metadata for version (e.g., "beta", "rc1")
	Alias string // Alias for version (e.g., "London")
}

// ParseVersion parses a version string into a Version struct.
// The format should be "Major.Minor.Patch-Meta", where Meta is optional.
// Returns the parsed Version struct and an error if parsing fails.
func ParseVersion(versionStr string) (Version, error) {
	var ver Version

	if versionStr[0] == 'v' {
		versionStr = versionStr[1:]
	}
	// Split the version string into parts
	parts := strings.Split(versionStr, ".")
	if len(parts) != 3 {
		return ver, errors.New("invalid version format")
	}

	// Parse Major version
	major, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return ver, fmt.Errorf("failed to parse Major version: %w", err)
	}
	ver.Major = uint(major)

	// Parse Minor version
	minor, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return ver, fmt.Errorf("failed to parse Minor version: %w", err)
	}
	ver.Minor = uint(minor)

	// Parse Patch version and Meta (if present)
	patchMeta := strings.Split(parts[2], "-")
	if len(patchMeta) > 2 {
		return ver, errors.New("invalid Patch and Meta format")
	}

	patch, err := strconv.ParseUint(patchMeta[0], 10, 64)
	if err != nil {
		return ver, fmt.Errorf("failed to parse Patch version: %w", err)
	}
	ver.Patch = uint(patch)

	if len(patchMeta) == 2 {
		ver.Meta = patchMeta[1]
	}

	return ver, nil
}

// StringWithAlias returns a string representation of the Version object with the alias.
func (v Version) StringWithAlias() string {
	if v.Alias == "" {
		return v.String()
	}

	return fmt.Sprintf("%s (%s)", v.String(), v.Alias)
}

// String returns a string representation of the Version object.
func (v Version) String() string {
	version := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.Meta != "" {
		version = fmt.Sprintf("%s-%s", version, v.Meta)
	}

	return version
}

// Compare compares the current version (v) with another version (other)
// and returns:
//
//	-1 if v < other
//	 0 if v == other
//	 1 if v > other
func (v Version) Compare(other Version) int {
	if v.Major != other.Major {
		return compareInt(v.Major, other.Major)
	}
	if v.Minor != other.Minor {
		return compareInt(v.Minor, other.Minor)
	}

	return compareInt(v.Patch, other.Patch)
}

func compareInt(a, b uint) int {
	if a < b {
		return -1
	} else if a > b {
		return 1
	}

	return 0
}
