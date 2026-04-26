package bootstrapper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"time"

	"axell.me/rblxutils/common"
	"axell.me/rblxutils/resources"
)

func LaunchBootstrapper() {
	LaunchUI()

	UiState.CurrentOperation = "Checking for updates."
	Println("checking clientsettingscdn for updates...")
	latestVersion := GetLatestVersion()
	UiState.Progress = 5
	UiState.Update()

	installDir := filepath.Join(common.LPath("./versions"), latestVersion.VersionGUID) 
	if RequiresInstallation(latestVersion.VersionGUID) {
		UiState.CurrentOperation = "Preparing to install " + latestVersion.Version
		Println("update/installation is required, creating directory")

		if _, err := os.Stat(installDir); err == nil {
			UiState.Progress = 10
			Println("deleting old install directory.")
			err := os.RemoveAll(installDir)
			if err != nil {
				common.FatalError(err)
			}
		}

		UiState.Progress = 15
		UiState.Update()
		Println("preparing directories")
		err := os.MkdirAll(installDir, 0666)
		if err != nil {
			common.FatalError(err)
		}

		UiState.CurrentOperation = "Fetching packages"
		Println("fetching packages...")
		pkgs, totalUc, totalC := GetPackages(latestVersion.VersionGUID)
		Println("total uncompressed size is", totalUc, "bytes, total compressed size is", totalC, "bytes")

		UiState.Progress = 20
		Println("fetched", len(pkgs), "packages, now starting to download them")

		for i, pkg := range pkgs {
			UiState.CurrentOperation = "Downloading package " + strconv.Itoa(i + 1) + "/" + strconv.Itoa(len(pkgs))
			Println("downloading package", pkg.Name)
			resp, err := http.Get("https://setup.rbxcdn.com/" + latestVersion.VersionGUID + "-" + pkg.Name)
			if err != nil {
				common.FatalError(err)
			}

			ba, err := io.ReadAll(resp.Body)
			if err != nil {
				common.FatalError(err)
			}

			UiState.CurrentOperation = "Extracting package " + strconv.Itoa(i + 1) + "/" + strconv.Itoa(len(pkgs))
			Println("extracting package", pkg.Name)

			pkgDir := filepath.Join(installDir, pkg.ExtractionDirectory)
			err = os.MkdirAll(pkgDir, 0666)
			if err != nil {
				common.FatalError(err)
			}

			rd := bytes.NewReader(ba)
			err = common.UnzipAt(rd, pkg.CompressedSize, pkgDir)
			if err != nil {
				common.FatalError(err)
			}
		}

		UiState.Progress = 45
		Println("all packages now downloaded, checking if webview2 is installed")

		if !IsWebView2Installed() {
			UiState.CurrentOperation = "Installing WebView2"
			Println("now installing webview2")
			InstallWebView2(installDir)
		} else {
			Println("webview2 already installed.")
		}

		Println("Now writing appsettings.xml") //what even is the purpose of this file
		err = os.WriteFile(filepath.Join(installDir, "AppSettings.xml"), resources.AppSettings, 0666)
		if err != nil {
			common.FatalError(err)
		}

		UiState.Progress = 50
		UiState.CurrentOperation = "Applying mods"
		Println("deleting cache db...")
		DeleteCacheDb()
		UiState.Progress = 60

		Println("roblox now installed, moving onto mods.")
		ApplyFileMods(installDir)
		UiState.Progress = 70
		Println("mods are now installed")

		common.LoadState()
		common.State.RequiresModApplication = false
		err = common.WriteState()
		if err != nil {return}
	}

	UiState.Progress = 80
	UiState.CurrentOperation = "Preparing Rblxutils for launch"
	Println("install ok, now doing final pre-launch procedures.")
	RunPreLaunchProcedures()

	UiState.CurrentOperation = "Launching client"
	UiState.Progress = 90
	Println("everything done, now launching client...")

	if !slices.Contains(common.Config.Misc.DebugOptions, "skip-launch") {
		cmd := exec.Command(filepath.Join(installDir, "RobloxPlayerBeta.exe"), os.Args[1:]...)
		err := cmd.Start()
		if err != nil {
			common.Error(err)
		}
	}


	UiState.CurrentOperation = "Roblox is now running"
	UiState.Progress = 100
	Println("roblox is now running, closing bootstrapper window in 5 seconds.")
	time.Sleep(5 * time.Second)
	UiState.CloseWindow()

	select {}
}

type clientsettingscdnResponse struct {
	Version string `json:"version"`
	VersionGUID string `json:"clientVersionUpload"`
}

func GetLatestVersion() clientsettingscdnResponse {
	resp, err := http.Get("https://clientsettingscdn.roblox.com/v2/client-version/WindowsPlayer")
	if err != nil {
		common.FatalErrorStr("Cannot check for updates, do you have an internet connection?")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		common.FatalError(err)
	}

	versionResponse := clientsettingscdnResponse{}
	err = json.Unmarshal(body, &versionResponse)
	if err != nil {
		common.FatalError(err)
	}
	return versionResponse
}

var BootstrapperLog = bytes.Buffer{}
func Println(args ...any) {
	s := fmt.Sprintln(args...)
	BootstrapperLog.Write([]byte(s))
	UiState.Update()
	fmt.Print(s)
}

func RequiresInstallation(latestVersion string) bool {
	Println("inspecting roblox installations...")
	versions, err := os.ReadDir(common.LPath("./versions"));
	if err != nil {
		return true
	}

	Println("latest version guid is", latestVersion)
	foundLatest := false
	for _, v := range versions {
		if v.Name() == latestVersion && v.IsDir() {
			foundLatest = true
		} /*else {
			err := os.Remove(filepath.Join(common.LPath("./versions"), v.Name()))
			if err != nil {
				common.ErrorStr("Could not remove junk!")
			}
		}*/
	}
	if !foundLatest {
		return true
	}

	Println("checking if mods need to be re-applied...")
	return common.State.RequiresModApplication
}
