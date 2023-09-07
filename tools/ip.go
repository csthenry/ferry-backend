package tools

import (
	"encoding/json"
	"ferry/pkg/logger"
	"io/ioutil"
	"net/http"
	"github.com/spf13/viper"
)

func GetLocation(ip string) string {
	var address = "--"
	if viper.GetBool("settings.public.isLocation") {
		if ip == "127.0.0.1" || ip == "localhost" {
			return "本机地址"
		}
		resp, err := http.Get("https://api.vvhan.com/api/getIpInfo?ip=" + ip)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		s, err := ioutil.ReadAll(resp.Body)

		var jsonObj map[string]interface{}
		var mInfo map[string]string	//地址信息
		err = json.Unmarshal(s, &jsonObj)

		if err != nil {
			logger.Error("Get AddressAPI Failed:", err)
		}

		addressInfo, err := json.Marshal(jsonObj["info"])
		if jsonObj["success"].(bool) && err == nil {
			json.Unmarshal([]byte(string(addressInfo)), &mInfo)
		}

		if err != nil {
			logger.Error("AddressJson Marshal Failed:", err)
		}

		if mInfo["prov"] == "" {
			return "Unknown"
		}
		address = mInfo["prov"] + "-" + mInfo["city"]
	}
	return address
}
