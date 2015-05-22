package main

import (
	"errors"
	"flag"
	"github.com/Congenital/log/v0.2/log"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
	"unsafe"
	//"path/filepath"
)

var Pool = &sync.Pool{}

type StringBuffer struct {
	data   []string
	length int
	Msg    chan int
}

func NewStringBuffer() *StringBuffer {
	return &StringBuffer{
		data:   make([]string, 0),
		length: 0,
		Msg:    make(chan int),
	}
}

func (this *StringBuffer) Write(buff []byte) (int, error) {
	var relative int

	for i := 0; i < len(buff); i++ {
		if buff[i] == '\n' || buff[i] == '\r' {
			buf := buff[relative:i]
			relative = i

			data := *(*string)(unsafe.Pointer(&buf))
			this.data = append(this.data, data)
			this.length++
		}
	}

	return len(buff), nil
}

func (this *StringBuffer) Read() ([]string, error) {
	defer this.Clear()
	if this.data == nil {
		return nil, errors.New("data is nil")
	}

	return this.data, nil
}

func (this *StringBuffer) Clear() {
	for i := 0; i < len(this.data); i++ {
		Pool.Put(this.data[i])
	}

	this.data = make([]string, 0)
	this.length = 0
}

type APP struct {
	User string
	Pid  string

	Program string
	Params  []string
}

func main() {

	if len(os.Args) < 3 {
		log.ErrorNoLine(`
	-root		:	server root path(can empty)
	-filter		:	process filter(if empty, no effect)
	-server		:	server relative path
	-time		:	search time	
	`)
		return
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	var root string
	var server string
	var filter string
	var t int

	flag.StringVar(&root, "root", "", "server root path")
	flag.StringVar(&server, "server", "", "server relative path")
	flag.StringVar(&filter, "filter", "", "process filter")
	flag.IntVar(&t, "time", 5, "search time")

	flag.Parse()

	if len(root) > 0 && root[len(root)-1:len(root)] != "/" {
		root += "/"
	}

	serverlist := strings.Split(server, ",")
	if len(serverlist) <= 0 {
		log.ErrorNoLine("Please Input Server Path")
		return
	}

	servers := make([]string, 0)
	for i := 0; i < len(serverlist); i++ {
		if serverlist[i] != "" {
			servers = append(servers, serverlist[i])
		}
	}

	waitgroup := &sync.WaitGroup{}
	buff := NewStringBuffer()
	timer := time.NewTimer(time.Second * time.Duration(t))

	waitgroup.Add(1)

	go func() {
		defer waitgroup.Done()

		for {
			fl := GetFilter(ParseFilter(filter))

			GetApp(buff, fl)
			buff.Msg <- 1

			<-timer.C
			timer.Reset(time.Second * time.Duration(t))
		}
	}()

	go func() {
		for {
			<-buff.Msg

			data, err := buff.Read()
			if err != nil {
				log.ErrorNoLine(err)
				break
			}

			apps := ParseApp(data)

			noStartServer, err := CheckApp(apps, root, servers)
			if err != nil {
				log.ErrorNoLine(err)
				return
			}

			log.Info("Need Start Server:\n", noStartServer)
			StartApps(root, noStartServer)
		}
	}()

	waitgroup.Wait()
	log.InfoNoLine("Exit")
}

func GetApp(write io.Writer, param string) {
	cmd := exec.Command("/bin/bash", "-c", `ps aux | awk '`+param+` {print $1" "$2" "$11}'`)

	cmd.Stdout = write
	err := cmd.Start()
	if err != nil {
		log.Error(err)
		return
	}
}

func StartApps(root string, app []string) {
	for _, v := range app {
		StartApp(root, v)
	}
}

func StartApp(root string, app string) {
	cmd := exec.Command("run", root+app)

	file, err := os.Create(root + app + ".log")
	if err != nil {
		log.Error(err)
		return
	}

	cmd.Stdout = file
	cmd.Stdin = os.Stdin
	cmd.Stderr = file

	err = cmd.Start()
	if err != nil {
		log.Error(err)
		return
	}

	return
}

func CheckApp(apps []*APP, root string, servers []string) ([]string, error) {
	noStartApp := make([]string, 0)

	appslength := len(apps)

	if len(apps) == 0 {
		return servers, nil
	}

	for _, server := range servers {
		for i, app := range apps {
			if app.Program == root+server {
				break
			}

			if i == appslength-1 {
				noStartApp = append(noStartApp, server)
			}
		}
	}

	return noStartApp, nil
}

func ParseFilter(filter string) []string {
	filters := make([]string, 0)

	fls := strings.Split(filter, ",")
	for i := 0; i < len(fls); i++ {
		if len(fls[i]) != 0 {
			filters = append(filters, fls[i])
		}
	}

	return filters
}

func GetFilter(filters []string) string {
	var filter string

	for i, v := range filters {
		filter += `$11~/` + ParsePath(v) + `/`
		if i != len(filters)-1 {
			filter += "||"
		}
	}

	return filter
}

func ParseApp(data []string) []*APP {
	apps := make([]*APP, 0)

	for i := 0; i < len(data); i++ {
		datas := strings.Split(data[i], " ")
		apps = append(apps, &APP{
			User:    datas[0],
			Pid:     datas[1],
			Program: datas[2],
		})
	}

	return apps
}

func ParsePath(path string) string {
	return strings.Replace(path, "/", "\\/", -1)
}

func UnparsePath(path string) string {
	return strings.Replace(path, "\\/", "/", -1)
}

func GetPath(path string) string {
	return path[:strings.LastIndex(path, "/")]
}
