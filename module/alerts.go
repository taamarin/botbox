package module

import (
	"fmt"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	healthAlertLock sync.Mutex
	healthAlertLast string
)

func collectHealthIssues() string {
	issues := make([]string, 0, 8)
	for _, name := range []string{"mihomo", "clash", "sing-box", "xray", "hysteria"} {
		if !isProcessRunning(name) {
			issues = append(issues, fmt.Sprintf("• %s: not running", strings.Title(name)))
		}
	}

	var version map[string]string
	if err := getJSON(mihomoAPI+"/version", &version); err != nil {
		issues = append(issues, "• Mihomo API: unreachable")
	}

	if len(issues) == 0 {
		return ""
	}
	return "🚨 *Health Alert*\n\n" + strings.Join(issues, "\n")
}

func StartHealthWatcher(bot *tgbotapi.BotAPI, ownerID int64, interval time.Duration) {
	if bot == nil || ownerID == 0 {
		return
	}
	if interval <= 0 {
		interval = 60 * time.Second
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			issueText := collectHealthIssues()
			healthAlertLock.Lock()
			if issueText != "" && issueText != healthAlertLast {
				bot.Send(tgbotapi.NewMessage(ownerID, issueText))
				healthAlertLast = issueText
			}
			if issueText == "" {
				healthAlertLast = ""
			}
			healthAlertLock.Unlock()
		}
	}()
}
