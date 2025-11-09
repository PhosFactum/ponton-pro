package models

// Request - заявка от клиента
type Request struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"not null" json:"name"`
	Phone       string `gorm:"not null" json:"phone"`
	Email       string `json:"email,omitempty"`
	Description string `json:"description,omitempty"`
	ProductID   *uint  `json:"product_id,omitempty"`

	Product *Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}
