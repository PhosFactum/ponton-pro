package bot

import (
	"fmt"
	"log"
	"strconv"

	"github.com/PhosFactum/TechnoLotos/backend/internal/database"
	"github.com/PhosFactum/TechnoLotos/backend/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const requestsPerPage = 5

// Ручка /start
func (h *Handlers) handleStart(message *tgbotapi.Message) {
	text := `**🤖 Бот-помощник 'ТехноЛотоса'**

Доступные команды:
/start - Вывести это руководство ещё раз
/test - Проверить работу бота
/requests - Показать последние заявки

Для навигации по заявкам используйте кнопки 'Вперёд' и 'Назад'.`

	h.bot.SendMessage(message.Chat.ID, text)
}

// Ручка /test
func (h *Handlers) handleTest(message *tgbotapi.Message) {
	text := "✅ Бот работает исправно, так как вы видите это сообщение!"
	h.bot.SendMessage(message.Chat.ID, text)
}

// Ручка /requests
func (h *Handlers) handleRequests(message *tgbotapi.Message, page int) {
	var requests []models.Request
	var total int64

	// Получаем общее число заявок
	database.DB.Model(&models.Request{}).Count(&total)

	// Получаем заявки на текущую страницу (для пагинации)
	offset := page * requestsPerPage
	result := database.DB.Preload("Product").
		Order("created_at DESC").
		Offset(offset).
		Limit(requestsPerPage).
		Find(&requests)

	if result.Error != nil {
		log.Printf("Error fetching requests: %v", result.Error)
		h.bot.SendMessage(message.Chat.ID, "Ошибка при получении заявок из базы данных!")
		return
	}

	if len(requests) == 0 {
		text := "Заявок пока нет."
		h.bot.SendMessage(message.Chat.ID, text)
		return
	}

	// Формируем сообщение
	text := fmt.Sprintf("**Заявки** (страница %d)\n\n", page+1)

	for i, req := range requests {
		productName := "Не указан"
		if req.Product != nil {
			productName = req.Product.Title
		}

		text += fmt.Sprintf("**Заявка** #%d**\n", offset+i+1)
		text += fmt.Sprintf("👤 Имя: %s\n", req.Name)
		text += fmt.Sprintf("📞 Телефон: %s\n", req.Phone)

		if req.Email != "" {
			text += fmt.Sprintf("📧 Email: %s\n", req.Email)
		}
		text += fmt.Sprintf("📝 Почта: %s\n", req.Email)

		if req.Description != "" {
			text += fmt.Sprintf("📝 Описание: %s\n", req.Description)
		}

		text += fmt.Sprintf("🛍️ Товар: %s\n", productName)
		text += fmt.Sprintf("🕐 Дата: %s\n\n", req.CreatedAt.Format("02.01.2006 15:04"))
	}

	totalPages := (int(total) + requestsPerPage - 1) / requestsPerPage
	if totalPages == 0 {
		totalPages = 1
	}

	text += fmt.Sprintf("Страница %d из %d", page+1, totalPages)

	// Создаём клавиатуру для пагинации
	keyboard := h.createPaginationKeyboard(page, int(total))

	h.bot.SendMessageWithKeyboard(message.Chat.ID, text, keyboard)
}

// Обработчик коллбэка для пагинации
func (h *Handlers) handleRequestsCallback(query *tgbotapi.CallbackQuery) {
	data := query.Data
	pageStr := data[14:] // "requests_page_123"
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		h.bot.SendMessage(query.Message.Chat.ID, "❌ Ошибка перемотки страниц")
		return
	}

	// Используем уже существующее сообщение для обновления
	h.handleRequests(query.Message, page)
}

// Обработчик неизвестной команды
func (h *Handlers) handleUnknown(message *tgbotapi.Message) {
	text := "❌ Неизвестная команда. Используйте /start для просмотра доступных команд."
	h.bot.SendMessage(message.Chat.ID, text)
}

// Обработчик создания клавиатуры
func (h *Handlers) createPaginationKeyboard(currentPage int, totalItems int) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	maxPage := (totalItems + requestsPerPage - 1) / requestsPerPage
	if maxPage <= 1 {
		return tgbotapi.InlineKeyboardMarkup{}
	}

	var buttons []tgbotapi.InlineKeyboardButton

	// Кнопка "Назад"
	if currentPage > 0 {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад",
			fmt.Sprintf("requests_page_%d", currentPage-1)))
	}

	// Кнопка "Вперёд"
	if currentPage < maxPage-1 {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData("Вперёд ➡️",
			fmt.Sprintf("requests_page_%d", currentPage-1)))
	}

	if len(buttons) > 0 {
		rows = append(rows, buttons)
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
