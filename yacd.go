package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// const (
	// MIHOMO_API = "http://192.168.1.1:9090"
	// API_SECRET = "123456"
// )

var HEADERS = map[string]string{
	"Authorization": "Bearer " + apiSecret,
}

type ProxyItem struct {
	Type string   `json:"type"`
	Now  string   `json:"now"`
	All  []string `json:"all"`
}

type ProxiesResponse struct {
	Proxies map[string]ProxyItem `json:"proxies"`
}

// helper GET
func getJSON(url string, target interface{}) error {
	req, _ := http.NewRequest("GET", url, nil)
	for k, v := range HEADERS {
		req.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(target)
}

// helper PUT/POST
func sendRequest(method, url string, payload []byte) (*http.Response, error) {
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(payload))
	for k, v := range HEADERS {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", "application/json")
	return http.DefaultClient.Do(req)
}

// keyboard utama
func yacdMenu(proxies map[string]ProxyItem) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton

	for name, item := range proxies {
		if item.Type == "Selector" {
			btn := tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("%s ➠ %s", name, item.Now),
				"select_"+name,
			)
			row = append(row, btn)
			if len(row) == 2 {
				rows = append(rows, row)
				row = []tgbotapi.InlineKeyboardButton{}
			}
		}
	}
	if len(row) > 0 {
		rows = append(rows, row)
	}

	rows = append(rows,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Status Delay", "status"),
			tgbotapi.NewInlineKeyboardButtonData("Reload Config", "reload"),
			tgbotapi.NewInlineKeyboardButtonData("Restart", "restart"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Upgrade", "upgrade"),
			tgbotapi.NewInlineKeyboardButtonData("Versi", "version"),
		),
	)

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// handler command /yacd
func handleYacd(bot *tgbotapi.BotAPI, chatID int64) {
	var proxiesResp ProxiesResponse
	err := getJSON(mihomoAPI+"/proxies", &proxiesResp)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "❌ Gagal ambil data proxies: "+err.Error()))
		return
	}

	msg := tgbotapi.NewMessage(chatID, "Pilih grup proxy atau perintah:")
	msg.ReplyMarkup = yacdMenu(proxiesResp.Proxies)
	bot.Send(msg)
}

