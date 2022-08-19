package models

type User struct {
	Model
	Username string `gorm:"unique;not null;uniqueIndex"`
	Password string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
	IsActive bool   `gorm:"not null;default:true"`
	Roles    []Role `gorm:"many2many:user_roles"`
}

type UserSafeDto struct {
	Model
	Username string `validate:"required" json:"username"`
	Email    string `validate:"required" json:"email"`
	IsActive bool   `json:"isActive"`
	Roles    []Role `json:"roles"`
}

type CreateUserDto struct {
	Username       string `validate:"required,printascii,min=5,max=20" json:"username"`
	Password       string `validate:"required,min=5" json:"password"`
	RepeatPassword string `validate:"required,eqfield=Password" json:"repeatPassword"`
	Email          string `validate:"required,email" json:"email"`
}

type UpdateUserDto struct {
	Password       string `validate:"min=5"`
	RepeatPassword string `validate:"required_with:Password,eqfield=Password"`
	Email          string `validate:"email"`
	IsActive       bool
	Roles          []uint
}

func ToUserSafeDto(user User) *UserSafeDto {
	return &UserSafeDto{
		Username: user.Username,
		Email:    user.Email,
		IsActive: user.IsActive,
		Roles:    user.Roles,
		Model:    user.Model,
	}
}
