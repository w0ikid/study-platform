package models

import "time"

type User struct {
    ID        int      `json:"id"`
    Username  string    `json:"username"`
    Name      string    `json:"name"`
    Surname   string    `json:"surname"`
    Email     string    `json:"email"`
    Password  string    `json:"password"`
    Role      string    `json:"role"` // roles: 0 - student, 1 - teacher, 2 - admin
    Level     int       `json:"level"` // 0 - beginner, 1 - intermediate, 2 - advanced 
    Xp        int       `json:"xp"`    // experience points
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
    Level    int    `json:"level"`
    Xp       int    `json:"xp"`
    CreatedAt time.Time `json:"created_at"`
}
