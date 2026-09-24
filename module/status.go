package module

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func runShell(cmd string, args ...string) string {
	out, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		if len(out) > 0 {
			return strings.TrimSpace(string(out))
		}
		return fmt.Sprintf("Error: %v", err)
	}
	return strings.TrimSpace(string(out))
}

func isProcessRunning(name string) bool {
	out := runShell("sh", "-c", "ps -A 2>/dev/null | grep -E \""+name+"\" | grep -v grep | head -n 1")
	return strings.TrimSpace(out) != ""
}

func readSelectedCore() string {
	data, err := os.ReadFile("/data/adb/box/settings.ini")
	if err != nil {
		return "not configured"
	}
	text := string(data)
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "bin_name=") {
			return strings.TrimPrefix(line, "bin_name=")
		}
	}
	return "unknown"
}

func HandleHealth(bot *tgbotapi.BotAPI, chatID int64) {
	var out strings.Builder
	out.WriteString("🩺 *Bot Health Check*\n\n")

	checks := []struct {
		name  string
		value string
	}{
		{"mihomo", ""},
		{"clash", ""},
		{"sing-box", ""},
		{"xray", ""},
		{"hysteria", ""},
	}

	for i := range checks {
		checks[i].value = map[bool]string{true: "✅ Running", false: "⚠️ Not running"}[isProcessRunning(checks[i].name)]
		out.WriteString(fmt.Sprintf("• %s: %s\n", strings.Title(checks[i].name), checks[i].value))
	}

	var version map[string]string
	if err := getJSON(mihomoAPI+"/version", &version); err != nil {
		out.WriteString("\n🌐 Mihomo API: ❌ unreachable\n")
	} else {
		apiVersion := version["version"]
		if apiVersion == "" {
			apiVersion = "ok"
		}
		out.WriteString(fmt.Sprintf("\n🌐 Mihomo API: ✅ reachable (%s)\n", apiVersion))
	}

	out.WriteString(fmt.Sprintf("\n🧩 Selected core: %s\n", readSelectedCore()))

	msg := tgbotapi.NewMessage(chatID, out.String())
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔄 Check Again", "health_check"),
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Kembali", "mainmenu"),
		),
	)
	bot.Send(msg)
}

func HandleHealthCallback(bot *tgbotapi.BotAPI, cq *tgbotapi.CallbackQuery) {
	chatID := cq.Message.Chat.ID
	msgID := cq.Message.MessageID

	var sb strings.Builder
	sb.WriteString("🩺 *Bot Health Check*\n\n")
	for _, name := range []string{"mihomo", "clash", "sing-box", "xray", "hysteria"} {
		status := map[bool]string{true: "✅ Running", false: "⚠️ Not running"}[isProcessRunning(name)]
		sb.WriteString(fmt.Sprintf("• %s: %s\n", strings.Title(name), status))
	}

	var version map[string]string
	if err := getJSON(mihomoAPI+"/version", &version); err == nil {
		sb.WriteString(fmt.Sprintf("\n🌐 Mihomo API: ✅ reachable (%s)\n", version["version"]))
	} else {
		sb.WriteString("\n🌐 Mihomo API: ❌ unreachable\n")
	}
	sb.WriteString(fmt.Sprintf("\n🧩 Selected core: %s\n", readSelectedCore()))

	edit := tgbotapi.NewEditMessageTextAndMarkup(chatID, msgID, sb.String(), tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔄 Check Again", "health_check"),
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Kembali", "mainmenu"),
		),
	))
	edit.ParseMode = "Markdown"
	bot.Send(edit)
}

func formatTail(lines []string, limit int) string {
	if len(lines) <= limit {
		return strings.Join(lines, "\n")
	}
	return strings.Join(lines[len(lines)-limit:], "\n")
}

func TailLog(path string, limit int) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return "❌ Log tidak ditemukan: " + err.Error()
	}
	text := strings.TrimRight(string(data), "\n")
	if text == "" {
		return "Log kosong."
	}
	lines := strings.Split(text, "\n")
	if limit <= 0 || limit > len(lines) {
		limit = len(lines)
	}
	return formatTail(lines, limit)
}

