package main

import (
    "log"
    "strings"
    "GriBotGoalLev/internal"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Временное хранение целей в памяти
var userGoals = make(map[int][]string)

func main() {
    // Создаем соединение с Telegram ботом
    botToken := internal.BotToken // Получаем токен из config.go
    bot, err := tgbotapi.NewBotAPI(botToken) // Используем botToken вместо telegramBotToken
    if err != nil {
        log.Panic(err)
    }
    bot.Debug = true

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60
    updates, _ := bot.GetUpdatesChan(u)

    // Обработка сообщений
    for update := range updates {
        if update.Message == nil { // игнорируем не-текстовые обновления
            continue
        }

        userID := update.Message.From.ID
        msgText := strings.ToLower(update.Message.Text)
        response := ""

        switch {
        case strings.HasPrefix(msgText, "/setgoal "):
            goal := strings.TrimSpace(strings.TrimPrefix(msgText, "/setgoal "))
            addGoal(userID, goal)
            response = "Цель добавлена: " + goal

        case msgText == "/listgoals":
            goals := listGoals(userID)
            if len(goals) > 0 {
                response = "Ваши цели:\n" + strings.Join(goals, "\n")
            } else {
                response = "У вас нет активных целей."
            }

        case strings.HasPrefix(msgText, "/removegoal "):
            goal := strings.TrimSpace(strings.TrimPrefix(msgText, "/removegoal "))
            if removeGoal(userID, goal) {
                response = "Цель удалена: " + goal
            } else {
                response = "Цель не найдена."
            }

        default:
            response = "Доступные команды:\n" +
                "/setgoal <цель> — добавить цель\n" +
                "/listgoals — показать все цели\n" +
                "/removegoal <цель> — удалить цель"
        }

        msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
        bot.Send(msg)
    }
}

// Добавление новой цели для пользователя
func addGoal(userID int, goal string) {
    userGoals[userID] = append(userGoals[userID], goal)
}

// Получение списка целей для пользователя
func listGoals(userID int) []string {
    return userGoals[userID]
}

// Удаление цели пользователя
func removeGoal(userID int, goal string) bool {
    goals, exists := userGoals[userID]
    if !exists {
        return false
    }

    for i, g := range goals {
        if g == goal {
            userGoals[userID] = append(goals[:i], goals[i+1:]...) // Удаляем цель
            return true
        }
    }
    return false
}

