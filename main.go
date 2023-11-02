package main

import (
	sql "anon_chat/database_settings"
	redis "anon_chat/redis_client"
	"fmt"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	updater tgbotapi.Update
	bot     *tgbotapi.BotAPI
	r       redis.Redis
}

/*ype D_message struct{
	message int
}*/

// предустановки
const (
	stop_find_button   = "❌ Остановить поиск"
	start_find_button  = "⚡ Поиск"
	profile            = "📱 Профиль"
	search_person      = "🔍 Поиск собеседника запущен"
	stop_search_perosn = "🔍 Поиск остановлен"
	Registration_fine  = "✅ Регистрация успешно пройдена\nНажми на /start"
	time_to_find       = "👽 Пора искать тебе собеседника"
	person_find        = "Собеседник найден\nВся переписка защищена 🔒"
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

// обработка команд
func (h Handler) Commands(db *sql.DB) {

	msg := tgbotapi.NewMessage(h.updater.Message.Chat.ID, "")

	switch h.updater.Message.Command() {
	case "start", "menu":
		if db.Check_person(h.updater.Message.From.UserName) == false {
			msg.Text = "Привет, " + h.updater.Message.From.FirstName + "\nДобро пожаловать в самого продвинутого бота анонимных чатов в Телеграм, пора искать тебе собеседника\n\nНо для начала укажите свой пол ⬇ "
			msg.ReplyMarkup = sex_menu

		} else {
			msg.Text = time_to_find
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
func (h Handler) Messages(db *sql.DB) {
	msg := tgbotapi.NewMessage(h.updater.Message.Chat.ID, "")
	switch h.updater.Message.Text {
	case start_find_button:
		if chat := h.r.Queue_pop(); chat != "" {
			chat_id_find, err := strconv.Atoi(chat)
			if err != nil {
				log.Fatal(err)
			}
			db.Create_chat(h.updater.Message.Chat.ID, int64(chat_id_find))
			msg.Text = person_find
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			fmt.Println(chat_id_find)
			h.bot.Send(msg)
			msg.ChatID = int64(chat_id_find)
			h.bot.Send(msg)

		} else {
			h.r.Queue_add(h.updater.Message.Chat.ID)
			msg.ReplyMarkup = stop_menu
			msg.Text = search_person
			h.bot.Send(msg)
		}

	case stop_find_button:
		h.r.Queue_rem(h.updater.Message.Chat.ID)
		msg.ReplyMarkup = Start_menu
		msg.Text = stop_search_perosn
		h.bot.Send(msg)
	case profile:
		msg.Text = fmt.Sprintf("Имя: %s\n", h.updater.Message.From.FirstName)
		msg.ReplyMarkup = Start_menu
		h.bot.Send(msg)
	default:
		msg.Text = "Неизвестная команда"
		h.bot.Send(msg)
	}


}
func (h Handler) Edit_message(chat_id int64, lastMessageID int) {
	edit_message := tgbotapi.NewEditMessageText(chat_id, lastMessageID, Registration_fine)
	h.bot.Send(edit_message)
}

func (h Handler) Send_message(msg tgbotapi.Chattable) {
	if _, err := h.bot.Send(msg); err != nil {
		log.Panic(err)
	}
}

func (h Handler)Start_find(db *sql.DB)  {
	msg := tgbotapi.NewMessage(h.updater.Message.Chat.ID, "")
	if chat := h.r.Queue_pop(); chat != "" {
		chat_id_find, err := strconv.Atoi(chat)
		if err != nil {
			log.Fatal(err)
		}
		db.Create_chat(h.updater.Message.Chat.ID, int64(chat_id_find))
		msg.Text = person_find
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		fmt.Println(chat_id_find)
		h.bot.Send(msg)
		msg.ChatID = int64(chat_id_find)
		h.bot.Send(msg)

	} else {
		h.r.Queue_add(h.updater.Message.Chat.ID)
		msg.ReplyMarkup = stop_menu
		msg.Text = search_person
		h.bot.Send(msg)
	}
}

func (h Handler) Callback(db *sql.DB) {
	// здесь будут calbacks
}

func main() {
	db := sql.Open_db() //db connection
	defer func() { db.Close() }()
	defer func() { fmt.Println("CONNECTION CLOSE") }()

	r := redis.Create_client() //redis connection

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
			} else {
				up.Messages(db)
			}

		} else if update.CallbackQuery != nil {
			// Respond to the callback query, telling Telegram to show the user
			// a message with the data received.
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
			if _, err := bot.Request(callback); err != nil {
				panic(err)
			}

			switch update.CallbackQuery.Data {
			case "man":
				if db.Check_person(update.CallbackQuery.From.UserName) == false {
					db.Create_person(update.CallbackQuery.From.UserName, update.CallbackQuery.From.FirstName, "m")
					up.Edit_message(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
				}
			case "girl":
				if db.Check_person(update.CallbackQuery.From.UserName) == false {
					db.Create_person(update.CallbackQuery.From.UserName, update.CallbackQuery.From.FirstName, "w")
					up.Edit_message(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
				}
			}

		}
	}
}
