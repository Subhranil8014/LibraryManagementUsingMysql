package entities
type LendingRecord struct{
	LenderName string `json:"name" gorm:"primaryKey;autoIncrement:false"`
	Token int `json:"token"`
}
func (LendingRecord) TableName() string {
    return "LenderTokenUpdate"
}