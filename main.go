package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	sql "anon_chat/database_settings"
)


type Handler struct{
	updater	tgbotapi.Update
	bot *tgbotapi.BotAPI
	
}

var Start_menu = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Найти собеседника"),
		
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Профиль"),
		
	),
)


func (h Handler)Commands(db *sql.DB) {
	
	msg := tgbotapi.NewMessage(h.updater.Message.Chat.ID, "")

	switch h.updater.Message.Command() {
	case "start":
		if db.Check_person(h.updater.Message.From.UserName) == false{
		db.Create_person(h.updater.Message.From.UserName)
		msg.Text = "Привет, "+h.updater.Message.From.UserName+"\nДобро пожаловать в самого продвинутого бота анонимных чатов в Телеграм"
		}else{
			msg.Text = "Привет"
		}
		
		msg.ReplyMarkup = Start_menu
	default:
		msg.Text = "Неизвестная команда"
	}
	if _, err := h.bot.Send(msg); err != nil {
		log.Panic(err)
	}
}

func main() {
	bot, err := tgbotapi.NewBotAPI("5837523403:AAFsyNVi6hEwo2pywLSrnkA54u3dqUTgHwU")
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
		if update.Message != nil { // If we got a message
			up := Handler{updater: update, bot: bot}

			if update.Message == nil { // ignore any non-Message updates
				continue
			}
			if update.Message.IsCommand() { 
				up.Commands(db)
			}
			
			
	}
}
}