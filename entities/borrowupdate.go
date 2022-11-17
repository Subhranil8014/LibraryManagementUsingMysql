package entities
type BorrowUpdate struct{
	BorrowerName string `json:"name" gorm:"primaryKey;autoIncrement:false"`
	Token int `json:"token"`
}
func (BorrowUpdate) TableName() string {
    return "BorrowTokenUpdate"
}