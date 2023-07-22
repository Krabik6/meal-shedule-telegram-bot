package meals

import (
	"context"
	"fmt"
	"github.com/Krabik6/meal-shedule-telegram-bot/api"
	"github.com/Krabik6/meal-shedule-telegram-bot/internal/interfaces"
	"github.com/Krabik6/meal-shedule-telegram-bot/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

type ListService struct {
	JwtManager interfaces.JwtManager
	Bot        *tgbotapi.BotAPI
	Api        *api.Api
}

func NewListService(
	jwtManager interfaces.JwtManager,
	bot *tgbotapi.BotAPI,
	api *api.Api,

) *ListService {
	return &ListService{
		JwtManager: jwtManager,
		Bot:        bot,
		Api:        api,
	}
}

func (ls *ListService) MealPlansList(ctx context.Context, userID int64) error {
	token, err := ls.JwtManager.GetUserJWTToken(ctx, userID)
	if err != nil {
		return err
	}
	schedules, err := ls.Api.GetMealPlans(token)
	if err != nil {
		return err
	}

	mealPlans := groupMealPlans(schedules)

	var allErrors []error
	for key, mealPlan := range mealPlans {
		splitKey := strings.Split(key, ":")
		msg := tgbotapi.NewMessage(userID, formatMealPlanMessage(splitKey[0], splitKey[1], mealPlan))
		msg.ParseMode = "Markdown"
		if _, err := ls.Bot.Send(msg); err != nil {
			allErrors = append(allErrors, err)
		}
	}

	if len(allErrors) > 0 {
		errs := make([]string, len(allErrors))
		for i, err := range allErrors {
			errs[i] = err.Error()
		}
		return fmt.Errorf("errors occurred while sending messages: %s", strings.Join(errs, "; "))
	}
	return nil
}

func groupMealPlans(schedules []model.ScheduleOutput) map[string][]model.ScheduleOutput {
	// Group schedules by meal plan id and time
	mealPlans := make(map[string][]model.ScheduleOutput)
	for _, schedule := range schedules {
		key := fmt.Sprintf("%d:%s", schedule.Id, schedule.AtTime)
		mealPlans[key] = append(mealPlans[key], schedule)
	}
	return mealPlans
}

func formatMealPlanMessage(id string, atTime string, mealPlan []model.ScheduleOutput) string {
	msg := fmt.Sprintf("\n*ID*: %s\n*At Date*: %s", id, atTime)

	// Add all recipes in the meal plan to the message
	for _, schedule := range mealPlan {
		msg += fmt.Sprintf("\n------------------\n*Title*: %s\n*Description*: %s\n*Public*: %v\n*Cost*: %.2f\n*Time to prepare*: %d\n*Healthy*: %d",
			schedule.Title, schedule.Description, schedule.Public, schedule.Cost, schedule.TimeToPrepare, schedule.Healthy)
	}
	return msg
}
