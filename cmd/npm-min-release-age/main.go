package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/klemjul/scripts/internal/npm"
)

func main() {
	minAgeDays := 7
	if len(os.Args) > 1 {
		if d, err := strconv.Atoi(os.Args[1]); err == nil {
			minAgeDays = d
		}
	}

	minAge := time.Duration(minAgeDays) * 24 * time.Hour
	now := time.Now()

	outdated, err := npm.GetOutdated()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error checking outdated packages: %v\n", err)
		os.Exit(1)
	}

	if len(outdated) == 0 {
		fmt.Println("All packages are up to date.")
		return
	}

	for pkg, info := range outdated {
		versionTimes, err := npm.GetVersionTimes(pkg)
		if err != nil {
			fmt.Printf("⚠️  %s: Could not check versions (%v)\n", pkg, err)
			continue
		}

		latestTime, err := npm.PublishTime(versionTimes, info.Latest)
		if err != nil {
			fmt.Printf("⚠️  %s: Could not find publish time for latest version\n", pkg)
			continue
		}

		age := now.Sub(latestTime)
		daysOld := npm.FormatAge(age)

		if age >= minAge {
			fmt.Printf("✅ %s: current=%s latest=%s (%d days old) SAFE\n", pkg, info.Current, info.Latest, daysOld)
		} else {
			latestSafe := npm.FindLatestSafeVersion(versionTimes, now, minAge)
			if latestSafe != "" {
				safeDays := npm.FormatAge(now.Sub(versionTimes[latestSafe]))
				fmt.Printf("⏳ %s: current=%s latest=%s (%d days old) NOT SAFE | latest safe: %s (%d days old)\n",
					pkg, info.Current, info.Latest, daysOld, latestSafe, safeDays)
			} else {
				fmt.Printf("⏳ %s: current=%s latest=%s (%d days old) NOT SAFE | no safe version found\n",
					pkg, info.Current, info.Latest, daysOld)
			}
		}
	}
}
