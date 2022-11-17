package entities


type BorrowRecord struct {
	ID           uint    `gorm:"primaryKey"`
	BookID       string  `json:"bookId"`
	BookName     string  `json:"bookname"`
	Author       string  `json:"author"`
	Genre        string  `json:"genre"`
	Rating       float32 `json:"bookrating"`
	Status       string  `json:"available"`
	BorrowerName string  `json:"name"`
}

func (BorrowRecord) TableName() string {
	return "BorrowerRecord"
}