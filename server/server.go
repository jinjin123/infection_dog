package server

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"infection/request"
	"io/ioutil"
	"math"
	"net/http"
	"sort"
	"strconv"
)

// Server information
type Server struct {
	URL      string `xml:"url,attr"`
	Lat      string `xml:"lat,attr"`
	Lon      string `xml:"lon,attr"`
	Name     string `xml:"name,attr"`
	Country  string `xml:"country,attr"`
	Sponsor  string `xml:"sponsor,attr"`
	ID       string `xml:"id,attr"`
	URL2     string `xml:"url2,attr"`
	Host     string `xml:"host,attr"`
	Distance float64
	DLSpeed  float64
	ULSpeed  float64
}

// ServerList : List of Server
type ServerList struct {
	Servers []Server `xml:"servers>server"`
}

// Servers : For sorting servers.
type Servers []Server

// ByDistance : For sorting servers.
type ByDistance struct {
	Servers
}

// Len : length of servers. For sorting servers.
func (s Servers) Len() int {
	return len(s)
}

// Swap : swap i-th and j-th. For sorting servers.
func (s Servers) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less : compare the distance. For sorting servers.
func (b ByDistance) Less(i, j int) bool {
	return b.Servers[i].Distance < b.Servers[j].Distance
}

func FetchServerList(lat string, Lon string) ServerList {
	// Fetch xml server data
	resp, _ := http.Get("http://www.speedtest.net/speedtest-servers-static.php")
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if len(body) == 0 {
		resp, _ = http.Get("http://c.speedtest.net/speedtest-servers-static.php")
		body, _ = ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
	}

	// Decode xml
	decoder := xml.NewDecoder(bytes.NewReader(body))
	list := ServerList{}
	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			decoder.DecodeElement(&list, &se)
		}
	}

	// Calculate distance
	for i := range list.Servers {
		server := &list.Servers[i]
		sLat, _ := strconv.ParseFloat(server.Lat, 64)
		sLon, _ := strconv.ParseFloat(server.Lon, 64)
		uLat, _ := strconv.ParseFloat(lat, 64)
		uLon, _ := strconv.ParseFloat(Lon, 64)
		server.Distance = distance(sLat, sLon, uLat, uLon)
	}

	// Sort by distance
	sort.Sort(ByDistance{list.Servers})

	return list
}

func distance(lat1 float64, lon1 float64, lat2 float64, lon2 float64) float64 {
	radius := 6378.137

	a1 := lat1 * math.Pi / 180.0
	b1 := lon1 * math.Pi / 180.0
	a2 := lat2 * math.Pi / 180.0
	b2 := lon2 * math.Pi / 180.0

	x := math.Sin(a1)*math.Sin(a2) + math.Cos(a1)*math.Cos(a2)*math.Cos(b2-b1)
	return radius * math.Acos(x)
}

// FindServer : find server by serverID
func (l *ServerList) FindServer(serverID []int) Servers {
	servers := Servers{}

	for _, sid := range serverID {
		for _, s := range l.Servers {
			id, _ := strconv.Atoi(s.ID)
			if sid == id {
				servers = append(servers, s)
			}
		}
	}

	if len(servers) == 0 {
		servers = append(servers, l.Servers[0])
	}

	return servers
}

// Show : show server list
func (l ServerList) Show() {
	for _, s := range l.Servers {
		fmt.Printf("[%4s] %8.2fkm ", s.ID, s.Distance)
		fmt.Printf(s.Name + " (" + s.Country + ") by " + s.Sponsor + "\n")
	}
}

// Show : show server information
func (s Server) Show() {
	fmt.Printf(" \n")
	fmt.Printf("Target Server: [%4s] %8.2fkm ", s.ID, s.Distance)
	fmt.Printf(s.Name + " (" + s.Country + ") by " + s.Sponsor + "\n")
}

// StartTest : start testing to the servers.
func (svrs Servers) StartTest() {
	for i, s := range svrs {
		s.Show()
		latency := request.PingTest(s.URL)
		dlSpeed := request.DownloadTest(s.URL, latency)
		ulSpeed := request.UploadTest(s.URL, latency)
		svrs[i].DLSpeed = dlSpeed
		svrs[i].ULSpeed = ulSpeed
	}
}

type Speed struct {
	Down string
	Up   string
}

// ShowResult : show testing result
func (svrs Servers) ShowResult() *Speed {
	if len(svrs) == 1 {
		speed := Speed{
			Down: fmt.Sprintf("%5.2f", svrs[0].DLSpeed),
			Up:   fmt.Sprintf("%5.2f", svrs[0].ULSpeed),
		}
		return &speed
		//fmt.Printf("Download: %5.2f Mbit/s,Upload: %5.2f Mbit/s", svrs[0].DLSpeed,svrs[0].ULSpeed)

	} else {
		for _, s := range svrs {
			//speed,_ :=fmt.Printf("[%4s] Download: %5.2f Mbit/s, Upload: %5.2f Mbit/s", s.ID, s.DLSpeed, s.ULSpeed)
			speed := Speed{
				Down: fmt.Sprintf("%5.2f", s.DLSpeed),
				Up:   fmt.Sprintf("%5.2f", s.ULSpeed),
			}
			return &speed
		}
		avgDL := 0.0
		avgUL := 0.0
		for _, s := range svrs {
			avgDL = avgDL + s.DLSpeed
			avgUL = avgUL + s.ULSpeed
		}
		//speed,_ := fmt.Printf("Download Avg: %5.2f Mbit/s,Upload Avg: %5.2f Mbit/s", avgDL/float64(len(svrs)), avgUL/float64(len(svrs)))
		speed := Speed{
			Down: fmt.Sprintf("%5.2f", avgDL/float64(len(svrs))),
			Up:   fmt.Sprintf("%5.2f", avgUL/float64(len(svrs))),
		}
		return &speed
	}
}

func (svrs Servers) checkResult() bool {
	errFlg := false
	if len(svrs) == 1 {
		s := svrs[0]
		errFlg = (s.DLSpeed*100 < s.ULSpeed) || (s.DLSpeed > s.ULSpeed*100)
	} else {
		for _, s := range svrs {
			errFlg = errFlg || (s.DLSpeed*100 < s.ULSpeed) || (s.DLSpeed > s.ULSpeed*100)
		}
	}
	return errFlg
}
