package database_settings

import (
	"database/sql"
	"fmt"
	"log"

	//"gorm.io/gorm"
	//"gorm.io/driver/sqlite"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	sql *sql.DB
}

type Chat struct {
	Id       int
	Chat_one int
	Chat_two int
}

func (db *DB) Create_table() {
	data, err := db.sql.Prepare("CREATE TABLE person ( id       INTEGER PRIMARY KEY, username TEXT    UNIQUE);")
	if err != nil {
		log.Fatal(err)
	}
	data.Exec()
}
func Open_db() *DB {
	db, err := sql.Open("sqlite3", "anon_db.db")
	if err != nil {
		panic(err)
	}

	data_b := DB{
		sql: db,
	}
	fmt.Println("Соединение с SQLite установлено")
	return &data_b
}

func (db *DB) Close() {
	db.sql.Close()
}

func (db *DB) Create_person(username string, name string, sex string) {
	db.Begin()
	data, err := db.sql.Prepare("INSERT INTO person (username, name, sex) VALUES (?, ?, ?)")
	if err != nil {
		panic(err)
	}
	data.Exec(username, name, sex)
	db.Commit()
}
func (db *DB) Check_person(username string) bool {
	rows, err := db.sql.Query("SELECT id, username FROM person WHERE username = ?", username)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	if rows.Next() {
		return true
	} else {
		return false
	}

}

func (db *DB) Create_chat(chat_one int64, chat_two int64) {
	db.Begin()
	data, err := db.sql.Prepare("INSERT INTO chats (chat_one, chat_two) VALUES (?, ?)")
	if err != nil {
		panic(err)
	}
	data.Exec(chat_one, chat_two)
	db.Commit()

}
func (db *DB) Check_chat(chat int64) bool {
	rows, err := db.sql.Query("SELECT * FROM chats WHERE chat_one = ? OR chat_two = ?", chat, chat)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	if rows.Next() {
		return true
	} else {
		return false
	}
}
func (db *DB) Get_active_chat(chat int64) Chat {
	rows, err := db.sql.Query("SELECT * FROM chats WHERE chat_one = ? OR chat_two = ?", chat, chat)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	c := Chat{}
	for rows.Next() {

		err = rows.Scan(&c.Id, &c.Chat_one, &c.Chat_two)
		if err != nil {
			fmt.Println(err)
		}

	}
	return c
}
func (db *DB) Delete_chat(chat int64) {
	db.Begin()
	data, err := db.sql.Prepare("DELETE FROM chats WHERE chat_one = ? OR chat_two = ?")
	if err != nil {
		panic(err)
	}
	data.Exec(chat, chat)
	db.Commit()
}
func (db *DB) Begin() error {
	stmt, err := db.sql.Prepare(`BEGIN`)
	if err != nil {
		return err
	}
	if _, err = stmt.Exec(); err != nil {
		return err
	}
	return nil
}

func (db *DB) Commit() error {
	stmt, err := db.sql.Prepare(`COMMIT`)
	if err != nil {
		return err
	}
	if _, err = stmt.Exec(); err != nil {
		return err
	}
	return nil
}
