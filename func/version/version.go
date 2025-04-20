package version

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const VersionControl = "1.2"

type GLOBAL_CNC_VERSIONS struct {
	LastVersion        json.Number `json:"last_version"`
	StableVersion      json.Number `json:"stable_version"`
	UnavailableVersion json.Number `json:"unavailable_version"`
	Download           string      `json:"download"`
}

func CompareVersions(currentVersion, latestVersion string) bool {
	partsCurrent := strings.Split(currentVersion, ".")
	partsLatest := strings.Split(latestVersion, ".")

	if len(partsCurrent) < 2 || len(partsLatest) < 2 {
		fmt.Println("Invalid version format")
		return false
	}

	majorCurrent, _ := strconv.Atoi(partsCurrent[0])
	minorCurrent, _ := strconv.Atoi(partsCurrent[1])
	majorLatest, _ := strconv.Atoi(partsLatest[0])
	minorLatest, _ := strconv.Atoi(partsLatest[1])

	return majorCurrent < majorLatest || (majorCurrent == majorLatest && minorCurrent < minorLatest)
}

func CheckVersion() (string, string, string, string, error) {
	url := "https://raw.githubusercontent.com/Royaloakap/Version/main/versions.json"
	resp, err := http.Get(url)
	if err != nil {
		return "", "", "", "", fmt.Errorf("\u001B[38;5;196mfailed\u001B[38;5;230m to check for Royal GRADIENT version: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", "", "", fmt.Errorf("\u001B[38;5;196mfailed\u001B[38;5;230m to check for Royal GRADIENT version. Received status code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", "", fmt.Errorf("\u001B[38;5;196mfailed\u001B[38;5;230m to read response body: %w", err)
	}

	var cncVersions GLOBAL_CNC_VERSIONS
	if err := json.Unmarshal(body, &cncVersions); err != nil {
		return "", "", "", "", fmt.Errorf("\u001B[38;5;196mfailed\u001B[38;5;230m to unmarshal response body: %w", err)
	}

	lastVersion := cncVersions.LastVersion.String()
	stableVersion := cncVersions.StableVersion.String()
	unavailableVersion := cncVersions.UnavailableVersion.String()

	return cncVersions.Download, lastVersion, stableVersion, unavailableVersion, nil
}
