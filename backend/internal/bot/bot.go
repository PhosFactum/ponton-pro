// AI-helped code
package bot

import (
	"fmt"
	"log"

	"github.com/PhosFactum/TechnoLotos/backend/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api      *tgbotapi.BotAPI
	handlers *Handlers
}

type Handlers struct {
	bot *Bot
}

// NewBot - Экземпляр нового бота
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

// Start - меню с кнопками
func (b *Bot) Start() {
	b.setBotCommands()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			if update.CallbackQuery != nil {
				b.HandleCallback(update.CallbackQuery)
			}
			continue
		}

		// Логируем сообщения
		log.Printf("[TG MSG] From: %s | Text: %s", update.Message.From.UserName, update.Message.Text)

		b.handleMessage(update.Message)
	}
}

// setBotCommands - настройка выпадающего меню команд
func (b *Bot) setBotCommands() {
	commands := []tgbotapi.BotCommand{
		{Command: "start", Description: "Главное меню / Помощь"},
		{Command: "requests", Description: "Список заявок"},
	}

	cfg := tgbotapi.NewSetMyCommands(commands...)
	if _, err := b.api.Request(cfg); err != nil {
		log.Printf("Ошибка настройки команд меню: %v", err)
	}
}

// handleMessage - обработка сообщения через / и с кнопок
func (b *Bot) handleMessage(message *tgbotapi.Message) {
	// 1. Обработка команд (через слэш /)
	if message.IsCommand() {
		switch message.Command() {
		case "start":
			b.handlers.handleStart(message)
		case "test":
			b.handlers.handleTest(message)
		case "requests":
			b.handlers.handleRequests(message, 0)
		default:
			b.handlers.handleUnknown(message)
		}
		return
	}

	// 2. Обработка НИЖНИХ КНОПОК (Reply Keyboard)
	switch message.Text {
	case "📋 Заявки":
		b.handlers.handleRequests(message, 0)
	case "❓ Помощь":
		// Кнопка Помощь делает то же самое, что и /start
		b.handlers.handleStart(message)
	default:
	}
}

// SendMessage - простая отправка
func (b *Bot) SendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	_, err := b.api.Send(msg)
	if err != nil {
		log.Printf("Error sending message to chat %d: %v", chatID, err)
	}
}

// SendMessageWithKeyboard - отправка с INLINE клавиатурой (под сообщением)
func (b *Bot) SendMessageWithKeyboard(chatID int64, text string, keyboard tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	msg.ParseMode = "Markdown"
	_, err := b.api.Send(msg)
	if err != nil {
		log.Printf("Error sending message with keyboard to chat %d: %v", chatID, err)
	}
}

// SendMessageWithReplyKeyboard - отправка с НИЖНЕЙ клавиатурой (кнопки меню)
func (b *Bot) SendMessageWithReplyKeyboard(chatID int64, text string, keyboard tgbotapi.ReplyKeyboardMarkup) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	msg.ParseMode = "Markdown"
	_, err := b.api.Send(msg)
	if err != nil {
		log.Printf("Error sending message with reply keyboard: %v", err)
	}
}

// EditMessage - редактирование сообщения (для пагинации)
func (b *Bot) EditMessage(chatID int64, messageID int, text string, keyboard *tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	msg.ParseMode = "Markdown"
	if keyboard != nil {
		msg.ReplyMarkup = keyboard
	}
	_, err := b.api.Send(msg)
	if err != nil {
		log.Printf("Error editing message: %v", err)
	}
}

// HandleCallback - обработка нажатий кнопок "Вперед/Назад"
func (b *Bot) HandleCallback(query *tgbotapi.CallbackQuery) {
	data := query.Data
	if len(data) > 14 && data[:14] == "requests_page_" {
		b.handlers.handleRequestsCallback(query)
	}
	callback := tgbotapi.NewCallback(query.ID, "")
	b.api.Request(callback)
}

// SendNewRequestNotification - уведомление о новом заказе
func (b *Bot) SendNewRequestNotification(chatID int64, req models.Request) {
	if chatID == 0 {
		return
	}
	productName := "Не выбран"
	if req.Product != nil {
		productName = req.Product.Title
	}
	text := fmt.Sprintf(
		"🔥 **НОВАЯ ЗАЯВКА!** 🔥\n\n"+
			"👤 **Имя:** %s\n"+
			"📞 **Телефон:** `%s`\n"+
			"📧 **Email:** %s\n"+
			"🛍️ **Товар:** %s\n"+
			"📝 **Комментарий:** %s\n\n"+
			"#заявка #new",
		req.Name, req.Phone, req.Email, productName, req.Description,
	)
	b.SendMessage(chatID, text)
}
