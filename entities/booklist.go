package entities
type BookList struct{
	Book `json:"Book"`
	LenderName string `json:"addedBy"`
	LenderRating float32 `json:"userRating"`
}
func (BookList) TableName() string {
    return "BookList"
}