package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	cmd := exec.Command("hugo", "server", "")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	log.Fatal(err)
}
