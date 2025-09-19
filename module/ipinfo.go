package module

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const IPINFO_TOKEN = ""
const IPINFO_URL = "https://ipinfo.io/%s?token=" + IPINFO_TOKEN

// Extract domain/URL dari teks
func extractDomains(text string) []string {
	re := regexp.MustCompile(`(https?://[^\s]+|[a-zA-Z0-9.-]+\.[a-zA-Z]{2,})`)
	matches := re.FindAllString(text, -1)
	domains := make(map[string]struct{})

	for _, match := range matches {
		if strings.HasPrefix(match, "http://") || strings.HasPrefix(match, "https://") {
			if parsed, err := url.Parse(match); err == nil && parsed.Hostname() != "" {
				domains[strings.ToLower(parsed.Hostname())] = struct{}{}
			}
		} else {
			if parsed, err := url.Parse("http://" + match); err == nil && parsed.Hostname() != "" {
				domains[strings.ToLower(parsed.Hostname())] = struct{}{}
			}
		}
	}

	var result []string
	for d := range domains {
		result = append(result, d)
	}
	return result
}

// Resolve domain jadi IP
func resolveDomain(domain string) (string, error) {
	ips, err := net.LookupIP(domain)
	if err != nil {
		return "", err
	}
	for _, ip := range ips {
		if ip.To4() != nil {
			return ip.String(), nil
		}
	}
	// fallback ke IPv6 kalau IPv4 nggak ada
	if len(ips) > 0 {
		return ips[0].String(), nil
	}
	return "", fmt.Errorf("tidak ditemukan IP")
}

// Fetch info dari ipinfo.io
func fetchIPInfo(ip string) (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(fmt.Sprintf(IPINFO_URL, ip))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// Handler untuk /ipinfo
func HandleIPInfo(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	var domainText string

	if update.Message.CommandArguments() != "" {
		domainText = update.Message.CommandArguments()
	} else if update.Message.ReplyToMessage != nil && update.Message.ReplyToMessage.Text != "" {
		domainText = strings.TrimSpace(update.Message.ReplyToMessage.Text)
	} else {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			"/ipinfo example.com\nAtau reply pesan yang berisi domain/URL")
		bot.Send(msg)
		return
	}

	domains := extractDomains(domainText)
	if len(domains) == 0 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Tidak ada domain valid ditemukan.")
		bot.Send(msg)
		return
	}

	loadingMsg, _ := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Mengambil data..."))

	var results []string
	for idx, domain := range domains {
		ip, err := resolveDomain(domain)
		if err != nil {
			results = append(results, fmt.Sprintf("**%s**: gagal resolusi domain", domain))
			continue
		}

		info, err := fetchIPInfo(ip)
		if err != nil {
			results = append(results, fmt.Sprintf("**%s** (%s): gagal ambil data dari ipinfo.io", domain, ip))
			continue
		}

		results = append(results, fmt.Sprintf("%s (%s):\n```\n%s\n```", domain, ip, info))

		editMsg := tgbotapi.NewEditMessageText(update.Message.Chat.ID, loadingMsg.MessageID,
			fmt.Sprintf("Memproses: `%s` (%d/%d)", domain, idx+1, len(domains)))
		editMsg.ParseMode = "Markdown"
		bot.Send(editMsg)
	}

	finalOutput := strings.Join(results, "\n\n")
	if len(finalOutput) > 4000 {
		finalOutput = finalOutput[:4000] + "..."
	}

	editMsg := tgbotapi.NewEditMessageText(update.Message.Chat.ID, loadingMsg.MessageID, finalOutput)
	editMsg.ParseMode = "Markdown"
	bot.Send(editMsg)
}
