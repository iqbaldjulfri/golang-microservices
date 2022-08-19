package services

import (
	"errors"
	"gateway/models"
	"gateway/services/db"
	"gateway/utils"
	"log"
)

func GetUserByUsername(username string) (*models.User, error) {
	user := new(models.User)
	result := db.Conn.Where("username = ?", username).First(user)

	if result.Error != nil {
		log.Println(result.Error.Error())
		return nil, errors.New("No user found")
	}

	return user, nil
}

func CreateUser(dto models.CreateUserDto) (*models.User, error) {
	user := models.User{Username: dto.Username, Password: utils.HashPassword(dto.Password), Email: dto.Email}

	result := db.Conn.Create(&user)
	if result.Error != nil {
		err := result.Error.Error()
		log.Println(err)
		switch err {
		case "UNIQUE constraint failed: users.username":
			return nil, errors.New("Username already exists")
		case "UNIQUE constraint failed: users.email":
			return nil, errors.New("Email already exists")
		default:
			return nil, errors.New("Error writing to database")
		}
	}

	return &user, nil
}

func FindUsers(username string) ([]models.User, error) {
	var users []models.User
	result := db.Conn.Where("username LIKE ?", "%"+username+"%").Find(&users)

	if result.Error != nil {
		log.Println(result.Error.Error())
		return nil, errors.New("Error when reading database")
	}

	return users, nil
}

func UpdateUser(id uint, dto models.UpdateUserDto) (*models.User, error) {
	var user = new(models.User)
	upData := models.User{Email: dto.Email, IsActive: dto.IsActive}
	if dto.Password != "" {
		upData.Password = utils.HashPassword(dto.Password)
	}
	result := db.Conn.Model(&user).Updates(upData)

	if result.Error != nil {
		log.Println(result.Error.Error())
		return nil, errors.New("Error when writing database")
	}

	return user, nil
}
