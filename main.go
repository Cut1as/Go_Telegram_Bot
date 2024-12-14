package main

import (
	"fmt"
	"os"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func main() {
	// Токен бота
	botToken := "7611078519:AAFJeJltikMhjsp74ezKMWui7TUn_PqUaV0"

	// Создание экземпляра бота
	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println("Ошибка при создании бота:", err)
		os.Exit(1)
	}

	// Создание кнопки с командой /start в меню
	err = bot.SetMyCommands(&telego.SetMyCommandsParams{
		Commands: []telego.BotCommand{
			{
				Command:     "start",
				Description: "Начать работу с ботом",
			},
		},
	})
	if err != nil {
		fmt.Println("Ошибка при добавлении команды:", err)
		os.Exit(1)
	}

	// Получение обновлений
	updates, _ := bot.UpdatesViaLongPolling(nil)

	defer bot.StopLongPolling()

	// Карта кнопок и соответствующих файлов
	files := map[string]string{
		"slice":      "slice.txt",
		"map":        "map.txt",
		"goroutines": "goroutines.txt",
		"channels":   "channels.txt",
	}

	// Кнопка "Назад в меню"
	backButton := telego.InlineKeyboardButton{
		Text:         "Назад в меню",
		CallbackData: "back_to_menu",
	}

	// Цикл обработки обновлений
	for update := range updates {
		// Обработка команды /start
		if update.Message != nil && update.Message.Text == "/start" {
			chatID := tu.ID(update.Message.Chat.ID)

			// Отправка кнопок выбора темы
			_, _ = bot.SendMessage(&telego.SendMessageParams{
				ChatID: chatID,
				Text:   "Выберите тему, чтобы получить информацию:",
				ReplyMarkup: &telego.InlineKeyboardMarkup{
					InlineKeyboard: [][]telego.InlineKeyboardButton{
						{
							{Text: "Slice", CallbackData: "slice"},
							{Text: "Map", CallbackData: "map"},
						},
						{
							{Text: "Goroutines", CallbackData: "goroutines"},
							{Text: "Channels", CallbackData: "channels"},
						},
					},
				},
			})
		}

		// Обработка нажатий на кнопки
		if update.CallbackQuery != nil {
			chatID := tu.ID(update.CallbackQuery.From.ID)
			data := update.CallbackQuery.Data

			// Если нажата кнопка "Назад в меню"
			if data == "back_to_menu" {
				_, _ = bot.SendMessage(&telego.SendMessageParams{
					ChatID: chatID,
					Text:   "Выберите тему, чтобы получить информацию:",
					ReplyMarkup: &telego.InlineKeyboardMarkup{
						InlineKeyboard: [][]telego.InlineKeyboardButton{
							{
								{Text: "Slice", CallbackData: "slice"},
								{Text: "Map", CallbackData: "map"},
							},
							{
								{Text: "Goroutines", CallbackData: "goroutines"},
								{Text: "Channels", CallbackData: "channels"},
							},
						},
					},
				})
				// Подтверждение нажатия кнопки "Назад в меню"
				if err := bot.AnswerCallbackQuery(&telego.AnswerCallbackQueryParams{
					CallbackQueryID: update.CallbackQuery.ID,
				}); err != nil {
					fmt.Printf("Ошибка подтверждения нажатия кнопки: %v\n", err)
				}
				continue
			}

			// Если нажата одна из кнопок с темами
			if filename, exists := files[data]; exists {
				// Проверка существования файла
				if _, err := os.Stat(filename); os.IsNotExist(err) {
					content := []byte("Ошибка: файл не найден.")
					_, _ = bot.SendMessage(&telego.SendMessageParams{
						ChatID: chatID,
						Text:   string(content),
					})
				} else {
					content, err := os.ReadFile(filename)
					if err != nil {
						content = []byte("Ошибка: не удалось прочитать файл.")
					}

					// Отправка содержимого файла
					_, _ = bot.SendMessage(&telego.SendMessageParams{
						ChatID: chatID,
						Text:   string(content),
						ReplyMarkup: &telego.InlineKeyboardMarkup{
							InlineKeyboard: [][]telego.InlineKeyboardButton{
								{backButton}, // Кнопка "Назад в меню"
							},
						},
					})
				}
				// Подтверждение нажатия кнопки
				if err := bot.AnswerCallbackQuery(&telego.AnswerCallbackQueryParams{
					CallbackQueryID: update.CallbackQuery.ID,
				}); err != nil {
					fmt.Printf("Ошибка подтверждения нажатия кнопки: %v\n", err)
				}
			}
		}
	}
}
