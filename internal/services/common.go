package services

type OrderItem struct {
	ID    uint `json:"id" binding:"required"`
	Order int  `json:"order" binding:"required"`
}
