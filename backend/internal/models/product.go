// internal/models/service.go:
package models

// Product - модель продукта (понтон, лодка)
type Product struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
