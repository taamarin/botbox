package module

import (
	"fmt"
	"os/exec"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// RunInfoCommand util (sudah ada, tapi taruh sini kalau perlu stand-alone)
func RunInfoCommand(name string, args ...string) string {
	out, err := exec.Command(name, args...).CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error: %v\nOutput: %s", err, string(out))
	}
	return string(out)
}

// HandleInfo menampilkan info status ponsel Android
func HandleInfo(bot *tgbotapi.BotAPI, chatID int64) {
	// Ambil info dari dumpsys
	battery := RunInfoCommand("su", "-c", "dumpsys battery")
	uptime := RunInfoCommand("su", "-c", "uptime")

	// Format hasil
	info := fmt.Sprintf(
		"📱 *Android Device Info*\n\n"+
			"🔋 *Battery:*\n```\n%s```\n"+
			"⏱ *Uptime:*\n```\n%s```\n",
		strings.TrimSpace(battery),
		strings.TrimSpace(uptime),
	)

	// Telegram limit
	if len(info) > 4000 {
		info = info[:4000] + "\n..."
	}

	msg := tgbotapi.NewMessage(chatID, info)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}