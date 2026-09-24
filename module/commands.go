package module

import (
    "strconv"
    "strings"

    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleLog(bot *tgbotapi.BotAPI, chatID int64, maxLines int) {
    path := "/data/adb/box/run/runs.log"
    text := TailLog(path, maxLines)
    if len(text) > 4000 {
        msg := tgbotapi.NewMessage(chatID, "📄 Output log terlalu panjang, kirim ringkasan terakhir:\n\n```\n"+text[:4000]+"\n```")
        msg.ParseMode = "Markdown"
        bot.Send(msg)
        return
    }

    msg := tgbotapi.NewMessage(chatID, "```\n"+text+"\n```")
    msg.ParseMode = "Markdown"
    bot.Send(msg)
}

func HandleLogCommand(bot *tgbotapi.BotAPI, chatID int64, rawArgs []string) {
    limit := 50
    if len(rawArgs) > 0 {
        if n, err := strconv.Atoi(rawArgs[0]); err == nil && n > 0 {
            limit = n
        }
    }
    HandleLog(bot, chatID, limit)
}

func HandleHelp(chatID int64, bot *tgbotapi.BotAPI) {
    helpText := `📖 *Daftar Command:*

/help       - Menampilkan bantuan
  └ Menampilkan daftar perintah lengkap
/menu       - Menampilkan menu utama
  └ Akses cepat ke menu sbfr
/status     - Menampilkan ringkasan status sistem Android
  └ Battery, uptime, dan memori
/health     - Cek status proses & koneksi Mihomo
  └ Mendeteksi apakah core dan API hidup
/log        - Menampilkan 50 baris log terakhir
  └ Gunakan /log 100 untuk ambil baris lebih banyak
/import     - <path> (default: /data/adb/box/)
  └ Import file ke box (reply ke file)
/export     - <path/file> (default: /data/adb/box/)
  └ Export file dari box
/sbfr       - Menu kontrountuk /system/bin/sbfr
  └ Jalankan, stop, restart, dan check status
/yacd       - Menu kontrodashboard YACD
  └ Pilih grup proxy, check delay, reload config
/core       - Pilih core untuk settings.ini
  └ Opsi: clash, sing-box, xray, v2fly, hysteria
/speedtest  - Pilih aksi SpeedTest
  └ Run SpeedTest
`

    msg := tgbotapi.NewMessage(chatID, helpText)
    msg.ParseMode = "Markdown"
    bot.Send(msg)
}

func HandleMenu(chatID int64, bot *tgbotapi.BotAPI) {
    msg := tgbotapi.NewMessage(chatID, "📌 *Menu Utama:*")
    msg.ParseMode = "Markdown"
    msg.ReplyMarkup = MainMenu()
    bot.Send(msg)
}

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
    case "/status":
        HandleStatus(bot, chatID)
    case "/health":
        HandleHealth(bot, chatID)
    }
}
