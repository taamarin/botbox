package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"main/module"
	"gopkg.in/ini.v1"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	botToken string
	ownerID int64
	mihomoAPI string
	apiSecret string
)

func loadConfig(path string) error {
	cfg, err := ini.Load(path)
	if err != nil { return err }
	botToken = strings.TrimSpace(cfg.Section("bot").Key("token").String())
	ownerText := strings.TrimSpace(cfg.Section("bot").Key("owner").String())
	ownerID, err = strconv.ParseInt(ownerText, 10, 64)
	if err != nil || ownerID == 0 { return fmt.Errorf("owner harus berupa Telegram ID yang valid") }
	mihomoAPI = strings.TrimRight(strings.TrimSpace(cfg.Section("mihomo").Key("api").String()), "/")
	if mihomoAPI == "" { mihomoAPI = strings.TrimRight(strings.TrimSpace(cfg.Section("bot").Key("mihomo_api").String()), "/") }
	apiSecret = cfg.Section("mihomo").Key("secret").String()
	if apiSecret == "" { apiSecret = cfg.Section("bot").Key("api_secret").String() }
	if botToken == "" || mihomoAPI == "" { return fmt.Errorf("config tidak lengkap: token, owner, dan mihomo API wajib diisi") }
	return nil
}

func main() {
	configPath := flag.String("c", "", "Path ke bot.ini")
	flag.Parse()
	if *configPath == "" {
		exe, err := os.Executable()
		if err != nil { log.Fatal(err) }
		*configPath = filepath.Join(filepath.Dir(exe), "bot.ini")
	}
	if err := loadConfig(*configPath); err != nil { log.Fatalf("Gagal membaca bot.ini: %v", err) }
	module.Init(mihomoAPI, apiSecret)
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil { log.Fatal(err) }
	log.Printf("Bot berjalan sebagai %s", bot.Self.UserName)
	module.StartHealthWatcher(bot, ownerID, time.Minute)
	startup := tgbotapi.NewMessage(ownerID, fmt.Sprintf("✅ Bot *%s* berhasil dijalankan! /help", bot.Self.UserName))
	startup.ParseMode = "Markdown"
	bot.Send(startup)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	updates, err := bot.GetUpdates(u)
	if err == nil && len(updates) > 0 { u.Offset = updates[len(updates)-1].UpdateID + 1 }
	for update := range bot.GetUpdatesChan(u) {
		if update.Message != nil && update.Message.Text != "" {
			handleMessage(bot, update)
		}
		if update.CallbackQuery != nil {
			handleCallback(bot, update.CallbackQuery)
		}
	}
}

func authorized(bot *tgbotapi.BotAPI, id int64, chatID int64) bool {
	if id == ownerID { return true }
	bot.Send(tgbotapi.NewMessage(chatID, "❌ Kamu tidak punya akses."))
	return false
}

func handleMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := update.Message
	if !authorized(bot, msg.From.ID, msg.Chat.ID) { return }
	args := strings.Fields(msg.Text)
	if len(args) == 0 { return }
	module.HandleBasicCommands(bot, msg.Chat.ID, msg.Text)
	switch strings.Split(args[0], "@")[0] {
	case "/myip": module.HandleMyIP(bot, msg.Chat.ID)
	case "/info": module.HandleInfo(bot, msg.Chat.ID)
	case "/status": module.HandleStatus(bot, msg.Chat.ID)
	case "/health": module.HandleHealth(bot, msg.Chat.ID)
	case "/ipinfo": module.HandleIPInfo(bot, update)
	case "/hostip": module.HandleHostIP(bot, update)
	case "/yacd": module.HandleYacd(bot, msg.Chat.ID)
	case "/speedtest": module.HandleSpeedTest(msg.Chat.ID, bot)
	case "/core":
		m := tgbotapi.NewMessage(msg.Chat.ID, "Pilih core yang ingin digunakan:")
		m.ReplyMarkup = module.CoreMenu(); bot.Send(m)
	case "/sbfr":
		m := tgbotapi.NewMessage(msg.Chat.ID, "Pilih aksi untuk sbfr:"); m.ReplyMarkup = module.MainMenu(); bot.Send(m)
	case "/log":
		lines := 50
		if len(args) > 1 { if n, e := strconv.Atoi(args[1]); e == nil && n > 0 && n <= 500 { lines = n } }
		text := module.TailLog("/data/adb/box/run/runs.log", lines)
		if len(text) > 4000 { text = text[:4000] }
		m := tgbotapi.NewMessage(msg.Chat.ID, "```\n"+text+"\n```"); m.ParseMode = "Markdown"; bot.Send(m)
	}
}

func handleCallback(bot *tgbotapi.BotAPI, cq *tgbotapi.CallbackQuery) {
	if cq.Message == nil { return }
	if !authorized(bot, cq.From.ID, cq.Message.Chat.ID) { bot.Request(tgbotapi.NewCallback(cq.ID, "Tidak ada akses")); return }
	data := cq.Data
	if data == "health_check" { module.HandleHealthCallback(bot, cq); return }
	if data == "sys_status" { module.HandleStatusCallback(bot, cq); return }
	if strings.HasPrefix(data, "select_") || strings.HasPrefix(data, "choose_") || strings.HasPrefix(data, "check_delay_") || data == "status" || data == "reload" || data == "restart" || data == "upgrade" || data == "version" || data == "traffic" || data == "back" { module.HandleYacdCallback(bot, cq); return }
	chatID, messageID := cq.Message.Chat.ID, cq.Message.MessageID
	switch data {
	case "submenu_service": edit(bot, chatID, messageID, "⚡ *Service Commands*", module.ServiceMenu())
	case "submenu_iptables": edit(bot, chatID, messageID, "🛡 *Iptables Commands*", module.IptablesMenu())
	case "submenu_tools": edit(bot, chatID, messageID, "🛠 *Tools Commands*", module.ToolsMenu())
	case "mainmenu": edit(bot, chatID, messageID, "Pilih aksi untuk sbfr:", module.MainMenu())
	case "core_clash", "core_sing-box", "core_xray", "core_v2fly", "core_hysteria":
		core := strings.TrimPrefix(data, "core_")
		edit(bot, chatID, messageID, fmt.Sprintf("✅ Core dipilih: `%s`", core), tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("⬅️ Kembali", "mainmenu"))))
	}
}

func edit(bot *tgbotapi.BotAPI, chatID int64, messageID int, text string, markup interface{}) {
	m := tgbotapi.NewEditMessageText(chatID, messageID, text)
	if k, ok := markup.(tgbotapi.InlineKeyboardMarkup); ok { m.ReplyMarkup = &k }
	m.ParseMode = "Markdown"
	bot.Send(m)
}
