package entities
type Book struct{
	ID string `json:"bookId" gorm:"primaryKey;autoIncrement:false"`
	BookName string `json:"bookname"`
	Author string `json:"author"`
	Genre string `json:"genre"`
	Rating float32 `json:"bookrating"`
	Status string `json:"available"`
}
