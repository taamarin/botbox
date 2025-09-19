package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"main/module"
	"gopkg.in/ini.v1"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	botToken	 string
	ownerID		int64
	mihomoAPI	string
	apiSecret	string
)

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

	// inisialisasi modul yacd
	module.Init(mihomoAPI, apiSecret)

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
			module.HandleBasicCommands(bot, update.Message.Chat.ID, update.Message.Text)

			switch args[0] {
			case "/myip":
				module.HandleMyIP(bot, update.Message.Chat.ID)
			case "/info":
				module.HandleInfo(bot, update.Message.Chat.ID)
			case "/ipinfo":
				module.HandleIPInfo(bot, update)
			case "/hostip":
				module.HandleHostIP(bot, update)
			case "/yacd":
				module.HandleYacd(bot, update.Message.Chat.ID)
			case "/speedtest":
				module.HandleSpeedTest(update.Message.Chat.ID, bot)
			case "/core":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Pilih core yang ingin digunakan:")
				msg.ReplyMarkup = module.CoreMenu()
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
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Pilih aksi untuk `/data/adb/modules/box_for_root/system/bin/sbfr`:")
				msg.ParseMode = "Markdown"
				msg.ReplyMarkup = module.MainMenu()
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
				output := module.RunCommand("sh", "-c", cmd)
			
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
				data == "traffic" ||
				data == "back" {
				module.HandleYacdCallback(bot, update.CallbackQuery)
				continue
			}

			// sbfr callback
			switch data {
			case "submenu_service":
				edit := tgbotapi.NewEditMessageTextAndMarkup(chatID, messageID, "⚡ *Service Commands*", module.ServiceMenu())
				edit.ParseMode = "Markdown"
				bot.Send(edit)

			case "submenu_iptables":
				edit := tgbotapi.NewEditMessageTextAndMarkup(chatID, messageID, "🛡 *Iptables Commands*", module.IptablesMenu())
				edit.ParseMode = "Markdown"
				bot.Send(edit)

			case "submenu_tools":
				edit := tgbotapi.NewEditMessageTextAndMarkup(chatID, messageID, "🛠 *Tools Commands*", module.ToolsMenu())
				edit.ParseMode = "Markdown"
				bot.Send(edit)

			case "mainmenu":
				edit := tgbotapi.NewEditMessageTextAndMarkup(chatID, messageID, "Pilih aksi untuk `/data/adb/modules/box_for_root/system/bin/sbfr`:", module.MainMenu())
				edit.ParseMode = "Markdown"
				bot.Send(edit)

			default:
				// tampilkan loading dulu
				loading := tgbotapi.NewEditMessageText(chatID, messageID, "⏳ Loading...")
				bot.Send(loading)

				var output string
				parts := strings.Fields(data)
				if len(parts) == 1 {
					output = module.RunCommand("/data/adb/modules/box_for_root/system/bin/sbfr", parts[0])
				} else if len(parts) == 2 {
					output = module.RunCommand("/data/adb/modules/box_for_root/system/bin/sbfr", parts[0], parts[1])
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

