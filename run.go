package main

import (
	"github.com/Congenital/log/v0.2/log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {

	if len(os.Args) < 2 {
		log.ErrorNoLine(`
	Please Input Server Name And Params:
	Absolute path
		`)
		return
	}

	filePath, err := filepath.Abs(os.Args[1])
	if err != nil {
		log.Error(err)
		return
	}

	path := GetPath(os.Args[1])

	if GetRelative(path) >= 0 {
		if err != nil {
			log.Error(err)
			return
		}
	}

	server := exec.Command(filePath, os.Args[2:]...)

	server.Stdout = os.Stdout
	server.Stderr = os.Stderr
	server.Stdin = os.Stdin
	server.Dir = path

	err = server.Start()
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Start Daemon Server Ok!\n", os.Args[1])

	return
}

func GetPath(path string) string {
	relative := strings.LastIndex(path, "/")
	if relative <= 0 {
		return path
	}

	return path[:strings.LastIndex(path, "/")]
}

func GetRelative(path string) int {
	return strings.LastIndex(path, "/")
}
