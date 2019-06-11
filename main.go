package main

import (
	"flag"
	"fmt"
	"github.com/leekchan/timeutil"
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	imageURL = "https://avatars1.githubusercontent.com/u/6083986"
)

func random(min, mac int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(mac-min) + min
}

func main() {
	cl, err := linebot.New(
		"secret",
		"access token",
	)
	if err != nil {
		log.Fatal(err)
	}
	var (
		httpAddr = flag.String("url", ":8080", "HTTP service address:localhost:8080")
	)
	flag.Parse()

	template := linebot.NewCarouselTemplate(
		linebot.NewCarouselColumn(
			imageURL, "bowwow.tips", "bowwow lin",
			linebot.NewURIAction("bowwow", "https://bowwow.tips"),
		),
	)
	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		res, err := cl.ParseRequest(req)
		if err != nil {
			log.Fatal(err)
		}
		for _, re := range res {
			if re.Type == linebot.EventTypeJoin {
				cl.ReplyMessage(re.ReplyToken, linebot.NewTextMessage("早上好！\n 請打「help」")).Do()
			}
			if re.Type == linebot.EventTypeFollow {
				n := time.Now()
				NowT := timeutil.Strftime(&n, "%Y年%m月%d日%H時%M分%S秒")
				p, _ := cl.GetProfile(re.Source.UserID).Do()
				cl.ReplyMessage(re.ReplyToken, linebot.NewTextMessage("謝謝你！\n"+p.DisplayName+"呃\n\n"+NowT)).Do()
				log.Println("DisplayName:" + p.DisplayName)
			}
			if re.Type == linebot.EventTypeMessage {
				switch msg := re.Message.(type) {
				case *linebot.TextMessage:
					if msg.Text == "test" {
						cl.ReplyMessage(re.ReplyToken, linebot.NewTextMessage("success")).Do()
					} else if msg.Text == "groupid" {
						cl.ReplyMessage(re.ReplyToken, linebot.NewTextMessage(string(re.Source.GroupID))).Do()
					} else if msg.Text == "byebye" {
						cl.ReplyMessage(re.ReplyToken, linebot.NewStickerMessage("3", "187")).Do()
						_, err := cl.LeaveGroup(re.Source.GroupID).Do()
						if err != nil {
							cl.LeaveRoom(re.Source.RoomID).Do()
						}
					} else if msg.Text == "help" {
						cl.ReplyMessage(re.ReplyToken, linebot.NewTextMessage("help\n・[image:画像url]=從圖片網址發送圖片\n・[speed]=測回話速度\n・[groupid]=發送GroupID\n・[roomid]=發送RoomID\n・[byebye]=取消訂閱\n・[about]=作者\n・[me]=發送發件人信息\n・[test]=test bowwow是否正常\n・[now]=現在時間\n・[mid]=mid\n・[sticker]=隨機圖片\n\n[其他機能]\n位置測試\n捉貼圖ID\n加入時發送消息")).Do()
					} else if msg.Text == "check" {
						fmt.Println(msg)
					} else if msg.Text == "now" {
						n := time.Now()
						NowT := timeutil.Strftime(&n, "%Y年%m月%d日%H時%M分%S秒")
						cl.ReplyMessage(re.ReplyToken, linebot.NewTextMessage(NowT)).Do()
					} else if msg.Text == "mid" {
						cl.ReplyMessage(re.ReplyToken, linebot.NewTextMessage(re.Source.UserID)).Do()
					} else if msg.Text == "roomid" {
						cl.ReplyMessage(re.ReplyToken, linebot.NewTextMessage(re.Source.RoomID)).Do()
					} else if msg.Text == "hidden" {
						cl.ReplyMessage(re.ReplyToken, linebot.NewTextMessage("hidden")).Do()
					} else if msg.Text == "bowwow" {
						_, err := cl.ReplyMessage(re.ReplyToken, linebot.NewImageMessage(imageURL, imageURL)).Do()
						if err != nil {
							log.Fatal(err)
						}
					} else if msg.Text == "sticker" {
						stid := random(180, 259)
						stidx := strconv.Itoa(stid)
						_, err := cl.ReplyMessage(re.ReplyToken, linebot.NewStickerMessage("3", stidx)).Do()
						if err != nil {
							log.Fatal(err)
						}
					} else if msg.Text == "me" {
						mid := re.Source.UserID
						p, err := cl.GetProfile(mid).Do()
						if err != nil {
							cl.ReplyMessage(re.ReplyToken, linebot.NewTextMessage("新增同意"))
						}

						cl.ReplyMessage(re.ReplyToken, linebot.NewTextMessage("mid:"+mid+"\nname:"+p.DisplayName+"\nstatusMessage:"+p.StatusMessage)).Do()
					} else if msg.Text == "speed" {
						replytoken := re.ReplyToken
						start := time.Now()
						cl.ReplyMessage(replytoken, linebot.NewTextMessage("..")).Do()
						end := time.Now()
						result := fmt.Sprintf("%f [sec]", (end.Sub(start)).Seconds())
						_, err := cl.PushMessage(re.Source.GroupID, linebot.NewTextMessage(result)).Do()
						if err != nil {
							_, err := cl.PushMessage(re.Source.RoomID, linebot.NewTextMessage(result)).Do()
							if err != nil {
								_, err := cl.PushMessage(re.Source.UserID, linebot.NewTextMessage(result)).Do()
								if err != nil {
									log.Fatal(err)
								}
							}
						}
					} else if res := strings.Contains(msg.Text, "hello"); res == true {
						cl.ReplyMessage(re.ReplyToken, linebot.NewTextMessage("hello!"), linebot.NewTextMessage("my name is bowwow")).Do()
					} else if res := strings.Contains(msg.Text, "image:"); res == true {
						image_url := strings.Replace(msg.Text, "image:", "", -1)
						cl.ReplyMessage(re.ReplyToken, linebot.NewImageMessage(image_url, image_url)).Do()
					} else if msg.Text == "about" {
						_, err := cl.ReplyMessage(re.ReplyToken, linebot.NewTemplateMessage("hi", template)).Do()
						if err != nil {
							log.Println(err)
						}
					}
				case *linebot.StickerMessage:
					cl.ReplyMessage(re.ReplyToken, linebot.NewTextMessage("StickerId:"+msg.StickerID+"\nPackageId:"+msg.PackageID)).Do()
				case *linebot.LocationMessage:
					cl.ReplyMessage(re.ReplyToken, linebot.NewLocationMessage(
						msg.Title,
						msg.Address,
						msg.Latitude,
						msg.Longitude)).Do()
				}
			}
		}
	})
	if err := http.ListenAndServe(*httpAddr, nil); err != nil {
		log.Fatal(err)
	}
}
