package config

import "gorm.io/gorm"

var (
	db     *gorm.DB
	logger *Logger
)

func Init() error {
	return nil
}

func GetLooger(p string) *Logger {
	//Initialize logger
	logger = NewLogger(p)
	return logger
}
