package models

type User struct {
	UserID     int64  `json:"user_id"`
	Phone      string `json:"phone"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patrynomic string `json:"patrynomic"`
	Mail       string `json:"mail"`
	Hash       string `json:"hash"`
	Photo      string `json:"userphoto"`
	Role       Role   `json:"Role"`
}

type Role struct {
	RoleID int    `json:"role_id"`
	Role   string `json:"role"`
}
