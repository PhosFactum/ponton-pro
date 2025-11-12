package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api      *tgbotapi.BotAPI
	handlers *Handlers
}

type Handlers struct {
	bot *Bot
}

func NewBot(token string) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	bot := &Bot{
		api: api,
	}

	bot.handlers = &Handlers{bot: bot}

	log.Printf("Authorized on account %s", api.Self.UserName)
	return bot, nil
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			b.handleMessage(update.Message)
		} else if update.CallbackQuery != nil {
			b.HandleCallback(update.CallbackQuery)
		}
	}
}

func (b *Bot) handleMessage(message *tgbotapi.Message) {
	switch message.Command() { // Тут прописывать команды боту
	case "start":
		b.handlers.handleStart(message)
	case "test":
		b.handlers.handleTest(message)
	case "requests":
		b.handlers.handleRequests(message, 0) // С первой странички
	default:
		b.handlers.handleUnknown(message)
	}
}

// Просто отправка текста в чат
func (b *Bot) SendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := b.api.Send(msg)
	if err != nil {
		log.Printf("Error sending message to chat %d: %v", chatID, err)
	}
}

// Отправка текста с клавиатурой
func (b *Bot) SendMessageWithKeyboard(chatID int64, text string, keyboard tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	_, err := b.api.Send(msg)
	if err != nil {
		log.Printf("Error sending message with keyboard to chat %d: %v", chatID, err)
	}
}

// Обработчик коллбэков для пагинации
func (b *Bot) HandleCallback(query *tgbotapi.CallbackQuery) {
	data := query.Data
	log.Printf("Callback received: %s", data)

	// Обработка пагинации заявок
	if len(data) > 14 && data[:14] == "requests_page_" {
		b.handlers.handleRequestsCallback(query)
	}

	// Ответ на коллбэк (убираем часы)
	callback := tgbotapi.NewCallback(query.ID, "")
	b.api.Request(callback)
}
