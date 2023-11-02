package database_settings

import (
	"database/sql"
	"log"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct{
	sql *sql.DB
}
func (db *DB)Create_table(){
	data, err := db.sql.Prepare("CREATE TABLE person ( id       INTEGER PRIMARY KEY, username TEXT    UNIQUE);")
	if err!=nil {
		log.Fatal(err)
	}
	data.Exec()
}
func Open_db() *DB{
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

func (db *DB)Close(){
	db.sql.Close()
}

func (db *DB)Create_person(username string, name string, sex string){
	data, err := db.sql.Prepare("INSERT INTO person (username, name, sex) VALUES (?, ?, ?)")
	if err != nil{
		panic(err)
	}
	data.Exec(username, name, sex)
}
func (db *DB)Check_person(username string) bool{
	rows, err := db.sql.Query("SELECT id, username FROM person WHERE username = ?", username)
	if err != nil{
		log.Fatal(err)
	}
	if rows.Next() {
		return true
	} else {
		return false
	}
}

func (db *DB)Create_chat(chat_one, chat_two int64) {
	data, err := db.sql.Prepare("INSERT INTO chats (chat_one, chat_two) VALUES (?, ?)")
	if err != nil{
		panic(err)
	}
	data.Exec(chat_one, chat_two)
}