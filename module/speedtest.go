package module

import (
	"fmt"

	"github.com/showwin/speedtest-go/speedtest"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleSpeedTest(chatID int64, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(chatID, "⏳ Running speed test, please wait...")
	sentMsg, _ := bot.Send(msg)

	// Fetch user info untuk ISP
	user, err := speedtest.FetchUserInfo()
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "❌ Failed to fetch user info: "+err.Error()))
		return
	}

	// Fetch servers
	serverList, err := speedtest.FetchServers()
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "❌ Failed to fetch servers: "+err.Error()))
		return
	}

	targets, err := serverList.FindServer([]int{})
	if err != nil || len(targets) == 0 {
		bot.Send(tgbotapi.NewMessage(chatID, "❌ No server found"))
		return
	}

	// Ambil server pertama
	srv := targets[0]

	// Run speed test
	if err := srv.PingTest(nil); err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "❌ PingTest failed: "+err.Error()))
		return
	}
	if err := srv.DownloadTest(); err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "❌ DownloadTest failed: "+err.Error()))
		return
	}
	if err := srv.UploadTest(); err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "❌ UploadTest failed: "+err.Error()))
		return
	}

	// Format hasil
	result := fmt.Sprintf(
		"📊 Speed Test Result:\nISP: %s\nServer: %s (%s)\nPing: %.2f ms\nDownload: %.2f Mbps\nUpload: %.2f Mbps",
		user.Isp,
		srv.Name,
		srv.Host,
		srv.Latency.Seconds()*1000,
		(srv.DLSpeed/1024/1024)*8,
		(srv.ULSpeed/1024/1024)*8,
	)

	// Reset context
	srv.Context.Reset()

	// Update pesan dengan hasil
	edit := tgbotapi.NewEditMessageText(chatID, sentMsg.MessageID, result)
	bot.Send(edit)
}