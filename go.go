package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func gocmd(goroot string, wd string, command []string) ([]byte, error) {

	cmd := exec.Command("go", command...)
	cmd.Dir = wd
	cmd.Env = []string{fmt.Sprintf("GOROOT=%s", goroot), fmt.Sprintf("PATH=%s/bin", goroot)}

	log.Printf("cmd: %v", cmd)

	output, err := cmd.Output()
	if err != nil {
		log.Print(output)
		panic(err)
	}

	return output, err
}

func gobuild(goroot string, src string, dest string) ([]byte, error) {
	var goCommand = []string{"build", "-o", dest, "."}
	return gocmd(goroot, src, goCommand)
}

func goversion(goroot string) ([]byte, error) {
	var goCommand = []string{"version"}
	wd, _ := os.Getwd()
	return gocmd(goroot, wd, goCommand)
}
