package module

import (
	"fmt"
	"os/exec"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func RunCommand(cmd string, args ...string) string {
	out, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error: %v\nOutput: %s", err, string(out))
	}
	return string(out)
}

func shellOutput(cmd string, args ...string) string {
	return strings.TrimSpace(RunCommand(cmd, args...))
}

func systemStatusText() string {
	battery := shellOutput("su", "-c", "dumpsys battery 2>/dev/null | sed -n '1,40p'")
	if battery == "" { battery = "N/A" }
	uptime := shellOutput("su", "-c", "uptime 2>/dev/null")
	if uptime == "" { uptime = "N/A" }
	memInfo := shellOutput("cat", "/proc/meminfo")
	if memInfo == "" { memInfo = "N/A" }
	text := fmt.Sprintf("📊 *System Status*\n\n🔋 *Battery*\n```\n%s\n```\n\n⏱ *Uptime*\n```\n%s\n```\n\n💾 *Memory*\n```\n%s\n```", battery, uptime, memInfo)
	if len(text) > 4000 { return text[:4000] + "\n..." }
	return text
}

func HandleStatus(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, systemStatusText())
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🔄 Refresh", "sys_status"),
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Kembali", "mainmenu"),
	))
	bot.Send(msg)
}

func HandleStatusCallback(bot *tgbotapi.BotAPI, cq *tgbotapi.CallbackQuery) {
	if cq.Message == nil { return }
	edit := tgbotapi.NewEditMessageTextAndMarkup(cq.Message.Chat.ID, cq.Message.MessageID, systemStatusText(), tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🔄 Refresh", "sys_status"),
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Kembali", "mainmenu"),
	)))
	edit.ParseMode = "Markdown"
	bot.Send(edit)
}

func MainMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("▶️ Start", "start"),
			tgbotapi.NewInlineKeyboardButtonData("⏹ Stop", "stop"),
			tgbotapi.NewInlineKeyboardButtonData("🔄 Restart", "r"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⚡ Service", "submenu_service"),
			tgbotapi.NewInlineKeyboardButtonData("🛡 Iptables", "submenu_iptables"),
			tgbotapi.NewInlineKeyboardButtonData("🛠 Tools", "submenu_tools"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬆️ Upgrade Core", "u"),
			tgbotapi.NewInlineKeyboardButtonData("⬆️ Upgrade UI", "x"),
		),
	)
}

func ServiceMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("▶️ Start", "s start"), tgbotapi.NewInlineKeyboardButtonData("⏹ Stop", "s stop")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔄 Restart", "s restart"), tgbotapi.NewInlineKeyboardButtonData("📊 Status", "s status")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("⏰ Cron", "s cron"), tgbotapi.NewInlineKeyboardButtonData("🛑 Kill Cron", "s kcron")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("⬅️ Kembali", "mainmenu")),
	)
}

func IptablesMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("✅ Enable", "i enable"), tgbotapi.NewInlineKeyboardButtonData("❌ Disable", "i disable"), tgbotapi.NewInlineKeyboardButtonData("🔄 Renew", "i renew")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("⬅️ Kembali", "mainmenu")),
	)
}

func ToolsMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🔍 Check", "t check"), tgbotapi.NewInlineKeyboardButtonData("🧠 Memcg", "t memcg"), tgbotapi.NewInlineKeyboardButtonData("⚙️ Cpuset", "t cpuset")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("💽 Blkio", "t blkio"), tgbotapi.NewInlineKeyboardButtonData("🔗 Bond0", "t bond0"), tgbotapi.NewInlineKeyboardButtonData("🔗 Bond1", "t bond1")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🌍 GeoSub", "t geosub"), tgbotapi.NewInlineKeyboardButtonData("🌐 GeoX", "t geox"), tgbotapi.NewInlineKeyboardButtonData("📦 Subs", "t subs")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("♻️ Reload", "t reload"), tgbotapi.NewInlineKeyboardButtonData("🌐 Webroot", "t webroot"), tgbotapi.NewInlineKeyboardButtonData("🚀 All", "t all")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("⬅️ Kembali", "mainmenu")),
	)
}

func CoreMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("clash", "core_clash"), tgbotapi.NewInlineKeyboardButtonData("sing-box", "core_sing-box")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("xray", "core_xray"), tgbotapi.NewInlineKeyboardButtonData("v2fly", "core_v2fly")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("hysteria", "core_hysteria")),
	)
}
