package tools

import (
	"bytes"
	"encoding/json"
	"ferry/pkg/logger"
	"io/ioutil"
	"net/http"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"github.com/spf13/viper"
)

func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, err := ioutil.ReadAll(reader)
	return d, err
}

func GetLocation(ip string) string {
	var address = "--"
	type IpAddress struct {
		Provence  string    `json:"pro"`
		City      string    `json:"city"`
	}
	if viper.GetBool("settings.public.isLocation") {
		if ip == "127.0.0.1" || ip == "localhost" {
			return "本机地址"
		}
		resp, err := http.Get("https://whois.pconline.com.cn/ipJson.jsp?ip=" + ip + "&json=true")
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		raw, err := ioutil.ReadAll(resp.Body)
		s, err := GbkToUtf8(raw)

		var mInfo map[string]string
		json.Unmarshal(s, &mInfo)

		var ipAddress IpAddress
		ipAddress.Provence = mInfo["pro"]
		ipAddress.City = mInfo["city"]

		if err != nil {
			logger.Error("Get IpAddress Failed:", err)
		}

		if ipAddress.Provence == "" {
			return "未知"
		}
		address = ipAddress.Provence + "-" + ipAddress.City
	}
	return address
}
