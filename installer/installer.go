package installer

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/axellse/rblxutils/common"
	"github.com/axellse/rblxutils/resources"
	"github.com/gen2brain/beeep"
)

func LaunchInstaller() {
	installDir := filepath.Join(common.LocalAppData, "rblxutils")
	if !common.YesNo("Rblxutils does not seem to be installed here, would you like to install rblxutils in " + installDir + "?") {
		return
	}

	err := os.MkdirAll(installDir, 0666)
	if err != nil {
		common.FatalError(err)
	}

	in, err := os.Open(common.BinPath)
	if err != nil {
		common.FatalError(err)
	}
	defer in.Close()

	out, err := os.Create(filepath.Join(installDir, "rblxutils.exe"))
	if err != nil {
		common.FatalError(err)
	}
	defer out.Close()

	n, err := io.Copy(out, in)
	if err != nil {
		common.FatalError(err)
	}

	fmt.Println("copied", n, "bytes")
	err = in.Close()
	if err != nil {
		common.FatalError(err)
	}

	err = out.Close()
	if err != nil {
		common.FatalError(err)
	}

	common.RegisterProtocolHandler()
	HelperInstallFlow()
	beeep.Alert("Rblxutils", "Rblxutils has been installed!", resources.ProgramLogo)

	cmd := exec.Command(filepath.Join(installDir, "rblxutils.exe"))
	err = cmd.Start()
	if err != nil {
		common.FatalError(err)
	}
	fmt.Println("start ok")

	err = cmd.Process.Release()
	if err != nil {
		common.FatalError(err)
	}
	fmt.Println("release ok")
}