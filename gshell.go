package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/gobuffalo/packr"
)

var (
	supportedPlatforms = []string{"windows", "linux", "darwin"}
)

func randomID(n int) string {
	bytes := make([]byte, n)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func main() {
	log.Printf("extracting assets ...")
	assetsBox := packr.NewBox("./assets")
	gshellDir := setup(assetsBox)
	defer func() {
		log.Printf("cleaning up assets in: %s", gshellDir)
		// os.RemoveAll(gshellDir)
	}()

	log.Printf("extracted assets to: %s", gshellDir)
	goroot := fmt.Sprintf("%s/go", gshellDir)
	output, _ := goversion(goroot)
	log.Printf("embedded %s", string(output))

	shellID := randomID(8)
	log.Printf("creating shell with ID: %s", shellID)

	shellGo, _ := assetsBox.MustString("shell/shell.go")
	shellDir := fmt.Sprintf("%s/shells/%s", gshellDir, shellID)
	err := os.MkdirAll(shellDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	log.Printf("shell dir: %s", shellDir)
	ioutil.WriteFile(fmt.Sprintf("%s/shell.go", shellDir), []byte(shellGo), 0600)
	for _, platform := range supportedPlatforms {
		shellPlatform, _ := assetsBox.MustBytes(fmt.Sprintf("shell/shell_%s.go", platform))
		ioutil.WriteFile(fmt.Sprintf("%s/shell_%s.go", shellDir, platform), shellPlatform, 0600)
	}

	cwd, _ := os.Getwd()
	output, _ = gobuild(goroot, shellDir, fmt.Sprintf("%s/implant.exe", cwd))
	log.Printf(" --- build:\n%s\n", string(output))
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}
