package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

type GoConfig struct {
	GOOS   string
	GOARCH string
	GOROOT string
}

func gocmd(config GoConfig, wd string, command []string) ([]byte, error) {

	cmd := exec.Command("go", command...)
	cmd.Dir = wd
	cmd.Env = []string{
		"CGO_ENABLED=0",
		fmt.Sprintf("GOOS=%s", config.GOOS),
		fmt.Sprintf("GOARCH=%s", config.GOARCH),
		fmt.Sprintf("GOROOT=%s", config.GOROOT),
		fmt.Sprintf("PATH=%s/bin", config.GOROOT),
	}

	log.Printf("cmd: %v", cmd)

	output, err := cmd.Output()
	if err != nil {
		log.Print(output)
		panic(err)
	}

	return output, err
}

func gobuild(config GoConfig, src string, dest string) ([]byte, error) {
	var goCommand = []string{"build", "-o", dest, "."}
	return gocmd(config, src, goCommand)
}

func goversion(config GoConfig) ([]byte, error) {
	var goCommand = []string{"version"}
	wd, _ := os.Getwd()
	return gocmd(config, wd, goCommand)
}
