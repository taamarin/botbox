package module

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ipApiResp struct {
	Query        string `json:"query"`
	Country      string `json:"country"`
	CountryCode  string `json:"countryCode"`
	Region       string `json:"regionName"`
	City         string `json:"city"`
	ISP          string `json:"isp"`
	Org          string `json:"org"`
	AS           string `json:"as"`
	Timezone     string `json:"timezone"`
	Lat          float64 `json:"lat"`
	Lon          float64 `json:"lon"`
}

// HandleMyIP shows local + public IP with ISP/location info.
func HandleMyIP(bot *tgbotapi.BotAPI, chatID int64) {
	// Kirim pesan loading
	loadingMsg, _ := bot.Send(tgbotapi.NewMessage(chatID, "⏳ Loading..."))

	// 1) Local IPs
	var localInfo []string
	ifs, err := net.Interfaces()
	if err == nil {
		for _, iface := range ifs {
			if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
				continue
			}
			addrs, _ := iface.Addrs()
			var ips []string
			for _, a := range addrs {
				ipStr := a.String()
				if strings.Contains(ipStr, "/") {
					ipStr = strings.SplitN(ipStr, "/", 2)[0]
				}
				if ipStr != "" {
					ips = append(ips, ipStr)
				}
			}
			if len(ips) > 0 {
				localInfo = append(localInfo, fmt.Sprintf("%s: %s", iface.Name, strings.Join(ips, ", ")))
			}
		}
	}

	// 2) Public IP + info via ip-api
	client := &http.Client{Timeout: 8 * time.Second}
	var ipData ipApiResp
	resp, err := client.Get("http://ip-api.com/json/")
	if err == nil {
		defer resp.Body.Close()
		json.NewDecoder(resp.Body).Decode(&ipData)
	}

	publicText := "(could not fetch public IP)"
	if ipData.Query != "" {
		publicText = fmt.Sprintf(
			"%s\nISP: %s\nOrg: %s\nASN: %s\nLoc: %s, %s (%s)\nTimezone: %s",
			ipData.Query,
			ipData.ISP,
			ipData.Org,
			ipData.AS,
			ipData.City,
			ipData.Country,
			ipData.CountryCode,
			ipData.Timezone,
		)
	}

	// 3) Build result
	result := "📡 *My IP Info*\n\n"
	result += "*Public IP Info:*\n```\n" + publicText + "\n```\n"

	if len(result) > 4000 {
		result = result[:3900] + "\n...(truncated)"
	}

	// Edit pesan loading jadi hasil
	edit := tgbotapi.NewEditMessageText(chatID, loadingMsg.MessageID, result)
	edit.ParseMode = "Markdown"
	bot.Send(edit)
}