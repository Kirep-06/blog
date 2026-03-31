package seed

import (
	"log"

	"blog/config"
	"blog/internal/database"
	"blog/internal/model"

	"golang.org/x/crypto/bcrypt"
)

func Run() error {
	cfg := config.C.Seed

	var existing model.User
	result := database.DB.Where("username = ?", cfg.AdminUsername).First(&existing)
	if result.Error == nil {
		log.Printf("seed: user %s already exists, skipping", cfg.AdminUsername)
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(cfg.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := model.User{
		Username:     cfg.AdminUsername,
		PasswordHash: string(hash),
	}
	if err := database.DB.Create(&user).Error; err != nil {
		return err
	}

	log.Printf("seed: user %s created successfully", cfg.AdminUsername)
	return nil
}
