package main

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandleHelp menampilkan daftar command
func HandleHelp(chatID int64, bot *tgbotapi.BotAPI) {
	helpText := `📖 *Daftar Command:*

/help - Menampilkan bantuan
/menu - Menampilkan menu utama
/import - Import file (reply ke file)
/export <namafile> - Export file dari box
/log - Baca isi runs.log
/sbfr - Menu kontrol untuk /system/bin/sbfr
/yacd - Menu kontrol dashboard
/core - Pilih core untuk settings.ini
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
