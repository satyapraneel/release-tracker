package repositories

import (
	"github.com/release-trackers/gin/models"
	"golang.org/x/crypto/bcrypt"
)

func PasswordHash(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return ""
	}
	return string(hashedPassword)
}

func (app *App) AuthenticateUser(email string, password string) (models.Users, bool) {
	user := models.Users{}
	app.Db.Where("email = ?", email).First(&user)
	hashedPassword := []byte(user.Password)
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		return models.Users{}, true
	}
	return user, false
}
