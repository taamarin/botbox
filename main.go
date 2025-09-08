package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/ini.v1"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	botToken	 string
	ownerID		int64
	mihomoAPI	string
	apiSecret	string
)

// Jalankan command shell
func runCommand(cmd string, args ...string) string {
	out, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error: %v\nOutput: %s", err, string(out))
	}
	return string(out)
}

// Menu utama
func mainMenu() tgbotapi.InlineKeyboardMarkup {
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

// Submenu Service
func serviceMenu() tgbotapi.InlineKeyboardMarkup {
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
func iptablesMenu() tgbotapi.InlineKeyboardMarkup {
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
func toolsMenu() tgbotapi.InlineKeyboardMarkup {
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

func coreMenu() tgbotapi.InlineKeyboardMarkup {
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

var configPath string

func main() {
	// baca path dari argumen -c
	flag.StringVar(&configPath, "c", "", "Path ke bot.ini")
	flag.Parse()

	if configPath == "" {
		// default: folder binary
		exPath, err := os.Executable()
		if err != nil {
			log.Fatalf("Gagal dapat path executable: %v", err)
		}
		exDir := filepath.Dir(exPath)
		configPath = filepath.Join(exDir, "bot.ini")
	}

	cfg, err := ini.Load(configPath)
	if err != nil {
		log.Fatalf("Gagal baca bot.ini di %s: %v", configPath, err)
	}

	botToken := cfg.Section("bot").Key("token").String()
	ownerStr := cfg.Section("bot").Key("owner").String()
	ownerID, _ = strconv.ParseInt(ownerStr, 10, 64)

	mihomoAPI = cfg.Section("mihomo").Key("api").String()
	apiSecret = cfg.Section("mihomo").Key("secret").String()

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Bot jalan sebagai %s", bot.Self.UserName)


	// Kirim notifikasi ke owner
	startupMsg := tgbotapi.NewMessage(ownerID, fmt.Sprintf("✅ Bot *%s* berhasil dijalankan! /help", bot.Self.UserName))
	startupMsg.ParseMode = "Markdown"
	bot.Send(startupMsg)
	
	// Ambil update terakhir dulu supaya mengabaikan pesan lama
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	
	updates, _ := bot.GetUpdates(u)
	if len(updates) > 0 {
			lastUpdate := updates[len(updates)-1]
			u.Offset = lastUpdate.UpdateID + 1
	}
	
	// Mulai channel update
	updatesChan := bot.GetUpdatesChan(u)

	for update := range updatesChan {
		// --- Handle pesan ---
		if update.Message != nil && update.Message.Text != "" {
			// cek akses
			if update.Message.From.ID != ownerID {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "❌ Kamu tidak punya akses."))
				continue
			}

			args := strings.Fields(update.Message.Text)
			HandleBasicCommands(bot, update.Message.Chat.ID, update.Message.Text)

			switch args[0] {
			case "/yacd":
				handleYacd(bot, update.Message.Chat.ID)

			case "/core":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Pilih core yang ingin digunakan:")
				msg.ReplyMarkup = coreMenu()
				bot.Send(msg)

			case "/import":
				// default folder
				defaultPath := "/data/adb/box"
				
				// target path: argumen pertama kalau ada
				targetPath := defaultPath
				if len(args) > 1 {
					targetPath = args[1]
				}
		
				if update.Message.ReplyToMessage != nil && update.Message.ReplyToMessage.Document != nil {
					doc := update.Message.ReplyToMessage.Document
					file, err := bot.GetFile(tgbotapi.FileConfig{FileID: doc.FileID})
					if err != nil {
						bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "❌ Failed to get file: "+err.Error()))
						break
					}
	
					downloadURL := file.Link(bot.Token)
					resp, err := http.Get(downloadURL)
					if err != nil {
						bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "❌ Failed to download file: "+err.Error()))
						break
					}
					defer resp.Body.Close()
	
					savePath := filepath.Join(targetPath, doc.FileName)
					out, err := os.Create(savePath)
					if err != nil {
						bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "❌ Failed to create file: "+err.Error()))
						break
					}
					defer out.Close()
	
					_, err = io.Copy(out, resp.Body)
					if err != nil {
						bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "❌ Failed to save file: "+err.Error()))
						break
					}
	
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "✅ File saved to "+savePath))
			} else {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Use `/import <path>` with reply to a file."))
				}

			case "/export":
				if len(args) < 2 {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Usage: /export <namafile>"))
					break
				}
				filePath := filepath.Join("/data/adb/box", args[1])
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "❌ File tidak ditemukan."))
					break
				}
				doc := tgbotapi.NewDocument(update.Message.Chat.ID, tgbotapi.FilePath(filePath))
				bot.Send(doc)

			case "/log":
				filePath := "/data/adb/box/run/runs.log"
				data, err := os.ReadFile(filePath)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "❌ Gagal baca log: "+err.Error()))
					break
				}
				if len(data) > 4000 {
					doc := tgbotapi.NewDocument(update.Message.Chat.ID, tgbotapi.FilePath(filePath))
					bot.Send(doc)
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, string(data))
					bot.Send(msg)
				}

			case "/sbfr":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Pilih aksi untuk `/system/bin/sbfr`:")
				msg.ParseMode = "Markdown"
				msg.ReplyMarkup = mainMenu()
				bot.Send(msg)
			}
		}

		// --- Handle callback query ---
		if update.CallbackQuery != nil {
			if update.CallbackQuery.From.ID != ownerID {
				bot.Request(tgbotapi.NewCallback(update.CallbackQuery.ID, "❌ Tidak ada akses."))
				continue
			}

			data := update.CallbackQuery.Data
			chatID := update.CallbackQuery.Message.Chat.ID
			messageID := update.CallbackQuery.Message.MessageID

			if strings.HasPrefix(data, "core_") {
				selecCore := strings.TrimPrefix(data, "core_")
				cmd := fmt.Sprintf("sed -i 's/bin_name=.*/bin_name=%s/g' /data/adb/box/settings.ini", selecCore)
				output := runCommand("sh", "-c", cmd)
			
				resultText := strings.TrimSpace(output)
				if resultText == "" {
					resultText = fmt.Sprintf("✅ Core berhasil diubah menjadi: %s", selecCore)
				}
			
				// tombol kembali
				backMenu := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("⬅️ Kembali", "mainmenu"),
					),
				)
			
				edit := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID,
					update.CallbackQuery.Message.MessageID,
					resultText,
					backMenu,
				)
				edit.ParseMode = "Markdown"
				bot.Send(edit)
				continue
			}

			// Yacd callback
			if strings.HasPrefix(data, "select_") ||
				strings.HasPrefix(data, "choose_") ||
				strings.HasPrefix(data, "check_delay_") ||
				data == "status" ||
				data == "reload" ||
				data == "restart" ||
				data == "upgrade" ||
				data == "version" ||
				data == "back" {
				handleYacdCallback(bot, update.CallbackQuery)
				continue
			}

			// sbfr callback
			switch data {
			case "submenu_service":
				edit := tgbotapi.NewEditMessageTextAndMarkup(chatID, messageID, "⚡ *Service Commands*", serviceMenu())
				edit.ParseMode = "Markdown"
				bot.Send(edit)

			case "submenu_iptables":
				edit := tgbotapi.NewEditMessageTextAndMarkup(chatID, messageID, "🛡 *Iptables Commands*", iptablesMenu())
				edit.ParseMode = "Markdown"
				bot.Send(edit)

			case "submenu_tools":
				edit := tgbotapi.NewEditMessageTextAndMarkup(chatID, messageID, "🛠 *Tools Commands*", toolsMenu())
				edit.ParseMode = "Markdown"
				bot.Send(edit)

			case "mainmenu":
				edit := tgbotapi.NewEditMessageTextAndMarkup(chatID, messageID, "Pilih aksi untuk `/system/bin/sbfr`:", mainMenu())
				edit.ParseMode = "Markdown"
				bot.Send(edit)

			default:
				// tampilkan loading dulu
				loading := tgbotapi.NewEditMessageText(chatID, messageID, "⏳ Loading...")
				bot.Send(loading)

				var output string
				parts := strings.Fields(data)
				if len(parts) == 1 {
					output = runCommand("/system/bin/sbfr", parts[0])
				} else if len(parts) == 2 {
					output = runCommand("/system/bin/sbfr", parts[0], parts[1])
				}

				// tombol kembali
				backMenu := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("⬅️ Kembali ke Menu", "mainmenu"),
					),
				)

				// hasil eksekusi
				resultText := strings.TrimSpace(output)
				if len(resultText) > 4000 {
					resultText = "❗ Output terlalu panjang untuk ditampilkan."
				} else {
					resultText = "```\n" + resultText + "\n```"
				}

				edit := tgbotapi.NewEditMessageTextAndMarkup(chatID, messageID, resultText, backMenu)
				edit.ParseMode = "MarkdownV2"
				bot.Send(edit)
			}
		}
	}
}