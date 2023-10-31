package main

import (
	"log"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	sql "anon_chat/database_settings"
)


type Handler struct{
	updater	tgbotapi.Update
	bot *tgbotapi.BotAPI
	
}
//предустановки
const(
	stop_find_button = "❌ Остановить поиск"
	start_find_button = "⚡ Найти собеседника"
	profile = "📱 Профиль"
	find_person = "🔔 Поиск собеседника запущен"
	stop_find_perosn = "🔔 Поиск остановлен"
)
var sex_menu = tgbotapi.NewInlineKeyboardMarkup(
    tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonData("Парень", "1"),
        tgbotapi.NewInlineKeyboardButtonData("Девушка", "0"),
    ),
)
var stop_menu = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(stop_find_button),
	
	),
)
var Start_menu = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(start_find_button),
		
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(profile),
		
	),
)

//обработка команд
func (h Handler)Commands(db *sql.DB) {
	
	msg := tgbotapi.NewMessage(h.updater.Message.Chat.ID, "")

	switch h.updater.Message.Command() {
	case "start":
		if db.Check_person(h.updater.Message.From.UserName) == false{
		db.Create_person(h.updater.Message.From.UserName, h.updater.Message.From.FirstName)
		msg.Text = "Привет, "+h.updater.Message.From.UserName+"\nДобро пожаловать в самого продвинутого бота анонимных чатов в Телеграм, пора искать тебе собеседника\nНо для начала выбери свой пол"
		msg.ReplyMarkup = sex_menu
		}else{
			msg.Text = "Привет"
			msg.ReplyMarkup = Start_menu
		}
		
	default:
		msg.Text = "Неизвестная команда"
	}
	if _, err := h.bot.Send(msg); err != nil {
		log.Panic(err)
	}
}

// обработка сообщений
func (h Handler)Messages(db *sql.DB){
	msg := tgbotapi.NewMessage(h.updater.Message.Chat.ID, "")
	switch h.updater.Message.Text {
	case start_find_button:
		msg.ReplyMarkup = stop_menu
		msg.Text = find_person
	case stop_find_button:
		msg.ReplyMarkup = Start_menu
		msg.Text = stop_find_perosn

	case profile:
		msg.Text = fmt.Sprintf("Имя: %s\n", h.updater.Message.From.FirstName)
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
	default:
		msg.Text = "Неизвестная команда"
	}

	if _, err := h.bot.Send(msg); err != nil {
		log.Panic(err)
	}
}
func main() {
	bot, err := tgbotapi.NewBotAPI("5837523403:AAEfOk3fyrn0tZJnWAO7TJhLxq0RUGbPyR4")
	if err != nil {
		log.Panic(err)
	}
	
	bot.Debug = true

	db := sql.Open_db()//db connection
	defer db.Close()

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 20

	updates := bot.GetUpdatesChan(u)
	
	for update := range updates {
		if update.Message != nil && update.Message.Chat.Type == "private" { // If we got a message
			up := Handler{updater: update, bot: bot}

			if update.Message == nil { // ignore any non-Message updates
				continue
			}
			if update.Message.IsCommand() { 
				up.Commands(db)
			}else{
				up.Messages(db)
			}
			
		}else if update.CallbackQuery != nil {
            // Respond to the callback query, telling Telegram to show the user
            // a message with the data received.
            callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
            if _, err := bot.Request(callback); err != nil {
                panic(err)
            }

            // And finally, send a message containing the data received.
            msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
            if _, err := bot.Send(msg); err != nil {
                panic(err)
            }
		}
	}
}