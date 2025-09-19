package module

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Handler untuk /hostip
func HandleHostIP(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	var domainText string

	// Ambil dari argumen /hostip <domain>
	if update.Message.CommandArguments() != "" {
		domainText = update.Message.CommandArguments()
	} else if update.Message.ReplyToMessage != nil && update.Message.ReplyToMessage.Text != "" {
		// Atau ambil dari reply message
		domainText = strings.TrimSpace(update.Message.ReplyToMessage.Text)
	} else {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "/hostip example.com\nAtau reply pesan yang berisi domain/URL")
		bot.Send(msg)
		return
	}

	// Regex untuk cari domain atau URL
	re := regexp.MustCompile(`[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}(?:/[^\s]*)?`)
	candidates := re.FindAllString(domainText, -1)
	domains := make(map[string]struct{})

	for _, item := range candidates {
		if strings.HasPrefix(item, "http://") || strings.HasPrefix(item, "https://") {
			if parsed, err := url.Parse(item); err == nil && parsed.Hostname() != "" {
				domains[parsed.Hostname()] = struct{}{}
			}
		} else {
			if parsed, err := url.Parse("http://" + item); err == nil && parsed.Hostname() != "" {
				domains[parsed.Hostname()] = struct{}{}
			}
		}
	}

	if len(domains) == 0 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Tidak ada domain valid ditemukan.")
		bot.Send(msg)
		return
	}

	var results []string
	for domain := range domains {
		ips, err := net.LookupIP(domain)
		if err != nil {
			results = append(results, fmt.Sprintf("%s\nTidak valid atau tidak dapat diakses.", domain))
			continue
		}

		var ipv4s []string
		var ipv6s []string
		for _, ip := range ips {
			if ip.To4() != nil {
				ipv4s = append(ipv4s, ip.String())
			} else {
				ipv6s = append(ipv6s, ip.String())
			}
		}

		res := fmt.Sprintf("%s", domain)
		if len(ipv4s) > 0 {
			res += "\nIPv4: " + strings.Join(ipv4s, ", ")
		}
		if len(ipv6s) > 0 {
			res += "\nIPv6: " + strings.Join(ipv6s, ", ")
		}
		if len(ipv4s) == 0 && len(ipv6s) == 0 {
			res += "\nTidak ditemukan alamat IP."
		}
		results = append(results, res)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, strings.Join(results, "\n\n"))
	bot.Send(msg)
}