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
	}

	server := Exec(filePath, os.Args[2:]...)
	err = server.Start()
	if err != nil {
		filePath = os.Args[1]
		server = Exec(filePath, os.Args[2:]...)
		err2 := server.Start()
		if err2 != nil {
			log.Error(err2)
			return
		}
		return
	}

	log.Info("Start Daemon Server Ok!\n", os.Args[1])

	return
}

func Exec(path string, args ...string) *exec.Cmd {
	server := exec.Command(path, args...)
	server.Stdout = os.Stdout
	server.Stderr = os.Stderr
	server.Stdin = os.Stdin

	var _path string

	if GetRelative(path) >= 0 {
		path = GetPath(os.Args[1])
		server.Dir = _path
	}

	return server
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
