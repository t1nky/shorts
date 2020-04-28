package models

// User structure
type User struct {
	ID       uint64 `json:"id" gorm:"primary_key"`
	Name     string `json:"name" gorm:"unique;not null" binding:"required,min=5"`
	Password string `json:"password" gorm:"not null" binding:"required,min=5,max=16"`
}

// AddUserData structure
// swagger:parameters addUser
type AddUserData struct {
	Name     string `json:"name" gorm:"unique;not null" binding:"required,min=5"`
	Password string `json:"password" gorm:"not null" binding:"required,min=5,max=16"`
}
