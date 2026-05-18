package bootstrapper

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/axellse/rblxutils/common"
	"github.com/axellse/rblxutils/resources"
)

type Package struct {
	Name                string //eg. RobloxApp.zip
	ExtractionDirectory string
	CompressedSize      int64
	UncompressedSize    int64
}

func GetRawPackageMap() []byte {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get("https://api.axell.me/rblxutils/v1/package-map")
	if err != nil {
		return resources.PackageMap
	}

	ba, err := io.ReadAll(resp.Body)
	if err != nil {
		return resources.PackageMap
	}
	return ba
}

func GetPackages(VersionGUID string) (packages []Package, totalUc int64, totalC int64) {
	resp, err := http.Get("https://setup.rbxcdn.com/" + VersionGUID + "-" + "rbxPkgManifest.txt")
	if err != nil {
		common.FatalError(err)
	}

	ba, err := io.ReadAll(resp.Body)
	if err != nil {
		common.FatalError(err)
	}

	rawPkgMap := GetRawPackageMap()
	packageMap := map[string]string{}
	err = json.Unmarshal(rawPkgMap, &packageMap)
	if err != nil {
		common.FatalError(err)
	}

	packages = []Package{}
	rbxPkgManifest := strings.ReplaceAll(string(ba), "\r", "")
	lines := strings.Split(rbxPkgManifest, "\n")
	for i, line := range lines {
		if !strings.Contains(line, ".") {
			continue
		} //checks if this line is the start of a new package

		if line == "RobloxPlayerInstaller.exe" {
			continue
		}

		compressedSize, err := strconv.ParseInt(lines[i+2], 10, 0)
		if err != nil {
			common.FatalError(err)
		}
		uncompressedSize, err := strconv.ParseInt(lines[i+3], 10, 0)
		if err != nil {
			common.FatalError(err)
		}

		extractionDir, foundDir := packageMap[line]
		if !foundDir {
			common.ErrorStr("Couldn't map directory to package " + line)
		}
		packages = append(packages, Package{
			Name:                line,
			CompressedSize:      compressedSize,
			UncompressedSize:    uncompressedSize,
			ExtractionDirectory: extractionDir,
		})
		totalC += compressedSize
		totalUc += uncompressedSize
	}
	return
}
