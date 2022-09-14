package models

type User struct {
	UserID      int64  `json:"user_id"`
	Phone       string `json:"phone"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patrynomic  string `json:"patrynomic"`
	Mail        string `json:"mail"`
	Hash        string `json:"hash"`
	Description string `json:"description"`
	Photo       string `json:"photo"`
	Role        Role   `json:"Role"`
}

type Role struct {
	RoleID int    `json:"role_id"`
	Role   string `json:"role"`
}

type UserAuth struct {
	UserPhone string `json:"username"`
	Password  string `json:"password"`
}

type RefreshToken struct {
	ID           int64  `json:"id"`
	UserID       int64  `json:"userid"`
	ExpiresIn    int64  `json:"expires_in"`
	UserAgent    string `json:"useragent"`
	RefreshToken string `json:"refresh_token"`
}
