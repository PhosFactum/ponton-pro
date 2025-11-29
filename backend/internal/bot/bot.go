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
			if update.CallbackQuery != nil {
				b.HandleCallback(update.CallbackQuery)
			}
			continue
		}

		// Логируем ID чата, чтобы узнать его
		log.Printf("[%s] ID чата: %d | Текст: %s", update.Message.From.UserName,
			update.Message.Chat.ID, update.Message.Text)

		b.handleMessage(update.Message)
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

// SendNessage - базовая отправка текста в чат
func (b *Bot) SendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	_, err := b.api.Send(msg)
	if err != nil {
		log.Printf("Error sending message to chat %d: %v", chatID, err)
	}
}

// EditMessage - изменение текста и клавиатуры для существующего сообщения
func (b *Bot) EditMessage(chatID int64, messageID int, text string,
	keyboard *tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	msg.ParseMode = "Markdown"

	if keyboard != nil {
		msg.ReplyMarkup = keyboard
	}

	_, err := b.api.Send(msg)
	if err != nil {
		log.Printf("Error editing message %d in chat %d: %v", messageID, chatID, err)
	}
}

// SendMessageWithKeyboard - отправка с кнопками
func (b *Bot) SendMessageWithKeyboard(chatID int64, text string, keyboard tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	msg.ParseMode = "Markdown"
	_, err := b.api.Send(msg)
	if err != nil {
		log.Printf("Error sending message with keyboard to chat %d: %v", chatID, err)
	}
}

// Обработчик коллбэков для пагинации
func (b *Bot) HandleCallback(query *tgbotapi.CallbackQuery) {
	data := query.Data
	// log.Printf("Callback received: %s", data)   // Для отладки

	// Обработка пагинации заявок
	if len(data) > 14 && data[:14] == "requests_page_" {
		b.handlers.handleRequestsCallback(query)
	}

	callback := tgbotapi.NewCallback(query.ID, "")
	b.api.Request(callback)
}

// SendNewRequestNotification - отправка уведомления о новом заказе
func (b *Bot) SendNewRequestNotification(chatID int64, req models.Request) {
	if chatID == 0 {
		log.Printf("AdminChatID не выбран, уведомление не будет отправлено!")
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
		req.Name,
		req.Phone,
		req.Email,
		productName,
		req.Description,
	)

	b.SendMessage(chatID, text)
}
