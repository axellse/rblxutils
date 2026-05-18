package bootstrapper

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/axellse/rblxutils/common"
	"golang.org/x/sys/windows/registry"
)

func IsWebView2Installed() bool {
	k1, err1 := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\EdgeUpdate\Clients\{F3017226-FE2A-4295-8BDF-00C3A9A7E4C5}`, registry.QUERY_VALUE)
	k2, err2 := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\WOW6432Node\Microsoft\EdgeUpdate\Clients\{F3017226-FE2A-4295-8BDF-00C3A9A7E4C5}`, registry.QUERY_VALUE)
	if err1 != nil && err2 != nil {
		fmt.Println("couldnt open either reg key")
		return false
	}

	val1, _, err1 := k1.GetStringValue("pv")
	val2, _, err2 := k2.GetStringValue("pv")
	if err1 != nil && err2 != nil {
		fmt.Println("couldnt read either value")
		return false
	}

	if (val1 == "" || val1 == "0.0.0.0") && (val2 == "" || val2 == "0.0.0.0") {
		fmt.Println("bad value")
		return false
	}
	return true
}

func InstallWebView2(installDir string) {
	installerDir := filepath.Join(installDir, "WebView2RuntimeInstaller", "MicrosoftEdgeWebview2Setup.exe")
	cmd := exec.Command(installerDir, "/silent", "/install")
	ba, err := cmd.CombinedOutput()
	fmt.Println(string(ba))
	if err != nil {
		common.FatalError(err)
	}
}
