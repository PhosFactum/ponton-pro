// AI-helped code
package bot

import (
	"fmt"
	"log"
	"strconv"

	"github.com/PhosFactum/TechnoLotos/backend/internal/database"
	"github.com/PhosFactum/TechnoLotos/backend/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Количество заявок на странице
const requestsPerPage = 3

// Ручка /start
func (h *Handlers) handleStart(message *tgbotapi.Message) {
	text := `**🤖 Приветствую! Я бот "ТехноЛотоса".**

Я буду присылать сюда новые заявки с сайта.
Нажмите кнопку внизу, чтобы открыть список заявок.`

	// Создаем НИЖНЮЮ клавиатуру (Reply Keyboard)
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📋 Заявки"),
			tgbotapi.NewKeyboardButton("❓ Помощь"),
		),
	)
	keyboard.ResizeKeyboard = true // Делаем кнопки компактными

	// Отправляем сообщение + клавиатуру
	h.bot.SendMessageWithReplyKeyboard(message.Chat.ID, text, keyboard)
}

// Ручка /test
func (h *Handlers) handleTest(message *tgbotapi.Message) {
	h.bot.SendMessage(message.Chat.ID, "✅ Бот работает исправно!")
}

// Ручка /requests (вызывается командой или кнопкой "Заявки")
func (h *Handlers) handleRequests(message *tgbotapi.Message, page int) {
	// Это новое сообщение, не редактирование
	h.sendOrEditRequests(message.Chat.ID, 0, page, false)
}

// handleRequestsCallback - пагинация по списку заявок
func (h *Handlers) handleRequestsCallback(query *tgbotapi.CallbackQuery) {
	data := query.Data
	pageStr := data[14:] // "requests_page_123"
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		h.bot.SendMessage(query.Message.Chat.ID, "❌ Ошибка перемотки страниц")
		return
	}

	// Это редактирование старого сообщения
	h.sendOrEditRequests(query.Message.Chat.ID, query.Message.MessageID, page, true)
}

// sendOrEditRequests - вывод списка заявок
func (h *Handlers) sendOrEditRequests(chatID int64, messageID int, page int, isEdit bool) {
	var requests []models.Request
	var total int64

	database.DB.Model(&models.Request{}).Count(&total)

	offset := page * requestsPerPage
	result := database.DB.Preload("Product").
		Order("created_at DESC").
		Offset(offset).
		Limit(requestsPerPage).
		Find(&requests)

	if result.Error != nil {
		log.Printf("Error fetching requests: %v", result.Error)
		h.bot.SendMessage(chatID, "Ошибка при получении заявок из базы данных!")
		return
	}

	if len(requests) == 0 && !isEdit {
		h.bot.SendMessage(chatID, "Заявок пока нет.")
		return
	}

	text := fmt.Sprintf("**📋 Список заявок** (Стр. %d)\n\n", page+1)

	for _, req := range requests {
		productName := "Не указан"
		if req.Product != nil {
			productName = req.Product.Title
		}

		text += fmt.Sprintf("**Заявка #%d**\n", req.ID)
		text += fmt.Sprintf("👤 Имя: %s\n", req.Name)
		text += fmt.Sprintf("📞 Тел: `%s`\n", req.Phone)

		if req.Email != "" {
			text += fmt.Sprintf("📧 Email: %s\n", req.Email)
		}

		if req.Description != "" {
			text += fmt.Sprintf("📝 Инфо: %s\n", req.Description)
		}

		text += fmt.Sprintf("🛍️ Товар: %s\n", productName)
		text += fmt.Sprintf("🕐 Дата: %s\n\n", req.CreatedAt.Format("02.01.2006 15:04"))
	}

	totalPages := (int(total) + requestsPerPage - 1) / requestsPerPage
	if totalPages == 0 {
		totalPages = 1
	}

	text += fmt.Sprintf("Страница %d из %d", page+1, totalPages)

	var keyboard *tgbotapi.InlineKeyboardMarkup
	if totalPages > 1 {
		kb := h.createPaginationKeyboard(page, int(total))
		keyboard = &kb
	}

	if isEdit {
		h.bot.EditMessage(chatID, messageID, text, keyboard)
	} else {
		if keyboard != nil {
			h.bot.SendMessageWithKeyboard(chatID, text, *keyboard)
		} else {
			h.bot.SendMessage(chatID, text)
		}
	}
}

// handleUnknown - Обработка неизвестной команды
func (h *Handlers) handleUnknown(message *tgbotapi.Message) {
	h.bot.SendMessage(message.Chat.ID, "❌ Неизвестная команда. Нажмите '❓ Помощь' для меню.")
}

// createPaginationKeyboard - создание клавиатуры для пагинации
func (h *Handlers) createPaginationKeyboard(currentPage int, totalItems int) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	maxPage := (totalItems + requestsPerPage - 1) / requestsPerPage

	var buttons []tgbotapi.InlineKeyboardButton

	if currentPage > 0 {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад",
			fmt.Sprintf("requests_page_%d", currentPage-1)))
	}

	if currentPage < maxPage-1 {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData("Вперёд ➡️",
			fmt.Sprintf("requests_page_%d", currentPage+1)))
	}

	if len(buttons) > 0 {
		rows = append(rows, buttons)
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
