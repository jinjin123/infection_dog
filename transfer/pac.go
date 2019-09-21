package transfer

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
)

var setWindowsPACRegistry func() = func() {}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, Data[0])
	switch runtime.GOOS {
	case "darwin":
	case "windows":
		setWindowsPACRegistry()
	case "linux":
	}
}
func PacHandle(PacPort string) {
	http.HandleFunc("/pac", handler)
	if err := http.ListenAndServe(PacPort, nil); err != nil {
		log.Println("pac Faild", err)
	}
}
