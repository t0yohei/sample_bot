package main

import (
    "log"
    "net/http"
    "fmt"
    "math/rand"
    "strconv"
    "strings"
    "time"
    "os"
    "github.com/leekchan/timeutil"
    "github.com/line/line-bot-sdk-go/linebot"
)


func random(min,mac int) int {
    rand.Seed(time.Now().Unix())
    return rand.Intn(mac - min) + min
}

func main(){
    cl, err := linebot.New(
                os.Getenv("LINE_CHANNEL_SECRET"),
                os.Getenv("LINE_CHANNEL_TOKEN"),
        )
    if err != nil{
        log.Fatal(err)
    }
    http.HandleFunc("/callback",func(w http.ResponseWriter, req *http.Request) {
        res ,err := cl.ParseRequest(req)
        if err != nil{
            log.Fatal(err)
        }
        for _, re := range res{
            if re.Type == linebot.EventTypeJoin{
                cl.ReplyMessage(re.ReplyToken,linebot.NewTextMessage("こんにちは！\n「help」と発言していただければコマンド一覧を出します!")).Do()
            }
            if re.Type == linebot.EventTypeFollow{
                n := time.Now()
                NowT := timeutil.Strftime(&n,"%Y年%m月%d日%H時%M分%S秒")
                p,_ := cl.GetProfile(re.Source.UserID).Do()
                cl.ReplyMessage(re.ReplyToken,linebot.NewTextMessage("よろしくお願いします！\n"+p.DisplayName+"さん\n\n"+NowT)).Do()
                log.Println("DisplayName:"+p.DisplayName)
            }
            if re.Type == linebot.EventTypeMessage{
                switch msg := re.Message.(type){
                case *linebot.TextMessage:
                    if msg.Text == "test"{
                        cl.ReplyMessage(re.ReplyToken,linebot.NewTextMessage("success")).Do()
                    }else if msg.Text == "groupid"{
                        cl.ReplyMessage(re.ReplyToken,linebot.NewTextMessage(string(re.Source.GroupID))).Do()
                    }else if msg.Text == "byebye"{
                        cl.ReplyMessage(re.ReplyToken,linebot.NewStickerMessage("3","187")).Do()
                        _,err := cl.LeaveGroup(re.Source.GroupID).Do()
                        if err != nil{
                            cl.LeaveRoom(re.Source.RoomID).Do()
                        }
                    }else if msg.Text == "help"{
                        cl.ReplyMessage(re.ReplyToken,linebot.NewTextMessage("helpです\n・[image:画像のurl]=画像のurlから画像を送信します\n・[speed]=返答速度の測定をします\n・[groupid]=GroupIDを送信します\n・[roomid]=RoomIDを送信します\n・[byebye]=グループから退会します\n・[author]=作者のtwitterを送信します\n・[me]=送信者の情報を送信します\n・[test]=動いてるか確認します\n・[now]=現在の時刻を送信します\n・[mid]=送信者のmidを送信します\n・[Sticker]=ランダムでスタンプを送信します\n\nその他機能\n位置情報を送ると同じのを返します\nスタンプを送るとidを返します\n追加時にメッセージを送ります")).Do()
                    }else if msg.Text == "check"{
                        fmt.Println(msg)
                    }else if msg.Text == "now"{
                        n := time.Now()
                        NowT := timeutil.Strftime(&n,"%Y年%m月%d日%H時%M分%S秒")
                        cl.ReplyMessage(re.ReplyToken,linebot.NewTextMessage(NowT)).Do()
                    }else if msg.Text == "mid"{
                        cl.ReplyMessage(re.ReplyToken,linebot.NewTextMessage(re.Source.UserID)).Do()
                    }else if msg.Text == "roomid"{
                        cl.ReplyMessage(re.ReplyToken,linebot.NewTextMessage(re.Source.RoomID)).Do()
                    }else if msg.Text == "ひで"{
                        cl.ReplyMessage(re.ReplyToken,linebot.NewTextMessage("ひでしね")).Do()
                    }else if msg.Text == "えろいさん"{
                        _,err := cl.ReplyMessage(re.ReplyToken,linebot.NewImageMessage("https://pbs.twimg.com/media/DP3CXalVwAArqm8.jpg:large","https://pbs.twimg.com/media/DP3CXalVwAArqm8.jpg:large")).Do()
                        if err != nil{
                            log.Fatal(err)
                        }
                    }else if msg.Text == "Sticker"{
                        stid := random(180,259)
                        stidx := strconv.Itoa(stid)
                        _,err := cl.ReplyMessage(re.ReplyToken,linebot.NewStickerMessage("3",stidx)).Do()
                        if err != nil{
                            log.Fatal(err)
                        }
                    }else if msg.Text == "me"{
                        mid := re.Source.UserID
                        p,err := cl.GetProfile(mid).Do()
                        if err != nil{
                            cl.ReplyMessage(re.ReplyToken,linebot.NewTextMessage("追加して同意してください"))
                        }

                        cl.ReplyMessage(re.ReplyToken,linebot.NewTextMessage("mid:"+mid+"\nname:"+p.DisplayName+"\nstatusMessage:"+p.StatusMessage)).Do()
                    }else if msg.Text == "speed"{
                        replytoken := re.ReplyToken
                        start := time.Now()
                        cl.ReplyMessage(replytoken,linebot.NewTextMessage("..")).Do()
                        end := time.Now()
                        result := fmt.Sprintf("%f [sec]",(end.Sub(start)).Seconds())
                        _,err := cl.PushMessage(re.Source.GroupID,linebot.NewTextMessage(result)).Do()
                        if err != nil{
                            _,err := cl.PushMessage(re.Source.RoomID,linebot.NewTextMessage(result)).Do()
                            if err != nil{
                                _,err := cl.PushMessage(re.Source.UserID,linebot.NewTextMessage(result)).Do()
                                if err != nil{
                                    log.Fatal(err)
                                }
                            }
                        }
                    }else if res := strings.Contains(msg.Text,"Hello");res == true{
                        cl.ReplyMessage(re.ReplyToken,linebot.NewTextMessage("Hello!"),linebot.NewTextMessage("my name is bot")).Do()
                    }else if res := strings.Contains(msg.Text,"image:");res == true{
                        image_url := strings.Replace(msg.Text,"image:","",-1)
                        cl.ReplyMessage(re.ReplyToken,linebot.NewImageMessage(image_url,image_url)).Do()
                    }else if msg.Text == "author"{
                        _,err := cl.ReplyMessage(re.ReplyToken,linebot.NewTemplateMessage(
                            "twitter",
                            linebot.NewButtonsTemplate(
                                "https://cdn.downdetector.com/static/uploads/c/300/661ed/twitter-logo_10.png",
                                "作者twitter",
                                "フォローしてね！",
                                linebot.NewURITemplateAction("follow","https://twitter.com/intent/follow?screen_name=HomoRiron"),
                                ),
                            )).Do()
                        if err != nil{
                            log.Println(err)
                        }
                    }
                case *linebot.StickerMessage:
                    cl.ReplyMessage(re.ReplyToken,linebot.NewTextMessage("StickerId:"+msg.StickerID+"\nPackageId:"+msg.PackageID)).Do()
                case *linebot.LocationMessage:
                    cl.ReplyMessage(re.ReplyToken,linebot.NewLocationMessage(
                        msg.Title,
                        msg.Address,
                        msg.Latitude,
                        msg.Longitude)).Do()
                }
            }
        }
        })
    if err := http.ListenAndServe(":9000", nil); err != nil {
        log.Fatal(err)
    }
}
