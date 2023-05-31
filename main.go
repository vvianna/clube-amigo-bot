package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Usuario struct {
	NomeUsuario  string `json:"NomeUsuario"`
	PrimeiroNome string `json:"PrimeiroNome"`
	UltimoNome   string `json:"UltimoNome"`
	Time         string `json:"Time"`
}

func findUser(string) {
	for _, v := range Usuarios {
		if v.Key == "key1" {
			// Found!
		}
	}
}

func main() {

	Usuarios := make([]Usuario, 0)
	Usuarios = append(Usuarios, Usuario{
		NomeUsuario:  "Victor",
		PrimeiroNome: "Victor",
		UltimoNome:   "Monteiro",
		Time:         "Fluminense",
	})
	/*Usuarios[0] = Usuario{
		NomeUsuario:  "Rafa",
		PrimeiroNome: "Rafael",
		UltimoNome:   "Monteiro",
		Time:         "Fluminense",
	} */
	log.Println(Usuarios)
	log.Println(Usuarios[0].Time)

	bot, err := tgbotapi.NewBotAPI("Token")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {

		if update.Message.IsCommand() {

			//findUser()

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
