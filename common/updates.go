package common

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	_ "embed"
	"encoding/pem"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"

	"github.com/axellse/rblxutils/resources"
)

type UpdateResponse struct {
	Title    string
	Text     string
	TextCols int
	Upgrade  struct {
		Version   int
		Download  string
		Signature []byte
	}
}

// Checks and performs any updates, then returns Uresp
func CheckForUpdates() UpdateResponse {
	var uResp UpdateResponse
	QueryAndUnmarshal("https://api.axell.me/rblxutils/v1/updates/json", &uResp)

	i, err := strconv.Atoi(resources.Version)
	if err != nil {
		FatalError(err)
	}

	if uResp.Upgrade.Version <= i || !YesNo("A new update is available ("+strconv.Itoa(uResp.Upgrade.Version)+"). Would you like to upgrade?") {
		return uResp
	}

	PerformUpdate(uResp.Upgrade.Download, uResp.Upgrade.Signature)
	return uResp
}

func PerformUpdate(link string, signature []byte) {
	u, err := url.Parse(link)
	if err != nil {
		FatalError(err)
	}

	if u.Host != "github.com" && u.Host != "api.axell.me" {
		if !YesNo("Rblxutils didn't recognize the url '" + u.Host + "'. Do you still want to continue with the upgrade?") {
			return
		}
	}

	resp, err := http.Get(link)
	if err != nil {
		FatalError(err)
	}

	ba, err := io.ReadAll(resp.Body)
	if err != nil {
		FatalError(err)
	}

	err = os.WriteFile(LPath("./rblxutils_fresh.exe"), ba, 0666)
	if err != nil {
		FatalError(err)
	}

	block, _ := pem.Decode(resources.UpdatePublicKey)
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		FatalError(err)
	}

	publicKey := pub.(*rsa.PublicKey)

	sum := sha256.Sum256(ba)
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, sum[:], signature)
	if err != nil {
		FatalErrorStr("Could not verify update: " + err.Error())
	} else {
		Notification("Verified signature of update package.")
	}

	cmd := exec.Command("powershell", "-WindowStyle", "Hidden", "-Command", "Start-Sleep -Seconds 3; Remove-Item .\\rblxutils.exe -Force; Rename-Item .\\rblxutils_fresh.exe rblxutils.exe")
	cmd.Dir = DotSlash
	err = cmd.Start()
	if err != nil {
		FatalError(err)
	}

	cmd.Process.Release()

	os.Exit(0)
}
