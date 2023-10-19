package initializers

import "jwt/models"

func SyncDatabase() {
	//миграция модели User
	DB.AutoMigrate(&models.User{})
}
