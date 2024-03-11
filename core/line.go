package core

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-backend-test/infra"
	"github.com/go-backend-test/model"
	"github.com/go-backend-test/repo"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

type LineCore struct {
	LineChannelToken  string ` mapstructure:"LINE_CHANNEL_TOKEN"`
	LineChannelSecret string `mapstructure:"LINE_CHANNEL_SECRET"`
	UserRepo          repo.IUserRepo
	bot               *messaging_api.MessagingApiAPI
}

func NewLineCore() ILine {
	lr := LineCore{}
	initEnv := func() {
		if err := infra.GetEnv(&lr); err != nil {
			log.Fatal(err)
		}
	}

	initEnv()
	bot, err := messaging_api.NewMessagingApiAPI(lr.LineChannelToken)
	if err != nil {
		log.Fatal(err)
	}
	lr.bot = bot

	lr.UserRepo = repo.NewUserRepo()

	return &lr

}

func (r *LineCore) HandleCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		cb, err := webhook.ParseRequest(r.LineChannelSecret, c.Request)
		if err != nil {
			if errors.Is(err, webhook.ErrInvalidSignature) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid signature"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "unknown error"})
			}
			return
		}
		for _, event := range cb.Events {
			// log.Printf("Got event %v", event)
			switch e := event.(type) {
			case webhook.MessageEvent:
				switch message := e.Message.(type) {
				case webhook.TextMessageContent:
					var userId string
					if user, ok := e.Source.(webhook.UserSource); ok {
						userId = user.UserId
					}
					text := strings.TrimSpace(message.Text)

					switch text {
					case string(model.LISTMYMESSAGE):
						lmsgs, err := r.UserRepo.ListUserMessageByUserId(userId)
						if err != nil {
							fmt.Printf("Handle message error:%s", err.Error())
						}

						replyMsg := &messaging_api.ReplyMessageRequest{
							ReplyToken: e.ReplyToken,
							Messages:   []messaging_api.MessageInterface{messaging_api.TextMessage{Text: "Your message is empty"}},
						}

						if len(lmsgs) > 0 {
							msgs := make([]messaging_api.MessageInterface, 0)
							for _, v := range lmsgs {
								msgs = append(msgs, messaging_api.TextMessage{Text: fmt.Sprintf("Message:%s, CreatedAt:%s", v.Message, v.CreatedAt)})
							}

							replyMsg.Messages = msgs
						}

						if _, err := r.bot.ReplyMessage(replyMsg); err != nil {
							fmt.Println(err)
						}

					default: // add user message
						request := model.LineMessage{
							UserId:    userId,
							Message:   text,
							CreatedAt: time.Now(),
						}

						if err := r.UserRepo.AddUserMessage(request); err != nil {
							fmt.Printf("AddUserMessage error:%s", err.Error())
						}
					}
				default:
					log.Printf("Unknown event: %v", event)
				}
			}
		}
	}

}

func (r *LineCore) HandleSendMessageByAPI() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// get userId
		var req model.HandleSendMessageByAPIReq

		// 从 POST 请求的 body 中获取 JSON 数据并解析到结构体中
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"HandleSendMessageByAPI bind json error :": err.Error()})
			return
		}

		_, err := r.bot.PushMessage(
			&messaging_api.PushMessageRequest{
				To: req.LineUserId,
				Messages: []messaging_api.MessageInterface{
					messaging_api.TextMessage{
						Text: req.Message,
					},
				},
			},
			"", // x-line-retry-key
		)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"HandleSendMessageByAPI push message error :": err.Error()})
			return
		}
	}
}
