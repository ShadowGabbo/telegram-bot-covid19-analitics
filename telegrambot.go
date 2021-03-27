package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var Covid = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Italy"),
		tgbotapi.NewKeyboardButton("United States"),
		tgbotapi.NewKeyboardButton("Spain"),
		tgbotapi.NewKeyboardButton("France"),
		tgbotapi.NewKeyboardButton("Russia"),
		tgbotapi.NewKeyboardButton("United Kingdom"),
		tgbotapi.NewKeyboardButton("Brazil"),
		tgbotapi.NewKeyboardButton("India"),
	),
)

func main() {
	bot, err := tgbotapi.NewBotAPI("YOUR-TOKEN-HERE")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 120

	updates, err := bot.GetUpdatesChan(u)
	updates.Clear()
	fmt.Print(".")
	for update := range updates {
		if update.Message == nil {
			continue
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		if update.CallbackQuery != nil {
			fmt.Print(update)
			bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))
			bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data))
		}
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				msg.Text = "List of Commands\n/covid - to see covid analytics for some country\n/credits - to see all the people involved in the menagement of the bot"
				bot.Send(msg)
			case "covid":
				msg.ReplyMarkup = Covid
				bot.Send(msg)
			case "credits":
				msg.Text = "Creator: @Sartigabriele \nHelper: @SetteMagic"
				bot.Send(msg)
			default:
				msg.Text = "I don't know that command"
			}
		}
		if update.Message.Text != "" {
			switch update.Message.Text {
			case "Italy":
				msg.Text = CovidAnalitics("italy")
				bot.Send(msg)
			case "United States":
				msg.Text = CovidAnalitics("us")
				bot.Send(msg)
			case "India":
				msg.Text = CovidAnalitics("india")
				bot.Send(msg)
			case "Brazil":
				msg.Text = CovidAnalitics("brazil")
				bot.Send(msg)
			case "United Kingdom":
				msg.Text = CovidAnalitics("uk")
				bot.Send(msg)
			case "France":
				msg.Text = CovidAnalitics("france")
				bot.Send(msg)
			case "Spain":
				msg.Text = CovidAnalitics("spain")
				bot.Send(msg)
			case "Russia":
				msg.Text = CovidAnalitics("russia")
				bot.Send(msg)
			}
		}
	}
}

func CovidAnalitics(state string) string {
	// Request the HTML page.
	fmt.Println(state)
	url := "https://www.worldometers.info/coronavirus/country/" + state + "/"
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	text := doc.Find(".news_li").Text()
	switch state {
	case "italy":
		return Text(text)
	case "us":
		n := strings.Index(text, "States")
		textNew := text[:n+6]
		return textNew
	case "india":
		return Text(text)
	case "brazil":
		return Text(text)
	case "france":
		return Text(text)
	case "russia":
		return Text(text)
	case "uk":
		return Text(text)
	case "spain":
		n := strings.Index(text, ".")
		textNew := text[:n]
		return textNew
	}
	return ""
}

func Text(text string) (new string) {
	n := strings.Index(text, "[")
	textNew := text[:n]
	return textNew
}
