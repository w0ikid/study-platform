package dto

// CreateUserInput represents the input for creating a user.
// swagger:model
type CreateUserInput struct {
    // Username of the user
    // required: true
    Username string `json:"username" validate:"required"`
    
    // First name of the user
    Name string `json:"name"`
    
    // Last name of the user
    Surname string `json:"surname"`
    
    // Email address
    // required: true
    Email string `json:"email" validate:"required,email"`
    
    // Password
    // required: true
    Password string `json:"password" validate:"required,min=8"`
    
    // Role of the user (student, teacher, admin)
    // required: true
    Role string `json:"role" validate:"required,oneof=student teacher admin"`
}

// swagger:model
type LoginUserInput struct {
    // User email
    // required: true
    Email string `json:"email" binding:"required,email"`

    // User password
    // required: true
    Password string `json:"password" binding:"required,min=6"`
}
