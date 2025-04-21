package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("Файл .env не найден")
	}
}

const csvFileName = "links.csv"

func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")

	if token == "" {
		log.Fatal("Неверный или отсутствующий токен")
		return
	}
	// Инициализация бота
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	log.Printf("Авторизация бота %s", bot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates, _ := bot.GetUpdatesChan(updateConfig)

	// Создание кнопок
	addButton := tgbotapi.NewKeyboardButton("Добавить ссылку")
	getButton := tgbotapi.NewKeyboardButton("Получить данные по категории")
	keyboard := tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(addButton, getButton))
	keyboard.ResizeKeyboard = true

	for update := range updates {
		if update.Message == nil {
			continue
		}
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		text := update.Message.Text
		chatID := update.Message.Chat.ID

		// Отправка клавиатуры
		if text == "/start" {
			msg := tgbotapi.NewMessage(chatID, "Добро пожаловать! Выберите действие:")
			msg.ReplyMarkup = keyboard
			bot.Send(msg)
			continue
		}

		switch text {
		case "Добавить ссылку":
			msg := tgbotapi.NewMessage(chatID, "Введите данные в формате: \n<категория>\n<ссылка>\n<описание>")
			bot.Send(msg)
			nextUpdate := <-updates
			if nextUpdate.Message != nil {
				text := nextUpdate.Message.Text
				chatID := nextUpdate.Message.Chat.ID

				parts := strings.Split(text, "\n")
				if len(parts) < 3 {
					msg := tgbotapi.NewMessage(chatID, "Неверный формат. Используйте:\n<категория>\n<ссылка>\n<описание>")
					bot.Send(msg)
					continue
				}

				category := strings.TrimSpace(parts[0])
				url := strings.TrimSpace(parts[1])
				description := strings.TrimSpace(parts[2])

				if err := saveToCSV(category, url, description); err != nil {
					msg := tgbotapi.NewMessage(chatID, "Ошибка при сохранении данных.")
					bot.Send(msg)
					continue
				}

				msg := tgbotapi.NewMessage(chatID, "Данные успешно сохранены!")
				bot.Send(msg)
			}
		case "Получить данные по категории":
			msg := tgbotapi.NewMessage(chatID, "Введите категорию")
			bot.Send(msg)
			nextUpdate := <-updates
			if nextUpdate.Message != nil {
				category := nextUpdate.Message.Text
				chatID := nextUpdate.Message.Chat.ID

				result, err := readFromCSV(category)
				if err != nil {
					msg := tgbotapi.NewMessage(chatID, "Ошибка при чтении данных.")
					bot.Send(msg)
					continue
				}

				msg := tgbotapi.NewMessage(chatID, result)
				bot.Send(msg)
			}
		default:
			msg := tgbotapi.NewMessage(chatID, "Используйте кнопки для добавления или получения данных.")
			bot.Send(msg)
		}
	}
}

// Функция для записи данных в CSV файл
func saveToCSV(category, url, description string) error {
	file, err := os.OpenFile(csvFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	record := []string{category, url, description}
	return writer.Write(record)
}

// Функция для чтения данных из CSV файла
func readFromCSV(category string) (string, error) {
	file, err := os.Open(csvFileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return "", err
	}

	var result strings.Builder
	for _, record := range records {
		if strings.Contains(strings.ToLower(record[0]), strings.ToLower(category)) {
			result.WriteString(fmt.Sprintf("Категория - %s\n%s\n\n%s\n\n", record[0], record[2], record[1]))
		}
	}

	if result.Len() == 0 {
		return "Нет записей для этой категории", nil
	}
	return result.String(), nil
}
