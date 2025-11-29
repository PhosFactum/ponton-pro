package handlers

import (
	"github.com/PhosFactum/TechnoLotos/backend/internal/bot"
	"github.com/PhosFactum/TechnoLotos/backend/internal/config"
)

// Handler - структура для хранения зависимостей
type Handler struct {
	bot *bot.Bot
	cfg *config.Config
}

// NewHandler - конструктор хэндлера
func NewHandler(b *bot.Bot, cfg *config.Config) *Handler {
	return &Handler{
		bot: b,
		cfg: cfg,
	}
}
