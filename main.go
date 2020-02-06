package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"sync"
)

var (
	port        = flag.Int("port", 11074, "Port")
	contentPath = flag.String("content", ".", "hugo content root path")
	commandPort = flag.Int("commandPort", 11075, "Command Port")
	baseUrl     = flag.String("baseUrl", "http://blog.xdean.cn/", "Base URL")
	appendPort  = flag.Bool("appendPort", true, "Append port to base url")
	extra       = flag.String("extra", "", "Extra hugo Param")
	key         = flag.String("key", "", "Key for admin command")

	waitGroup = sync.WaitGroup{}
	adminExit = make(chan int)
)

func main() {
	flag.Parse()
	waitGroup.Add(1)
	go hugo()
	go command()
	waitGroup.Wait()
	os.Exit(0)
}

func hugo() {
	defer waitGroup.Done()
	cmd := exec.Command("hugo", "server",
		"-p", strconv.Itoa(*port),
		"--watch", "--disableLiveReload",
		"--appendPort="+strconv.FormatBool(*appendPort),
		"--baseURL", *baseUrl,
		"--source", *contentPath,
		*extra)
	fmt.Println(cmd.Args)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Print("Hugo process start: ", cmd.Start())

	hugoExit := make(chan error)

	go func(cmd *exec.Cmd) {
		err := cmd.Wait()
		hugoExit <- err
	}(cmd)

	select {
	case <-adminExit:
		fmt.Println("To kill hugo process")
		log.Print(cmd.Process.Kill())
		os.Exit(0)
	case err := <-hugoExit:
		fmt.Println("Hugo exit with:", err.Error())
		os.Exit(1)
	}
}

func command() {
	http.HandleFunc("/admin/pull", func(writer http.ResponseWriter, request *http.Request) {
		if matchKey(request) {
			fmt.Println("To pull new code")
			cmd := exec.Command("git", "pull")
			cmd.Stdout = writer
			cmd.Stderr = writer
			err := cmd.Run()
			if err != nil {
				_, _ = writer.Write([]byte("\n\n"))
				http.Error(writer, err.Error(), http.StatusInternalServerError)
			} else {
				_, _ = writer.Write([]byte("\n\nPull Success\n\nHugo will reload later"))
				writer.WriteHeader(http.StatusOK)
			}
		} else {
			http.Error(writer, "Wrong Key", http.StatusForbidden)
		}
	})

	http.HandleFunc("/admin/exit", func(writer http.ResponseWriter, request *http.Request) {
		if matchKey(request) {
			fmt.Println("To Exit")
			defer func() { adminExit <- 0 }()
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte("Exit"))
		} else {
			http.Error(writer, "Wrong Key", http.StatusForbidden)
		}
	})

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*commandPort), nil))
}

func matchKey(request *http.Request) bool {
	return request.URL.Query().Get("key") == *key
}