// handler callback dashboard
func handleYacdCallback(bot *tgbotapi.BotAPI, cq *tgbotapi.CallbackQuery) {
	data := cq.Data
	chatID := cq.Message.Chat.ID
	msgID := cq.Message.MessageID

	// tampilkan loading
	edit := tgbotapi.NewEditMessageText(chatID, msgID, "⏳ Loading...")
	bot.Send(edit)

	switch {
	case strings.HasPrefix(data, "select_"):
		group := strings.TrimPrefix(data, "select_")
		var groupInfo ProxyItem
		if err := getJSON(mihomoAPI+"/proxies/"+group, &groupInfo); err != nil {
			bot.Send(tgbotapi.NewEditMessageText(chatID, msgID, "❌ Gagal ambil data group: "+err.Error()))
			return
		}

		var rows [][]tgbotapi.InlineKeyboardButton
		var row []tgbotapi.InlineKeyboardButton
		for _, node := range groupInfo.All {
			label := node
			if node == groupInfo.Now {
				label = "✿ " + node
			}
			btn := tgbotapi.NewInlineKeyboardButtonData(label, "choose_"+group+"_"+node)
			row = append(row, btn)
			if len(row) == 2 {
				rows = append(rows, row)
				row = []tgbotapi.InlineKeyboardButton{}
			}
		}
		if len(row) > 0 {
			rows = append(rows, row)
		}
		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Cek Delay", "check_delay_"+group),
				tgbotapi.NewInlineKeyboardButtonData("⬅️ Kembali", "back"),
			),
		)

		edit := tgbotapi.NewEditMessageTextAndMarkup(
			chatID, msgID,
			fmt.Sprintf("*%s* ➠ *%s*\nList node:", group, groupInfo.Now),
			tgbotapi.NewInlineKeyboardMarkup(rows...),
		)
		edit.ParseMode = "Markdown"
		bot.Send(edit)

	case strings.HasPrefix(data, "check_delay_"):
		group := strings.TrimPrefix(data, "check_delay_")
		var groupInfo ProxyItem
		if err := getJSON(mihomoAPI+"/proxies/"+group, &groupInfo); err != nil {
			bot.Send(tgbotapi.NewEditMessageText(chatID, msgID, "❌ Gagal ambil data group: "+err.Error()))
			return
		}

		text := fmt.Sprintf("*Delay untuk grup %s:*\n", group)
		for _, node := range groupInfo.All {
			var delayResp map[string]interface{}
			delayURL := fmt.Sprintf("%s/proxies/%s/delay?url=%s&timeout=5000&name=%s",
				mihomoAPI,
				url.PathEscape(group),
				url.QueryEscape("http://www.gstatic.com/generate_204"),
				url.QueryEscape(node),
			)
			_ = getJSON(delayURL, &delayResp)

			delay := "timeout"
			if d, ok := delayResp["delay"].(float64); ok {
				delay = strconv.Itoa(int(d))
			}
			text += fmt.Sprintf("❀ %s: %sms\n", node, delay)
		}

		edit := tgbotapi.NewEditMessageText(chatID, msgID, text)
		edit.ParseMode = "Markdown"
		bot.Send(edit)

	case strings.HasPrefix(data, "choose_"):
		parts := strings.SplitN(data, "_", 3)
		group, proxy := parts[1], parts[2]
		if _, err := sendRequest("PUT", mihomoAPI+"/proxies/"+group, []byte(fmt.Sprintf(`{"name":"%s"}`, proxy))); err != nil {
			bot.Send(tgbotapi.NewEditMessageText(chatID, msgID, "❌ Gagal set proxy: "+err.Error()))
			return
		}
		edit := tgbotapi.NewEditMessageText(chatID, msgID, fmt.Sprintf("Proxy *%s* ➠ *%s*", group, proxy))
		edit.ParseMode = "Markdown"
		bot.Send(edit)

	case data == "status":
		var proxiesResp ProxiesResponse
		if err := getJSON(mihomoAPI+"/proxies", &proxiesResp); err != nil {
			bot.Send(tgbotapi.NewEditMessageText(chatID, msgID, "❌ Gagal ambil status: "+err.Error()))
			return
		}

		text := "*Delay semua grup:*\n"
		for name, item := range proxiesResp.Proxies {
			if item.Type == "Selector" {
				var delayResp map[string]interface{}
				_ = getJSON(mihomoAPI+"/proxies/"+name+"/delay?url=http://www.gstatic.com/generate_204&timeout=3000", &delayResp)
				delay := "timeout"
				if d, ok := delayResp["delay"].(float64); ok {
					delay = strconv.Itoa(int(d))
				}
				text += fmt.Sprintf("❀ %s (%sms)\n", name, delay)
			}
		}
		edit := tgbotapi.NewEditMessageText(chatID, msgID, text)
		edit.ParseMode = "Markdown"
		bot.Send(edit)

	case data == "reload":
		sendRequest("PUT", mihomoAPI+"/configs?force=true", []byte(`{"path":"","payload":""}`))
		bot.Send(tgbotapi.NewEditMessageText(chatID, msgID, "Config berhasil di-reload."))

	case data == "restart":
		sendRequest("POST", mihomoAPI+"/restart", nil)
		bot.Send(tgbotapi.NewEditMessageText(chatID, msgID, "Restart selesai."))

	case data == "upgrade":
		resp, _ := sendRequest("POST", mihomoAPI+"/upgrade", nil)
		if resp != nil && resp.StatusCode == 200 {
			bot.Send(tgbotapi.NewEditMessageText(chatID, msgID, "Berhasil diperbarui."))
		} else {
			bot.Send(tgbotapi.NewEditMessageText(chatID, msgID, "Tidak ada pembaruan."))
		}

	case data == "version":
		var v map[string]string
		getJSON(mihomoAPI+"/version", &v)
		bot.Send(tgbotapi.NewEditMessageText(chatID, msgID, "`"+v["version"]+"`"))

	case data == "back":
		var proxiesResp ProxiesResponse
		getJSON(mihomoAPI+"/proxies", &proxiesResp)
		edit := tgbotapi.NewEditMessageTextAndMarkup(
			chatID, msgID,
			"Kembali ke menu utama:",
			yacdMenu(proxiesResp.Proxies),
		)
		bot.Send(edit)
	}
}