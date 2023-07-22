package main

import (
	"context"
	"fmt"
	"github.com/Krabik6/meal-shedule-telegram-bot/api"
	"github.com/Krabik6/meal-shedule-telegram-bot/internal/bot_buttons"
	manager2 "github.com/Krabik6/meal-shedule-telegram-bot/internal/services/manager"
	"github.com/Krabik6/meal-shedule-telegram-bot/internal/services/states/expectation"
	"github.com/Krabik6/meal-shedule-telegram-bot/internal/services/states/login"
	"github.com/Krabik6/meal-shedule-telegram-bot/internal/services/states/meals"
	"github.com/Krabik6/meal-shedule-telegram-bot/internal/services/states/recipes"
	"github.com/Krabik6/meal-shedule-telegram-bot/internal/services/states/registration"
	"github.com/Krabik6/meal-shedule-telegram-bot/internal/transport/statehandlers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

const (
	expiration = 7 * 24 * time.Hour
)

func main() {
	// Создание клиента Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // Пароль Redis, если применимо
		DB:       0,  // Номер базы данных Redis, если применимо
	})

	//check manager connection
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Error connecting to Redis:", err)
	}
	//g
	botToken := "*******"
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	// Получение обновлений от Telegram API
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	jwtManager := manager2.NewJwtManager(redisClient, expiration)
	botMenu := bot_buttons.NewBotMenu(bot, jwtManager)
	botMenu.CreateMainMenu(ctx, 0)
	apis := api.NewApi()

	userStateManager := manager2.NewUserStateManager(redisClient)
	if botMenu == nil {
		log.Fatal("botMenu is nil")
	}

	recipesService := recipes.NewRecipesService(bot, redisClient, userStateManager, jwtManager, botMenu, apis)
	mealServices := meals.NewMealsService(bot, userStateManager, jwtManager, redisClient, botMenu, recipesService, apis)
	loginService := login.NewLoginService(redisClient, bot, jwtManager, userStateManager, botMenu, apis)
	registerService := registration.NewRegistrationService(bot, redisClient, userStateManager, botMenu, apis)
	noStateService := expectation.NewNoStateHandler(bot, userStateManager, jwtManager, botMenu, mealServices, recipesService)

	sh := statehandlers.NewStateHandler(redisClient, bot, userStateManager, recipesService, mealServices, loginService, registerService, noStateService)
	log.Println("Bot is running...")

	// Обработка входящих сообщений
	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}

		if update.Message != nil {

			// Handle message updates
			err := sh.HandleMessage(ctx, update.Message.Chat.ID, update.Message.Text, update)
			if err != nil {
				// В случае ошибки отправляем сообщение пользователю и логируем ошибку
				message := fmt.Sprintf("Произошла ошибка: %s\n Пожалуйста попробуйте ещё раз.", err)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
				_, err := bot.Send(msg)
				if err != nil {
					log.Println("Error sending message:", err)
				}
				log.Println("Error handling message:", err)
			}
		} else if update.CallbackQuery != nil {
			// Handle callback query updates
			err := sh.HandleCallbackQuery(ctx, update.CallbackQuery.Message.Chat.ID, update)
			if err != nil {
				// В случае ошибки отправляем сообщение пользователю и логируем ошибку
				message := fmt.Sprintf("Произошла ошибка: %s\n Пожалуйста попробуйте ещё раз.", err)
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, message)
				_, err := bot.Send(msg)
				if err != nil {
					log.Println("Error sending message:", err)
				}
				log.Println("Error handling callback query:", err)
			}
		}

	}
}
