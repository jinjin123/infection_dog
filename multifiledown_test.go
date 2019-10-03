package main

import (
	"infection/util/lib"
	"log"
	"testing"
)

func MultiFileDown(files []string, step string) {
	if len(files) == 0 && step == "init" {
		var fileinit = []struct {
			Name string
		}{
			{"WindowsDaemon.exe"},
			{"sqlite3_386.dll"},
			{"sqlite3_amd64.dll"},
		}
		for _, name := range fileinit {
			log.Print(lib.MIDFILE+name.Name, name.Name)
			//Get(MIDFILE+name.Name,name.Name)
		}
	}
}
func TestMultiFileDown(t *testing.T) {
	MultiFileDown([]string{}, "init")
}
