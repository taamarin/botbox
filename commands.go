package main

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandleHelp menampilkan daftar command
func HandleHelp(chatID int64, bot *tgbotapi.BotAPI) {
	helpText := `📖 *Daftar Command:*

/help       - Menampilkan bantuan
  └ Menampilkan daftar perintah lengkap

/menu       - Menampilkan menu utama
  └ Akses cepat ke menu sbfr

/import     - <path> (default: /data/adb/box/)
  └ Import file ke box (reply ke file)

/export     - <path/file> (default: /data/adb/box/)
  └ Export file dari box

/log        - Membaca isi runs.log
  └ Lihat log terakhir run sbfr

/sbfr       - Menu kontrountuk /system/bin/sbfr
  └ Jalankan, stop, restart, dan check status

/yacd       - Menu kontrodashboard YACD
  └ Pilih grup proxy, check delay, reload config

/core       - Pilih core untuk settings.ini
  └ Opsi: clash, sing-box, xray, v2fly, hysteria
`

	msg := tgbotapi.NewMessage(chatID, helpText)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

// HandleMenu menampilkan menu utama
func HandleMenu(chatID int64, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(chatID, "📌 *Menu Utama:*")
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = mainMenu()
	bot.Send(msg)
}

// Router untuk command text sederhana
func HandleBasicCommands(bot *tgbotapi.BotAPI, chatID int64, text string) {
	args := strings.Fields(text)
	if len(args) == 0 {
		return
	}

	switch args[0] {
	case "/help":
		HandleHelp(chatID, bot)
	case "/menu":
		HandleMenu(chatID, bot)
	}
}
