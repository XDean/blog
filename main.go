package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
)

var (
	port       = flag.Int("port", 1313, "Port")
	baseUrl    = flag.String("baseUrl", "http://blog.xdean.com/", "Base URL")
	appendPort = flag.Bool("appendPort", true, "Append port to base url")
	extra      = flag.String("extra", "", "Extra hugo Param")
)

func main() {
	flag.Parse()

	cmd := exec.Command("hugo", "server",
		"-p", strconv.Itoa(*port),
		"--watch", "--disableLiveReload",
		"--appendPort="+strconv.FormatBool(*appendPort),
		"--baseURL", *baseUrl,
		*extra)
	fmt.Println(cmd.Args)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	log.Fatal(err)
}
