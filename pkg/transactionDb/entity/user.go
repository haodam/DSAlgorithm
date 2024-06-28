package entity

type User struct {
	ID    *string `json:"id,omitempty"`
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty" gorm:"unique"`
	Phone *string `json:"phone,omitempty" gorm:"index:idx_id;unique"`
}

func (User) TableName() string {
	return "users"
}
