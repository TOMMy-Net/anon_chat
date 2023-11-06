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
	db      *sql.DB
}

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
	queue              = "❗ Вы уже находитесь в поиске!"
	stop_chat          = "🚶 Вы покинули чат"
	leave_chat         = "😓 Собеседник покинул чат"
	no_available_chat  = "🚫 У вас не было активного чата"
	no_access 	       = "⚠ Пройдите регистрацию"
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
func (h Handler) Commands() {

	msg := tgbotapi.NewMessage(h.updater.Message.Chat.ID, "")

	switch h.updater.Message.Command() {
	case "start", "menu":
		if h.db.Check_person(h.updater.Message.From.UserName) == false {
			msg.Text = "Привет, " + h.updater.Message.From.FirstName + "\nДобро пожаловать в самого продвинутого бота анонимных чатов в Телеграм, пора искать тебе собеседника\n\nНо для начала укажите свой пол ⬇ "
			msg.ReplyMarkup = sex_menu

		} else if chat := h.db.Get_active_chat(h.updater.Message.Chat.ID); chat != (sql.Chat{}) {
			h.db.Delete_chat(h.updater.Message.Chat.ID)
			msg.Text = stop_chat
			msg.ReplyMarkup = Start_menu
			if _, err := h.bot.Send(msg); err != nil {
				log.Panic(err)
			}
			msg.Text = leave_chat
			if chat_one := int64(chat.Chat_one); chat_one != h.updater.Message.Chat.ID {
				msg.ChatID = chat_one
			} else if chat_two := int64(chat.Chat_two); chat_two != h.updater.Message.Chat.ID {
				msg.ChatID = chat_two
			}
		} else if h.r.Queue_exist(h.updater.Message.Chat.ID) == true {
			h.r.Queue_rem(h.updater.Message.Chat.ID)
			msg.Text = stop_search_perosn
			msg.ReplyMarkup = Start_menu
		} else {
			msg.Text = time_to_find
			msg.ReplyMarkup = Start_menu
		}

	case "stop":
		if h.db.Check_person(h.updater.Message.From.UserName) == true {

			if chat := h.db.Get_active_chat(h.updater.Message.Chat.ID); chat != (sql.Chat{}) {
				h.db.Delete_chat(h.updater.Message.Chat.ID)
				msg.ReplyMarkup = Start_menu
				msg.Text = stop_chat
				h.Send(msg)
				msg.Text = leave_chat
				if chat_one := int64(chat.Chat_one); chat_one != h.updater.Message.Chat.ID {
					msg.ChatID = chat_one
				} else if chat_two := int64(chat.Chat_two); chat_two != h.updater.Message.Chat.ID {
					msg.ChatID = chat_two
				}

			} else if h.r.Queue_exist(h.updater.Message.Chat.ID) == true {
				h.r.Queue_rem(h.updater.Message.Chat.ID)
				msg.Text = stop_search_perosn
				msg.ReplyMarkup = Start_menu
			} else {
				msg.Text = no_available_chat
			}
		}else{
			msg.Text = no_access
		}
	default:
		msg.Text = "Неизвестная команда"
	}
	h.Send(msg)
}

// обработка сообщений
func (h Handler) Messages() {
	msg := tgbotapi.NewMessage(h.updater.Message.Chat.ID, "")
	switch h.updater.Message.Text {
	case start_find_button:
		if h.r.Queue_exist(h.updater.Message.Chat.ID) == false {
			if chat, err := strconv.Atoi(h.r.Queue_pop()); chat != 0 && chat != int(h.updater.Message.Chat.ID) {

				if err != nil {
					log.Fatal(err)
				}

				h.db.Create_chat(h.updater.Message.Chat.ID, int64(chat))

				msg.Text = person_find
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				fmt.Println(chat)
				h.Send(msg)
				msg.ChatID = int64(chat)
				h.Send(msg)

			} else {
				h.r.Queue_add(h.updater.Message.Chat.ID)
				msg.ReplyMarkup = stop_menu
				msg.Text = search_person
				h.Send(msg)
			}

		} else {
			msg.Text = queue
		}

	case stop_find_button:
		h.r.Queue_rem(h.updater.Message.Chat.ID)
		msg.ReplyMarkup = Start_menu
		msg.Text = stop_search_perosn
		h.Send(msg)
	case profile:
		msg.Text = fmt.Sprintf("Имя: %s\n", h.updater.Message.From.FirstName)
		msg.ReplyMarkup = Start_menu
		h.Send(msg)

	default:
		if chat := h.db.Get_active_chat(h.updater.Message.Chat.ID); chat != (sql.Chat{}) {
			if h.updater.Message.Voice != nil {
				voice := tgbotapi.NewVoice(h.updater.Message.Chat.ID, nil)
				voice.File = tgbotapi.FileID(h.updater.Message.Voice.FileID)

				if chat_one := int64(chat.Chat_one); chat_one != h.updater.Message.Chat.ID {
					voice.ChatID = chat_one
				} else if chat_two := int64(chat.Chat_two); chat_two != h.updater.Message.Chat.ID {
					voice.ChatID = chat_two
				}
				h.Send(voice)

			} else if h.updater.Message.Text != "" {
				msg.Text = h.updater.Message.Text
				if chat_one := int64(chat.Chat_one); chat_one != h.updater.Message.Chat.ID {
					msg.ChatID = chat_one
				} else if chat_two := int64(chat.Chat_two); chat_two != h.updater.Message.Chat.ID {
					msg.ChatID = chat_two
				}
				h.Send(msg)
			} else if h.updater.Message.Sticker != nil{
				sticker := tgbotapi.NewSticker(h.updater.Message.Chat.ID, tgbotapi.FileID(h.updater.Message.Sticker.FileID))
				if chat_one := int64(chat.Chat_one); chat_one != h.updater.Message.Chat.ID {
					sticker.ChatID = chat_one
				} else if chat_two := int64(chat.Chat_two); chat_two != h.updater.Message.Chat.ID {
					sticker.ChatID = chat_two
				}
				h.Send(sticker)
			}

		} else {
			msg.Text = "Неизвестная команда"
			h.Send(msg)
		}

	}

}

func (h Handler) Send(msg tgbotapi.Chattable) {
	if _, err := h.bot.Send(msg); err != nil {
		log.Panic(err)
	}
}

func (h Handler) Edit_message(chat_id int64, lastMessageID int) {
	edit_message := tgbotapi.NewEditMessageText(chat_id, lastMessageID, Registration_fine)
	h.Send(edit_message)
}

func main() {
	db := sql.Open_db() //db connection
	defer func() { db.Close() }()
	defer func() { fmt.Println("CONNECTION CLOSE") }()
	r := redis.Create_client() //redis connection
	defer func() { r.Client.Close() }()

	bot, err := tgbotapi.NewBotAPI("5837523403:AAExgbJOdXRFCJRow0Mw6j0Tqx_oR3J1F0Q")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 20

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		up := Handler{updater: update, bot: bot, r: r, db: db}

		if update.Message != nil && update.Message.Chat.Type == "private" { // If we got a message

			if update.Message == nil { // ignore any non-Message updates
				continue
			}
			if update.Message.IsCommand() {
				up.Commands()
			} else {
				up.Messages()
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
