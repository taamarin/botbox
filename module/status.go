package module

import (
	"fmt"
	"os/exec"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func shellOutput(cmd string, args ...string) string {
	out, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		if len(out) > 0 {
			return strings.TrimSpace(string(out))
		}
		return fmt.Sprintf("Error: %v", err)
	}
	return strings.TrimSpace(string(out))
}

func systemStatusText() string {
	battery := shellOutput("su", "-c", "dumpsys battery 2>/dev/null | sed -n '1,40p'")
	if battery == "" {
		battery = "N/A"
	}

	uptime := shellOutput("su", "-c", "uptime 2>/dev/null")
	if uptime == "" {
		uptime = "N/A"
	}

	memInfo := shellOutput("cat", "/proc/meminfo")
	memSummary := memInfo
	if memSummary == "" {
		memSummary = "N/A"
	}

	text := fmt.Sprintf(
		"📊 *System Status*\n\n"+
			"🔋 *Battery*\n```\n%s\n```\n\n"+
			"⏱ *Uptime*\n```\n%s\n```\n\n"+
			"💾 *Memory*\n```\n%s\n```",
			battery,
			uptime,
			memSummary,
	)

	if len(text) > 4000 {
		return text[:4000] + "\n..."
	}
	return text
}

func HandleStatus(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, systemStatusText())
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔄 Refresh", "sys_status"),
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Kembali", "mainmenu"),
		),
	)
	bot.Send(msg)
}

func HandleStatusCallback(bot *tgbotapi.BotAPI, cq *tgbotapi.CallbackQuery) {
	chatID := cq.Message.Chat.ID
	msgID := cq.Message.MessageID

	edit := tgbotapi.NewEditMessageTextAndMarkup(
		chatID,
		msgID,
		systemStatusText(),
		tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🔄 Refresh", "sys_status"),
				tgbotapi.NewInlineKeyboardButtonData("⬅️ Kembali", "mainmenu"),
			),
		),
	)
	edit.ParseMode = "Markdown"
	bot.Send(edit)
}
