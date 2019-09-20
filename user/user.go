package user

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
)

// User information
type User struct {
	IP  string `xml:"ip,attr"`
	Lat string `xml:"lat,attr"`
	Lon string `xml:"lon,attr"`
	Isp string `xml:"isp,attr"`
}

// Users : for decode xml
type Users struct {
	Users []User `xml:"client"`
}

func FetchUserInfo() User {
	// Fetch xml user data
	resp, _ := http.Get("http://speedtest.net/speedtest-config.php")
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	// Decode xml
	decoder := xml.NewDecoder(bytes.NewReader(body))
	users := Users{}
	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			decoder.DecodeElement(&users, &se)
		}
	}
	if users.Users == nil {
		log.Println("Warning: Cannot fetch user information. http://www.speedtest.net/speedtest-config.php is temporarily unavailable.")
		return User{}
	}
	return users.Users[0]
}

type Outboud struct {
	OIp string `json:"ip"`
	Isp string `json:"isp"`
	Lat string `json:"lat"`
	Lon string `json:"Lon"`
}

// Show user location
func (u *User) Show() *Outboud {
	if u.IP != "" {
		out := Outboud{
			OIp: u.IP,
			Isp: u.Isp,
			Lat: u.Lat,
			Lon: u.Lon,
		}
		//outboud,_ :=log.Println("Testing From IP: " + u.IP + " (" + u.Isp + ") [" + u.Lat + ", " + u.Lon + "]")
		return &out
	}
	return &Outboud{}
}
