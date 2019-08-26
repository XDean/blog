package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
)

var (
	port = flag.Int("port", 11073, "Port")
)

func main() {
	cmd := exec.Command("hugo", "server", "-p", string(*port), "--watch", "--disableLiveReload")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	log.Fatal(err)
}
