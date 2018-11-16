package main

import (
	"flag"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"time"
	"os"
	"github.com/baitulakova/TelegramBotController/ffmpeg"
	"github.com/baitulakova/TelegramBotController/youtube"
)

var token=os.Getenv("TOKEN")
var key=os.Getenv("KEY")
var videoLink ="https://www.youtube.com/watch?v="

type Bot struct{
	BotAPI *tgbotapi.BotAPI
	Update tgbotapi.UpdateConfig
	UpdatesChannel tgbotapi.UpdatesChannel
}

func NewBot(token string)(Bot,error){
	Bot:=Bot{}
	if token==""{
		logrus.Fatal("Token is empty: ",token)
	}
	botApi,err:=tgbotapi.NewBotAPI(token)
	if err!=nil{
		return Bot,err
	}
	Bot.BotAPI=botApi
	Bot.BotAPI.Debug=false

	return Bot,err
}

func (bot *Bot)GetUpdates()error{
	u:=tgbotapi.NewUpdate(0)
	u.Timeout=60

	updates,err:=bot.BotAPI.GetUpdatesChan(u)
	if err!=nil{
		return err
	}
	bot.UpdatesChannel=updates
	updates.Clear()
	return nil
}

func (b *Bot) SendMessage(chatId int64,msg string){
	Message:=tgbotapi.NewMessage(chatId,msg)
	b.BotAPI.Send(Message)
}

func (b *Bot) SendAudio(chatId int64,filename string){
	audio:=tgbotapi.NewAudioUpload(chatId,filename)
	b.BotAPI.Send(audio)
}

func (b *Bot) SendNotification(update tgbotapi.Update,channel chan bool,msg string){
	for{
		b.SendMessage(update.Message.Chat.ID,msg)
		time.Sleep(time.Second*5)
		b.SendMessage(update.Message.Chat.ID,msg)
		time.Sleep(time.Second*5)
		b.SendMessage(update.Message.Chat.ID,msg)
		time.Sleep(time.Second*5)
		if <-channel{
			break
		}
	}
}

var stopped chan bool

func main(){
	flag.Parse()
	bot,err:=NewBot(token)
	if err!=nil{
		logrus.Fatal("Error creating bot")
	}

	logrus.Info("Bot started")
	go bot.start()

	<-stopped
}

func (bot *Bot) start() {
	err:=bot.GetUpdates()
	if err!=nil{
		logrus.Error("Error getting updates")
	}
	logrus.Info("Successfully got updates")
	time.Sleep(time.Millisecond*500)
	bot.UpdatesChannel.Clear()

	for update:=range bot.UpdatesChannel{
		if update.Message==nil{
			continue
		}
		go bot.HandleUpdate(update)
	}
	//send true to main
	stopped<-true
}

func (bot *Bot) HandleUpdate(update tgbotapi.Update){
	channel1:=make(chan bool)
	channel2:=make(chan bool)
	if update.Message.Text!="/start" {
		for {
			logrus.Info("Client sent: ", update.Message.Text)
			bot.SendMessage(update.Message.Chat.ID, "Started searching")
			client, err := youtube.NewYoutubeClient(key)
			if err != nil {
				logrus.Error("Error creating new youtube client: ", err)
				bot.SendMessage(update.Message.Chat.ID, "Unexpected error. Please try again")
				continue
			}

			logrus.Info("Started sending notifications to user")
			go bot.SendNotification(update, channel1, "I am searching video")

			logrus.Info("Started searching video")
			vid, err := client.Search(update.Message.Text)
			if err != nil {
				logrus.Error("Error searching: ", err)
				bot.SendMessage(update.Message.Chat.ID, "Unexpected error. Please try again")
				continue
			}
			logrus.Info("Successfully find video ", vid.Title)
			//send true to stop notifications
			channel1 <- true
			logrus.Info("Started getting video info")
			vidInfo, err := ffmpeg.GetVideoInfo(videoLink + vid.Id)
			if err != nil {
				logrus.Error("Error getting video info: ", err)
				bot.SendMessage(update.Message.Chat.ID, "Unexpected error. Please try again")
				continue
			}
			bot.SendMessage(update.Message.Chat.ID, "Started converting")

			logrus.Info("Started sending notifications to user")
			go bot.SendNotification(update, channel2, "Please wait. I am converting video")

			audioName := vid.Title + ".mp3"
			errConvert := vidInfo.GetDownloadLinkAndConvert(audioName)
			if errConvert != nil {
				logrus.Error("Converting error: ", err)
				bot.SendMessage(update.Message.Chat.ID, "I didn't convert this video. Please try enter other text")
				continue
			}
			logrus.Info("Finished  converting process")
			//send true to stop notifications
			channel2 <- true
			bot.SendAudio(update.Message.Chat.ID, audioName)
			logrus.Info("Sent message to client")
			os.Remove(audioName)
			logrus.Info("Removed file ", audioName)
			break
		}
	}
}
