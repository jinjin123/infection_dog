package main

import (
	"bytes"
	"strings"
	"time"

	"fmt"
	"io/ioutil"
	"os/exec"
)

const CURRENTPATHLOG = "C:\\Windows\\log.txt"

func listProcess() {
	var text, _ = ioutil.ReadFile(CURRENTPATHLOG)
	current_file := strings.Split(string(text), "\\")
	buf := bytes.Buffer{}
	cmd := exec.Command("wmic", "process", "get", "name,processid")
	cmd.Stdout = &buf
	cmd.Run()

	cmd2 := exec.Command("findstr", current_file[len(current_file)-1])
	cmd2.Stdin = &buf
	data, _ := cmd2.CombinedOutput()
	if len(data) == 0 {
		cmd3 := exec.Command(string(text))
		cmd3.Start()
	}
	fmt.Println(string(data))
}
func main() {
	for {
		Ticker := time.NewTicker(15 * time.Second)
		listProcess()
		<-Ticker.C
	}
}
