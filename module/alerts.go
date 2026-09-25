package module

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// StartHealthWatcher is kept as a compatibility no-op. Automatic alerts are disabled;
// health checks remain available through the manual /health command.
func StartHealthWatcher(_ *tgbotapi.BotAPI, _ int64, _ time.Duration) {}
