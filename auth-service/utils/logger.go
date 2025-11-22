package utils

import (
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Log глобальная переменная логгера
var Log = logrus.New()

// InitLogger инициализирует логгер
func InitLogger() {
	// Настройка формата
	Log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Настройка уровня логирования
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	Log.SetLevel(level)

	// Настройка вывода в зависимости от среды
	env := os.Getenv("APP_ENV")
	if env == "production" {
		SetupLogRotation()
	} else {
		// В development выводим в консоль
		Log.SetOutput(os.Stdout)
		// В development используем текстовый формат для лучшей читаемости
		Log.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
		})
	}
}

// SetupLogRotation настраивает ротацию логов
func SetupLogRotation() {
	logRotation := &lumberjack.Logger{
		Filename:   "/var/log/auth-service.log",
		MaxSize:    100, // MB
		MaxBackups: 3,
		MaxAge:     28, // days
		Compress:   true,
	}

	Log.SetOutput(logRotation)
}
