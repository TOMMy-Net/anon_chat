package main

import (
	"log"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	sql "anon_chat/database_settings"
	redis "anon_chat/redis_client"
)


type Handler struct{
	updater	tgbotapi.Update
	bot *tgbotapi.BotAPI
	r redis.Redis
}
/*ype D_message struct{
	message int
}*/

//предустановки
const(
	stop_find_button = "❌ Остановить поиск"
	start_find_button = "⚡ Поиск"
	profile = "📱 Профиль"
	find_person = "🔔 Поиск собеседника запущен"
	stop_find_perosn = "🔔 Поиск остановлен"
	Registration_fine = "✅ Регистрация успешно пройдена\nНажми на /start"
	Time_to_find = "👽 Пора искать тебе собеседника"
)

var sex_menu = tgbotapi.NewInlineKeyboardMarkup(
    tgbotapi.NewInlineKeyboardRow(
        tgbotapi.NewInlineKeyboardButtonData("Парень", "man"),
        tgbotapi.NewInlineKeyboardButtonData("Девушка", "girl"),
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
	case "start", "menu":
		if db.Check_person(h.updater.Message.From.UserName) == false{
		msg.Text = "Привет, "+h.updater.Message.From.FirstName+"\nДобро пожаловать в самого продвинутого бота анонимных чатов в Телеграм, пора искать тебе собеседника\n\nНо для начала укажите свой пол ⬇ "
		msg.ReplyMarkup = sex_menu
		
		}else{
			msg.Text = Time_to_find
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
		msg.ReplyMarkup = Start_menu
	default:
		msg.Text = "Неизвестная команда"
	}

	if _, err := h.bot.Send(msg); err != nil {
		log.Panic(err)
	}
}
func (h Handler)Edit_message(chat_id int64, lastMessageID int){
	edit_message := tgbotapi.NewEditMessageText(chat_id, lastMessageID, Registration_fine)

  	if _, err := h.bot.Send(edit_message); err != nil {
		log.Panic(err)
	}
    
}

func (h Handler)Callback(db *sql.DB){
	// здесь будут calbacks
}

func main() {
	db := sql.Open_db()//db connection
	defer func ()  {db.Close()}()
	defer func ()  {fmt.Println("CONNECTION CLOSE")}()

	r := redis.Create_client()//redis connection

	bot, err := tgbotapi.NewBotAPI("5837523403:AAEfOk3fyrn0tZJnWAO7TJhLxq0RUGbPyR4")
	if err != nil {
		log.Panic(err)
	}
	
	bot.Debug = true


	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 20

	updates := bot.GetUpdatesChan(u)
	
	for update := range updates {
		up := Handler{updater: update, bot: bot, r: r}

		if update.Message != nil && update.Message.Chat.Type == "private" { // If we got a message
			

			if update.Message == nil { // ignore any non-Message updates
				continue
			}
			if update.Message.IsCommand() { 
				up.Commands(db)
			}else{
				up.Messages(db)
			}
			
		}else if update.CallbackQuery != nil{
            // Respond to the callback query, telling Telegram to show the user
            // a message with the data received.
            callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
            if _, err := bot.Request(callback); err != nil {
                panic(err)
            }
			
			switch update.CallbackQuery.Data{
			case "man":
				if db.Check_person(update.CallbackQuery.From.UserName) == false{
				db.Create_person(update.CallbackQuery.From.UserName, update.CallbackQuery.From.FirstName, "m")
				up.Edit_message(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
				
				}
			case "girl":
				if db.Check_person(update.CallbackQuery.From.UserName) == false{
				db.Create_person(update.CallbackQuery.From.UserName, update.CallbackQuery.From.FirstName, "w")
				up.Edit_message(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
				}
			}
			
		}
	}
}