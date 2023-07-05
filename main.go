package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type User struct {
	ID        string `json:"ID"`
	Username  string `json:"Username"`
	Firstname string `json:"Firstname"`
	Lastname  string `json:"Lastname"`
	Team      string `json:"Team"`
	ChatId    int64  `json:"ChatId"`
}

func NewUser(username, firstname, lastname, team string, chat_id int64) *User {
	return &User{
		ID:        uuid.NewString(),
		Username:  username,
		Firstname: firstname,
		Lastname:  lastname,
		Team:      team,
		ChatId:    chat_id,
	}
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := sql.Open("mysql", os.Getenv("dbUser")+":"+os.Getenv("dbPassword")+"@tcp("+os.Getenv("dbIP")+":"+os.Getenv("dbPort")+")/"+os.Getenv("dbDatabase"))
	if err != nil {
		log.Fatal("Error opening databse connection")
	}
	defer db.Close()

	bot, err := tgbotapi.NewBotAPI(os.Getenv("telegramToken"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {

		if update.Message != nil {
			log.Println("Message received...")
			user := NewUser(update.Message.From.UserName, update.Message.From.FirstName, update.Message.From.LastName, "", update.Message.Chat.ID)
			user, err := selectUser(db, user.ChatId)
			if err != nil {
				if err.Error() == "sql: no rows in result set" {
					log.Println("ChatID not found on database, creating...")
					user = NewUser(update.Message.From.UserName, update.Message.From.FirstName, update.Message.From.LastName, "", update.Message.Chat.ID)
					err = insertUser(db, user)
					if err != nil {
						panic(err)
					}
					log.Println("Created succefull")
				} else {
					log.Println("Some database problem: " + err.Error())
					panic(err)
				}

			} else {
				log.Println("ChatID found")
				log.Println(user)
			}

		}

		if update.Message.IsCommand() {

			switch commandReceived := update.Message.Command(); commandReceived {
			case "start":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "E ai campeão!! \nBem vindo ao BOT Clube Amigo! \nComandos suportados: \n/cadastrarTime \n/meuTime")
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
			case "cadastrarTime":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Informe o seu time agora: ")
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
			case "meuTime":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "O seu time é: ")
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Poxa, comando não suportado")
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
			}
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Nesse momento estou aceitando texto livre, mas não sei o que fazer com isso.. \nVou repetir o que vc me mandou: \n"+update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}

	}
}

func insertUser(db *sql.DB, user *User) error {
	stmt, err := db.Prepare("insert into cba_users(id, username, firstname, lastname, team, chat_id) values (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(user.ID, user.Username, user.Firstname, user.Lastname, user.Team, user.ChatId)
	if err != nil {
		return err
	}
	return nil
}

func updateUser(db *sql.DB, user *User) error {
	stmt, err := db.Prepare("update cba_users set username = ?, firstname = ?, lastname = ?, team = ?, chat_id = ? where id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(user.Username, user.Firstname, user.Lastname, user.Team, user.ChatId, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func selectUser(db *sql.DB, chat_id int64) (*User, error) {
	stmt, err := db.Prepare("select * from cba_users where chat_id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var us User
	err = stmt.QueryRow(chat_id).Scan(&us.ID, &us.Username, &us.Firstname, &us.Lastname, &us.Team, &us.ChatId)
	if err != nil {
		return nil, err
	}
	return &us, nil
}
