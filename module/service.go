package module

import (
	"fmt"
	"os/exec"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Jalankan command shell
func RunCommand(cmd string, args ...string) string {
	out, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error: %v\nOutput: %s", err, string(out))
	}
	return string(out)
}

// Menu utama
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
			tgbotapi.NewInlineKeyboardButtonData("📊 Status", "sys_status"),
			tgbotapi.NewInlineKeyboardButtonData("⬆️ Upgrade Core", "u"),
			tgbotapi.NewInlineKeyboardButtonData("⬆️ Upgrade UI", "x"),
		),
	)
}

// Submenu Service
func ServiceMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("▶️ Start", "s start"),
			tgbotapi.NewInlineKeyboardButtonData("⏹ Stop", "s stop"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔄 Restart", "s restart"),
			tgbotapi.NewInlineKeyboardButtonData("📊 Status", "s status"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⏰ Cron", "s cron"),
			tgbotapi.NewInlineKeyboardButtonData("🛑 Kill Cron", "s kcron"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Kembali", "mainmenu"),
		),
	)
}

// Submenu Iptables
func IptablesMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Enable", "i enable"),
			tgbotapi.NewInlineKeyboardButtonData("❌ Disable", "i disable"),
			tgbotapi.NewInlineKeyboardButtonData("🔄 Renew", "i renew"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Kembali", "mainmenu"),
		),
	)
}

// Submenu Tools
func ToolsMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔍 Check", "t check"),
			tgbotapi.NewInlineKeyboardButtonData("🧠 Memcg", "t memcg"),
			tgbotapi.NewInlineKeyboardButtonData("⚙️ Cpuset", "t cpuset"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💽 Blkio", "t blkio"),
			tgbotapi.NewInlineKeyboardButtonData("🔗 Bond0", "t bond0"),
			tgbotapi.NewInlineKeyboardButtonData("🔗 Bond1", "t bond1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🌍 GeoSub", "t geosub"),
			tgbotapi.NewInlineKeyboardButtonData("🌐 GeoX", "t geox"),
			tgbotapi.NewInlineKeyboardButtonData("📦 Subs", "t subs"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬆️ UpKernel", "t upkernel"),
			tgbotapi.NewInlineKeyboardButtonData("⬆️ UpXUI", "t upxui"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬆️ UpYQ", "t upyq"),
			tgbotapi.NewInlineKeyboardButtonData("⬆️ UpCurl", "t upcurl"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("♻️ Reload", "t reload"),
			tgbotapi.NewInlineKeyboardButtonData("🌐 Webroot", "t webroot"),
			tgbotapi.NewInlineKeyboardButtonData("🚀 All", "t all"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Kembali", "mainmenu"),
		),
	)
}

func CoreMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("clash", "core_clash"),
			tgbotapi.NewInlineKeyboardButtonData("sing-box", "core_sing-box"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("xray", "core_xray"),
			tgbotapi.NewInlineKeyboardButtonData("v2fly", "core_v2fly"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("hysteria", "core_hysteria"),
		),
	)
}

// NOTE: Status function is implemented in status.go, while this file keeps only shell/menu helpers.
