package module

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	mihomoAPI string
	apiSecret string
)

var mihomoClient = &http.Client{Timeout: 10 * time.Second}

func Init(api, secret string) {
	mihomoAPI = strings.TrimRight(api, "/")
	apiSecret = secret
}

type ProxyItem struct {
	Type string   `json:"type"`
	Now  string   `json:"now"`
	All  []string `json:"all"`
}

type ProxiesResponse struct {
	Proxies map[string]ProxyItem `json:"proxies"`
}

func newAPIRequest(method, requestURL string, payload []byte) (*http.Request, error) {
	req, err := http.NewRequest(method, requestURL, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	if apiSecret != "" {
		req.Header.Set("Authorization", "Bearer "+apiSecret)
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}

// helper GET
func getJSON(requestURL string, target interface{}) error {
	req, err := newAPIRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return err
	}
	resp, err := mihomoClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("Mihomo API returned %s", resp.Status)
	}
	return json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(target)
}

// helper PUT/POST
func sendRequest(method, requestURL string, payload []byte) (*http.Response, error) {
	req, err := newAPIRequest(method, requestURL, payload)
	if err != nil {
		return nil, err
	}
	return mihomoClient.Do(req)
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
			tgbotapi.NewInlineKeyboardButtonData("Traffic", "traffic"), // 👈 tambahan
		),
	)

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// handler command /yacd
func HandleYacd(bot *tgbotapi.BotAPI, chatID int64) {
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
func HandleYacdCallback(bot *tgbotapi.BotAPI, cq *tgbotapi.CallbackQuery) {
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
		if err := getJSON(mihomoAPI+"/proxies/"+url.PathEscape(group), &groupInfo); err != nil {
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
		if err := getJSON(mihomoAPI+"/proxies/"+url.PathEscape(group), &groupInfo); err != nil {
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
		if len(parts) != 3 {
			bot.Send(tgbotapi.NewEditMessageText(chatID, msgID, "❌ Callback tidak valid."))
			return
		}
		group, proxy := parts[1], parts[2]
		payload, err := json.Marshal(map[string]string{"name": proxy})
		if err != nil {
			bot.Send(tgbotapi.NewEditMessageText(chatID, msgID, "❌ Gagal set proxy."))
			return
		}
		resp, err := sendRequest(http.MethodPut, mihomoAPI+"/proxies/"+url.PathEscape(group), payload)
		if err == nil {
			defer resp.Body.Close()
		}
		if err != nil || resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
			if err == nil {
				err = fmt.Errorf("Mihomo API returned %s", resp.Status)
			}
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
				_ = getJSON(mihomoAPI+"/proxies/"+url.PathEscape(name)+"/delay?url=http://www.gstatic.com/generate_204&timeout=3000", &delayResp)
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
		resp, err := sendRequest(http.MethodPut, mihomoAPI+"/configs?force=true", []byte(`{"path":"","payload":""}`))
		if err == nil {
			defer resp.Body.Close()
		}
		if err != nil || resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
			bot.Send(tgbotapi.NewEditMessageText(chatID, msgID, "❌ Gagal reload config."))
		} else {
			bot.Send(tgbotapi.NewEditMessageText(chatID, msgID, "Config berhasil di-reload."))
		}

	case data == "restart":
		resp, err := sendRequest(http.MethodPost, mihomoAPI+"/restart", nil)
		if err == nil {
			defer resp.Body.Close()
		}
		if err != nil || resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
			bot.Send(tgbotapi.NewEditMessageText(chatID, msgID, "❌ Gagal restart."))
		} else {
			bot.Send(tgbotapi.NewEditMessageText(chatID, msgID, "Restart selesai."))
		}

	case data == "upgrade":
		resp, err := sendRequest(http.MethodPost, mihomoAPI+"/upgrade", nil)
		if resp != nil {
			defer resp.Body.Close()
		}
		if err == nil && resp.StatusCode == http.StatusOK {
			bot.Send(tgbotapi.NewEditMessageText(chatID, msgID, "Berhasil diperbarui."))
		} else {
			bot.Send(tgbotapi.NewEditMessageText(chatID, msgID, "Tidak ada pembaruan."))
		}

	case data == "version":
		var v map[string]string
		getJSON(mihomoAPI+"/version", &v)
		bot.Send(tgbotapi.NewEditMessageText(chatID, msgID, "`"+v["version"]+"`"))

	case data == "traffic":
		var conns map[string]interface{}
		if err := getJSON(mihomoAPI+"/connections", &conns); err != nil {
			bot.Send(tgbotapi.NewEditMessageText(chatID, msgID, "❌ Gagal ambil status: "+err.Error()))
			return
		}
		
		total := 0
		if connections, ok := conns["connections"].([]interface{}); ok {
			total = len(connections)
		}
		
		// Get download/upload stats
		download := 0.0
		upload := 0.0
		if traffic, ok := conns["downloadTotal"].(float64); ok {
			download = traffic / 1024 / 1024 // Convert to MB
		}
		if traffic, ok := conns["uploadTotal"].(float64); ok {
			upload = traffic / 1024 / 1024 // Convert to MB
		}

		statusText := fmt.Sprintf("📊 *Status Koneksi:*\n"+
			"• Koneksi aktif: %d\n"+
			"• Download: %.2f MB\n"+
			"• Upload: %.2f MB", total, download, upload)
		
		edit := tgbotapi.NewEditMessageText(chatID, msgID, statusText)
		edit.ParseMode = "Markdown"
		bot.Send(edit)

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