package setu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/config"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
)

var instance *setu
var logger = utils.GetModuleLogger("com.aimerneige.setu")
var privateEnabled bool = false
var r18Enabled bool = false
var blacklistUser []int64
var allowedList []int64

type setu struct {
}

func init() {
	instance = &setu{}
	bot.RegisterModule(instance)
}

func (s *setu) MiraiGoModule() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       "com.aimerneige.setu",
		Instance: instance,
	}
}

// Init 初始化过程
// 在此处可以进行 Module 的初始化配置
// 如配置读取
func (s *setu) Init() {
	privateEnabled = config.GlobalConfig.GetBool("aimerneige.setu.private")
	r18Enabled = config.GlobalConfig.GetBool("aimerneige.setu.r18")
	blacklistUserSlice := config.GlobalConfig.GetIntSlice("aimerneige.setu.blacklist")
	for _, user := range blacklistUserSlice {
		blacklistUser = append(blacklistUser, int64(user))
	}
	logger.Info("blacklist user list:", blacklistUser)
	allowedListSlice := config.GlobalConfig.GetIntSlice("aimerneige.setu.allowed")
	for _, groupCode := range allowedListSlice {
		allowedList = append(allowedList, int64(groupCode))
	}
	logger.Info("allowed group list:", allowedList)
}

// PostInit 第二次初始化
// 再次过程中可以进行跨 Module 的动作
// 如通用数据库等等
func (s *setu) PostInit() {
}

// Serve 注册服务函数部分
func (s *setu) Serve(b *bot.Bot) {
	b.GroupMessageEvent.Subscribe(func(c *client.QQClient, msg *message.GroupMessage) {
		if !isAllowed(msg.GroupCode) {
			return
		}
		if msg.Sender.IsAnonymous() {
			return
		}
		if inBlacklist(msg.Sender.Uin) {
			return
		}
		if msg.ToString() == "来点色图" {
			c.SendGroupMessage(msg.GroupCode, message.NewSendingMessage().Append(message.NewText("不可以色色！")))
			sendSetu(c, msg.Sender.Uin, false, "")
		}
		if r18Enabled && msg.ToString() == "来点r18色图" {
			c.SendGroupMessage(msg.GroupCode, message.NewSendingMessage().Append(message.NewText("太色了！不可以！")))
			sendSetu(c, msg.Sender.Uin, true, "")
		}
		tag := parseTag(msg.ToString())
		if tag != "" {
			c.SendGroupMessage(msg.GroupCode, message.NewSendingMessage().Append(message.NewText("不可以色色！")))
			sendSetu(c, msg.Sender.Uin, false, tag)
		}
	})
	b.PrivateMessageEvent.Subscribe(func(c *client.QQClient, msg *message.PrivateMessage) {
		if !privateEnabled {
			return
		}
		if inBlacklist(msg.Sender.Uin) {
			return
		}
		if msg.ToString() == "来点色图" {
			sendSetu(c, msg.Sender.Uin, false, "")
		}
		if r18Enabled && msg.ToString() == "来点r18色图" {
			sendSetu(c, msg.Sender.Uin, true, "")
		}
		tag := parseTag(msg.ToString())
		if tag != "" {
			sendSetu(c, msg.Sender.Uin, false, tag)
		}
	})
}

// Start 此函数会新开携程进行调用
// ```go
//
//	go exampleModule.Start()
//
// ```
// 可以利用此部分进行后台操作
// 如 http 服务器等等
func (s *setu) Start(b *bot.Bot) {
}

// Stop 结束部分
// 一般调用此函数时，程序接收到 os.Interrupt 信号
// 即将退出
// 在此处应该释放相应的资源或者对状态进行保存
func (s *setu) Stop(b *bot.Bot, wg *sync.WaitGroup) {
	// 别忘了解锁
	defer wg.Done()
}

func parseTag(msg string) string {
	if msg == "" || msg == "来点色图" || msg == "来点r18色图" {
		return ""
	}
	if len(msg) <= 4 {
		return ""
	}
	msgRune := []rune(msg)
	length := len(msgRune)
	beforeTag := msgRune[:2]
	tag := msgRune[2 : length-2]
	afterTag := msgRune[length-2:]
	if string(beforeTag) == "来点" && string(afterTag) == "色图" {
		return string(tag)
	}

	return ""
}

func sendSetu(c *client.QQClient, id int64, r18 bool, tag string) {
	imgData, err := getSetuImg(r18, tag)
	if err != nil {
		logger.WithError(err).Error("Unable to get img from Lolicon API.")
	}
	imgMsgElement, err := c.UploadPrivateImage(id, imgData)
	if err != nil {
		logger.WithError(err).Error("Unable to Upload img.")
	}
	imgMsg := message.NewSendingMessage().Append(imgMsgElement)
	c.SendPrivateMessage(id, imgMsg)
}

func getSetuImg(r18 bool, tag string) (io.ReadSeeker, error) {
	apiURL := "https://api.lolicon.app/setu/v2"
	queryList := make([][]string, 0)
	if r18 {
		queryList = append(queryList, []string{"r18", "1"})
	} else {
		queryList = append(queryList, []string{"r18", "0"})
	}
	if tag != "" {
		queryList = append(queryList, []string{"tag", tag})
	}
	type loliconResponse struct {
		Error string `json:"error"`
		Data  []struct {
			Pid        int      `json:"pid"`
			P          int      `json:"p"`
			UID        int      `json:"uid"`
			Title      string   `json:"title"`
			Author     string   `json:"author"`
			R18        bool     `json:"r18"`
			Width      int      `json:"width"`
			Height     int      `json:"height"`
			Tags       []string `json:"tags"`
			Ext        string   `json:"ext"`
			UploadDate int64    `json:"uploadDate"`
			Urls       struct {
				Original string `json:"original"`
			} `json:"urls"`
		} `json:"data"`
	}
	var apiResponse loliconResponse
	apiResp, err := getRequest(apiURL, queryList)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(apiResp, &apiResponse); err != nil {
		return nil, err
	}
	if apiResponse.Error != "" {
		return nil, fmt.Errorf(apiResponse.Error)
	}
	imgURL := apiResponse.Data[0].Urls.Original
	imgBytes, err := getRequest(imgURL, [][]string{})
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(imgBytes), nil
}

func getRequest(url string, queryList [][]string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	for _, queryItem := range queryList {
		if len(queryItem) != 2 {
			return nil, fmt.Errorf("%v is not a valid query", queryItem)
		}
		q.Add(queryItem[0], queryItem[1])
	}
	req.URL.RawQuery = q.Encode()
	logger.Info(req.URL.String())
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func inBlacklist(userUin int64) bool {
	for _, v := range blacklistUser {
		if userUin == v {
			return true
		}
	}
	return false
}

func isAllowed(groupCode int64) bool {
	for _, v := range allowedList {
		if groupCode == v {
			return true
		}
	}
	return false
}
