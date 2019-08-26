package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"strconv"
)

var (
	port    = flag.Int("port", 1313, "Port")
	baseUrl = flag.String("baseUrl", "http://blog.xdean.com/", "Base URL")
)

func main() {
	flag.Parse()

	cmd := exec.Command("hugo", "server",
		"-p", strconv.Itoa(*port),
		"--watch", "--disableLiveReload",
		"--baseURL", *baseUrl)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	log.Fatal(err)
}
