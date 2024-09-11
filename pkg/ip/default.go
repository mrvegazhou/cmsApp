package ip

import (
	"errors"
	"github.com/oschwald/geoip2-golang"
	"net"
	"os"
	"path/filepath"
)

type IPInfo struct {
	IP       string   `json:"ip"`       //ip地址
	Country  string   `json:"country"`  //国家
	Province string   `json:"province"` //省份
	City     string   `json:"city"`     //城市
	Location struct { //经纬度
		AccuracyRadius uint16  `json:"accuracy_radius"`
		Latitude       float64 `json:"latitude"`
		Longitude      float64 `json:"longitude"`
		MetroCode      uint    `json:"metro_code"`
		TimeZone       string  `json:"time_zone"`
	} `json:"location"`
}

var localNetworkNames = map[string]string{
	"zh-CN": "局域网",
	"en":    "local network",
}

//var mmdb embed.FS

func init() {

	//fs, err := mmdb.Open(dbPath)
	//if err != nil {
	//	panic(err)
	//}
	//data, err := io.ReadAll(fs)
	//if err != nil {
	//	panic(err)
	//}
	//db, err = geoip2.FromBytes(data)

}

func IsPublicIP(IP net.IP) bool {
	if IP.IsLoopback() || IP.IsLinkLocalMulticast() || IP.IsLinkLocalUnicast() {
		return false
	}
	if ip4 := IP.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false

		default:
			return true
		}
	}
	return false
}

func HandleIPInfo(ipStr string, language string) (*IPInfo, error) {
	if language != "cn" && language != "en" {
		language = "cn"
	}
	if language == "cn" {
		language = "zh-CN"
	}
	ipInfo := &IPInfo{}
	if ipStr == "" {
		return ipInfo, errors.New("ip is error")
	}
	var err error
	dir, err := os.Getwd()
	dbPath := filepath.Join(dir, "pkg", "ip", "GeoLite2-City.mmdb")
	db, err := geoip2.Open(dbPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	ip := net.ParseIP(ipStr)
	if !IsPublicIP(ip) {
		ipInfo = &IPInfo{
			IP:   ip.String(),
			City: localNetworkNames[language],
		}
		return ipInfo, nil
	}
	city, err := db.City(ip)
	if err != nil {
		return ipInfo, err
	}
	ipInfo.IP = ip.String()
	ipInfo.Country = city.Country.Names[language]
	if len(city.Subdivisions) > 0 {
		ipInfo.Province = city.Subdivisions[0].Names[language]
	} else {
		ipInfo.Province = city.City.Names[language]
	}
	ipInfo.City = city.City.Names[language]
	ipInfo.Location.AccuracyRadius = city.Location.AccuracyRadius
	ipInfo.Location.Latitude = city.Location.Latitude
	ipInfo.Location.Longitude = city.Location.Longitude
	ipInfo.Location.MetroCode = city.Location.MetroCode
	ipInfo.Location.TimeZone = city.Location.TimeZone
	return ipInfo, nil
}
