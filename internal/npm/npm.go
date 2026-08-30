package npm

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"sort"
	"time"
)

// OutdatedInfo represents the output of `npm outdated --json` for a single package.
type OutdatedInfo struct {
	Current string `json:"current"`
	Wanted  string `json:"wanted"`
	Latest  string `json:"latest"`
	Depends string `json:"depends"`
	Dev     bool   `json:"dev"`
}

type versionTime struct {
	version string
	t       time.Time
}

// GetOutdated runs `npm outdated --json` and returns a map of package names to info.
func GetOutdated() (map[string]OutdatedInfo, error) {
	cmd := exec.Command("npm", "outdated", "--json")
	output, err := cmd.CombinedOutput()
	if err != nil {
		// npm outdated exits with code 1 when there are outdated packages,
		// but still writes valid JSON to stdout.
		if _, ok := err.(*exec.ExitError); !ok {
			return nil, err
		}
	}

	var result map[string]OutdatedInfo
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetVersionTimes fetches the publish timestamps for every version of a package.
func GetVersionTimes(pkg string) (map[string]time.Time, error) {
	cmd := exec.Command("npm", "view", pkg, "time", "--json")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var timeStrs map[string]string
	if err := json.Unmarshal(output, &timeStrs); err != nil {
		return nil, err
	}

	result := make(map[string]time.Time)
	for ver, ts := range timeStrs {
		// npm time object includes metadata keys like "created" and "modified"
		if ver == "created" || ver == "modified" {
			continue
		}
		t, err := time.Parse(time.RFC3339, ts)
		if err != nil {
			continue
		}
		result[ver] = t
	}

	return result, nil
}

// FindLatestSafeVersion scans the version history and returns the most recent
// version whose age is greater than or equal to minAge.
func FindLatestSafeVersion(versionTimes map[string]time.Time, now time.Time, minAge time.Duration) string {
	var versions []versionTime
	for v, t := range versionTimes {
		versions = append(versions, versionTime{version: v, t: t})
	}

	// Sort by publish time descending (newest first)
	sort.Slice(versions, func(i, j int) bool {
		return versions[i].t.After(versions[j].t)
	})

	for _, v := range versions {
		if now.Sub(v.t) >= minAge {
			return v.version
		}
	}
	return ""
}

// FormatAge formats a duration as "X days old".
func FormatAge(age time.Duration) int {
	return int(age.Hours() / 24)
}

// PublishTime returns the publish time for a specific version.
func PublishTime(versionTimes map[string]time.Time, version string) (time.Time, error) {
	t, ok := versionTimes[version]
	if !ok {
		return time.Time{}, fmt.Errorf("version %s not found in time data", version)
	}
	return t, nil
}
