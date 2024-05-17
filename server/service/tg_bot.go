package service

import (
	"context"
	"errors"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/utils/encrypt_plugin"
	"github.com/ppoonk/AirGo/utils/font_plugin"
	"github.com/ppoonk/AirGo/utils/net_plugin"
	"github.com/ppoonk/AirGo/utils/time_plugin"
	"github.com/vicanso/go-charts/v2"
	"gorm.io/gorm"

	"regexp"
	"strconv"
	"strings"
	"time"
)

func init() {
	err := charts.InstallFont("WenQuanWeiMiHei_modify", font_plugin.WenQuanWeiMiHei_modify)
	if err == nil {
		font, _ := charts.GetFont("WenQuanWeiMiHei_modify")
		charts.SetDefaultFont(font)
	}
}

type TgBotService struct {
	bot    *tgbotapi.BotAPI
	cancel context.CancelFunc
	enable bool
}

var EmailRegexp = regexp.MustCompile(`\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`)
var TgBotSvc TgBotService

func InitTgBotSvc() {
	TgBotSvc.StartTask()
	TgBotSvc.TGBotStart()
}

func (ts *TgBotService) TGBotStart() {
	if global.Server.Notice.BotToken == "" || !global.Server.Notice.EnableTGBot {
		ts.enable = false //机器人死了
		return
	}
	bot, err := NewTGBot(global.Server.Notice.BotToken)
	if err != nil {
		return
	}
	ts.bot = bot
	ctx, cancel := context.WithCancel(context.Background())
	ts.cancel = cancel
	ts.enable = true //机器人存活
	go ts.TGBotListen(ctx)
}
func (ts *TgBotService) TGBotListen(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := ts.bot.GetUpdatesChan(u)
	for {
		select {
		case <-ctx.Done():
			ts.bot.StopReceivingUpdates()
			global.Logrus.Info("bot exit")
			return
		case update := <-updates:
			if update.Message != nil {
				//fmt.Println("msg.From.ID:", update.Message.From.ID)
				//fmt.Println("update.Message.Text:", update.Message.Text)
				//fmt.Println("update.Message.Command():", update.Message.Command())
				//fmt.Println("update.Message.Chat.ID:", update.Message.Chat.ID)
				var msg tgbotapi.Chattable
				ok := IsAdmin(update.Message)
				if ok { //管理员
					msg = MessageHandlerForAdmin(&update)
				} else { //普通用户
					msg = MessageHandlerForUser(&update)
				}
				//消息入队
				if msg != nil {
					ts.TGBotSendMessage(msg)
				} else {
					ts.TGBotSendMessage(returnMsgForNil(&update))
				}
			}
		}
	}
}
func (ts *TgBotService) TGBotCloseListen() {
	if ts.cancel != nil {
		ts.cancel()
	}
}
func (ts *TgBotService) TGBotSendMessage(msg tgbotapi.Chattable) {
	global.Queue.Publish(constant.TG_BOT_SEND_MESSAGE, msg)
}

func (ts *TgBotService) StartTask() {
	ch, err := global.Queue.Subscribe(constant.TG_BOT_SEND_MESSAGE)
	if err != nil {
		global.Logrus.Error("Tg bot StartTask error:", err)
		return
	}
	go func() {
		for v := range ch {
			msg := v.(tgbotapi.Chattable)
			// 判断bot是否存活
			if ts.enable {
				_, err = ts.bot.Send(msg)
				if err != nil {
					global.Logrus.Error("Tg bot send msg error:", err)
				}
			}
		}
	}()
}

func NewTGBot(token string) (*tgbotapi.BotAPI, error) {
	var bot *tgbotapi.BotAPI
	var err error
	if global.Server.Notice.TGSocks5 != "" {
		socks := strings.Split(global.Server.Notice.TGSocks5, ":")
		add := socks[0]
		port := socks[1]
		portInt, _ := strconv.ParseInt(port, 10, 64)
		c := net_plugin.ClientWithSocks5(add, int(portInt), 10*time.Second)
		bot, err = tgbotapi.NewBotAPIWithClient(token, tgbotapi.APIEndpoint, c)
	} else {
		bot, err = tgbotapi.NewBotAPI(token)
	}
	if err != nil {
		global.Logrus.Error("TGBotListen error:", err)
		return nil, err
	}
	bot.Debug = false
	return bot, nil
}

func IsAdmin(msg *tgbotapi.Message) bool {
	_, ok := global.Server.Notice.AdminIDCacheWithTGID[msg.From.ID]
	if ok {
		return true
	}
	return false
}
func MessageHandlerForUser(update *tgbotapi.Update) tgbotapi.Chattable {
	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case "start":
			return ShowMenuForUser(update)
		case "bind":
			return BindTG(update)
		default:
			return nil
		}
	} else {
		switch update.Message.Text {
		case "我的订阅":
			return CreateSubChart(update)
		case "绑定":
			return BindTGPre(update)
		case "解绑":
			return UnbindTG(update)
		case "获取TG ID":
			return tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("您的tg id：%d", update.Message.Chat.ID))
		case "官网":
			return tgbotapi.NewMessage(update.Message.Chat.ID, "官网："+global.Server.Website.FrontendUrl)
		case "每日抽奖":
			return Lottery(update)

		default:
			return nil
		}
	}
}
func MessageHandlerForAdmin(update *tgbotapi.Update) tgbotapi.Chattable {
	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case "start":
			return ShowMenuForAdmin2(update)
		case "bind":
			return BindTG(update)
		case "find":
			return FindUser(update)
		default:
			return nil
		}
	} else {
		switch update.Message.Text {
		case "我的订阅":
			return CreateSubChart(update)
		case "绑定":
			return BindTGPre(update)
		case "解绑":
			return UnbindTG(update)
		case "获取TG ID":
			return tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("您的tg id：%d", update.Message.Chat.ID))
		case "官网":
			return tgbotapi.NewMessage(update.Message.Chat.ID, "官网："+global.Server.Website.FrontendUrl)
		case "查询用户":
			return tgbotapi.NewMessage(update.Message.Chat.ID, "查询用户格式：/find xxx@email.com")
		case "用户分析":
			return CreateUserSummaryChart(update)
		case "收入概览":
			return CreateOrderSummaryChart(update)
		case "节点状态":
			return CreateNodeStatusChart(update)
		case "切换用户菜单":
			return ShowMenuForAdmin1(update)
		case "切换管理员菜单":
			return ShowMenuForAdmin2(update)
		case "每日抽奖":
			return Lottery(update)
		default:
			return nil
		}
	}
}

func ShowMenuForUser(update *tgbotapi.Update) tgbotapi.Chattable {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	msg.ReplyToMessageID = update.Message.MessageID
	msg.Text = "菜单"
	bt1 := tgbotapi.NewKeyboardButton("我的订阅")
	bt2 := tgbotapi.NewKeyboardButton("绑定")
	bt3 := tgbotapi.NewKeyboardButton("解绑")
	bt4 := tgbotapi.NewKeyboardButton("获取TG ID")
	bt5 := tgbotapi.NewKeyboardButton("官网")
	bt6 := tgbotapi.NewKeyboardButton("每日抽奖")

	row1 := tgbotapi.NewKeyboardButtonRow(bt1)
	row2 := tgbotapi.NewKeyboardButtonRow(bt2, bt3, bt4)
	row3 := tgbotapi.NewKeyboardButtonRow(bt5, bt6)

	keyboard := tgbotapi.NewReplyKeyboard(row1, row2, row3)

	//keyboard.ResizeKeyboard = true
	msg.ReplyMarkup = keyboard
	//msg.ParseMode = tgbotapi.ModeMarkdown
	return msg
}
func ShowMenuForAdmin1(update *tgbotapi.Update) tgbotapi.Chattable {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	msg.ReplyToMessageID = update.Message.MessageID
	msg.Text = "菜单"
	bt1 := tgbotapi.NewKeyboardButton("我的订阅")
	bt2 := tgbotapi.NewKeyboardButton("绑定")
	bt3 := tgbotapi.NewKeyboardButton("解绑")
	bt4 := tgbotapi.NewKeyboardButton("获取TG ID")
	bt5 := tgbotapi.NewKeyboardButton("官网")
	bt6 := tgbotapi.NewKeyboardButton("每日抽奖")
	bt7 := tgbotapi.NewKeyboardButton("切换管理员菜单")

	row1 := tgbotapi.NewKeyboardButtonRow(bt1)
	row2 := tgbotapi.NewKeyboardButtonRow(bt2, bt3, bt4)
	row3 := tgbotapi.NewKeyboardButtonRow(bt5, bt6)
	row4 := tgbotapi.NewKeyboardButtonRow(bt7)

	keyboard := tgbotapi.NewReplyKeyboard(row1, row2, row3, row4)

	//keyboard.ResizeKeyboard = true
	msg.ReplyMarkup = keyboard
	msg.ParseMode = tgbotapi.ModeMarkdown
	return msg
}
func ShowMenuForAdmin2(update *tgbotapi.Update) tgbotapi.Chattable {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	msg.ReplyToMessageID = update.Message.MessageID
	msg.Text = "菜单"

	bt1 := tgbotapi.NewKeyboardButton("查询用户")
	bt2 := tgbotapi.NewKeyboardButton("用户分析")

	bt3 := tgbotapi.NewKeyboardButton("收入概览")
	bt4 := tgbotapi.NewKeyboardButton("节点状态")

	bt5 := tgbotapi.NewKeyboardButton("切换用户菜单")

	row1 := tgbotapi.NewKeyboardButtonRow(bt1, bt2)
	row2 := tgbotapi.NewKeyboardButtonRow(bt3, bt4)
	row3 := tgbotapi.NewKeyboardButtonRow(bt5)

	keyboard := tgbotapi.NewReplyKeyboard(row1, row2, row3)
	keyboard.ResizeKeyboard = true

	msg.ReplyMarkup = keyboard
	msg.ParseMode = tgbotapi.ModeMarkdown
	return msg
}
func BindTGPre(update *tgbotapi.Update) tgbotapi.Chattable {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	msg.ReplyToMessageID = update.Message.MessageID
	//查询用户
	user, _ := UserSvc.FirstUser(&model.User{TgID: int64(update.Message.From.ID)})
	if user != nil {
		if user.TgID != 0 {
			msg.Text = fmt.Sprintf("已经绑定账户: %s\n绑定其他用户请先解绑", user.UserName)
			return msg
		}
	}
	msg.Text = "绑定格式：/bind xxx@email.com|your_password"
	return msg
}
func BindTG(update *tgbotapi.Update) tgbotapi.Chattable {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	msg.ReplyToMessageID = update.Message.MessageID

	userText := strings.TrimSpace(update.Message.Text[strings.LastIndex(update.Message.Text, "/bind")+6:])
	userName := strings.Split(userText, "|")[0]
	pwd := strings.Split(userText, "|")[1]

	ok := EmailRegexp.MatchString(userName)
	if !ok {
		msg.Text = "邮箱格式错误！绑定格式：/bind xxx@email.com|your_password"
		return msg
	} else {
		user, err := UserSvc.FirstUser(&model.User{UserName: userName})
		if errors.Is(err, gorm.ErrRecordNotFound) {
			msg.Text = "账户不存在！绑定格式：/bind xxx@email.com|your_password"
			return msg
		}
		if err != nil {
			msg.Text = "账户查询错误！绑定格式：/bind xxx@email.com|your_password"
			return msg
		}
		if user != nil {
			if user.TgID != 0 {
				msg.Text = fmt.Sprintf("已经绑定账户: %s\n绑定其他用户请先解绑", user.UserName)
				return msg
			}
			if err = encrypt_plugin.BcryptDecode(pwd, user.Password); err != nil {
				msg.Text = "密码错误！绑定格式：/bind xxx@email.com|your_password"
				return msg

			}
			user.TgID = update.Message.From.ID
			UserSvc.SaveUser(user)
			msg.Text = fmt.Sprintf("TG ID: %d\n绑定账户：%s", update.Message.From.ID, userName)
			return msg
		}
		return nil
	}
}
func UnbindTG(update *tgbotapi.Update) tgbotapi.Chattable {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	msg.ReplyToMessageID = update.Message.MessageID
	user, err := UserSvc.FirstUser(&model.User{TgID: int64(update.Message.From.ID)})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		msg.Text = "未绑定账户"
		return msg
	}
	if err != nil {
		msg.Text = "账户查询错误！"
		return msg
	}
	if user != nil {
		if user.TgID != 0 {
			user.TgID = 0
			UserSvc.SaveUser(user)
			msg.Text = "已经解绑账户：" + user.UserName
			return msg
		}
	}
	return nil
}
func FindUser(update *tgbotapi.Update) tgbotapi.Chattable {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	msg.ReplyToMessageID = update.Message.MessageID

	text := strings.TrimSpace(update.Message.Text[strings.LastIndex(update.Message.Text, "/bind")+6:])
	user, err := UserSvc.FirstUser(&model.User{UserName: text})
	if err != nil {
		return returnMsgForErr(update, err)
	}
	msg.Text = strings.Join([]string{
		fmt.Sprintf("[ID]: %d", user.ID),
		fmt.Sprintf("[Name]: %s", user.UserName),
		fmt.Sprintf("[Balance]: %.2f", user.Balance),
		fmt.Sprintf("[Enable]: %v", user.Enable),
	}, "\n")
	return msg
}
func Lottery(update *tgbotapi.Update) tgbotapi.Chattable {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	msg.ReplyToMessageID = update.Message.MessageID

	user, err := UserSvc.FirstUser(&model.User{TgID: int64(update.Message.From.ID)})
	if err != nil {
		msg.Text = err.Error()
		return msg
	}
	_, balance, err := UserSvc.ClockIn(user.ID)
	if err != nil {
		msg.Text = err.Error()
		return msg
	}
	msg.Text = fmt.Sprintf("恭喜中奖！余额 + %.2f", balance)
	return msg
}

func CreateSubChart(update *tgbotapi.Update) tgbotapi.Chattable {
	header := []string{
		"Name",
		"Expiration",
		"Total",
		"Remain",
	}
	spans := map[int]int{
		0: 3,
		1: 4,
		2: 2,
		3: 2,
	}
	var data [][]string
	user, err := UserSvc.FirstUser(&model.User{TgID: int64(update.Message.From.ID)})
	if err != nil {
		return returnMsgForErr(update, err)
	}
	list, err := AdminCustomerServiceSvc.GetCustomerServiceList(&model.CustomerService{UserID: user.ID})
	if err != nil {
		return returnMsgForErr(update, err)
	}
	for _, v := range *list {
		//过期时间
		date := v.ServiceEndAt.Format("2006-01-02 15:04:05")
		//总流量
		remain1 := (float64(v.TotalBandwidth - v.UsedUp - v.UsedDown)) / 1024 / 1024 / 1024
		remain2 := strconv.FormatFloat(remain1, 'f', 2, 64)
		//剩余流量
		total1 := (float64(v.TotalBandwidth)) / 1024 / 1024 / 1024
		total2 := strconv.FormatFloat(total1, 'f', 2, 64)
		data = append(data, []string{v.Subject, date, total2 + " GB", remain2 + " GB"})
	}
	return returnMsgForPhoto(update, header, spans, data)
}

func CreateNodeStatusChart(update *tgbotapi.Update) tgbotapi.Chattable {
	header := []string{
		"ID",
		"Name",
		"Status",
		"UserAmount",
		"TrafficRate",
	}
	spans := map[int]int{
		0: 1,
		1: 4,
		2: 2,
		3: 3,
		4: 3,
	}
	var data, onlineNodes, offlineNodes [][]string
	nodeStatusList := AdminNodeSvc.GetNodesStatus()
	for _, v := range *nodeStatusList {
		if v.Status {
			onlineNodes = append(onlineNodes, []string{
				fmt.Sprintf("%d", v.ID), v.Name, "online", fmt.Sprintf("%d", v.UserAmount),
			})

		} else {
			offlineNodes = append(offlineNodes, []string{
				fmt.Sprintf("%d", v.ID), v.Name, "offline", fmt.Sprintf("%d", v.UserAmount),
			})
		}
	}
	data = append(data, onlineNodes...)
	data = append(data, offlineNodes...)
	if len(data) == 0 {
		return nil
	}
	return returnMsgForPhoto(update, header, spans, data)
}

func CreateUserSummaryChart(update *tgbotapi.Update) tgbotapi.Chattable {
	header := []string{
		"Date",
		"Number of registrations",
	}
	spans := map[int]int{
		0: 1,
		1: 2,
	}
	var data [][]string

	params := model.QueryParams{
		FieldParamsList: []model.FieldParamsItem{
			{
				ConditionValue: "",
			},
			{
				ConditionValue: "",
			},
		},
	}
	start, end := time_plugin.GetFirstToTodayForMonth()
	params.FieldParamsList[0].ConditionValue = start.Format("2006-01-02 15:04:05")
	params.FieldParamsList[1].ConditionValue = end.Format("2006-01-02 15:04:05")

	list, err := AdminUserSvc.UserSummary(&params)
	if err != nil {
		return returnMsgForErr(update, err)
	}
	for _, v := range *list {
		data = append(data, []string{
			v.Date, fmt.Sprintf("%d", v.RegisterTotal),
		})
	}
	return returnMsgForPhoto(update, header, spans, data)
}

func CreateOrderSummaryChart(update *tgbotapi.Update) tgbotapi.Chattable {
	header := []string{
		"Date",
		"General",
		"Recharge",
		"Subscribe",
		"Order",
		"Income",
	}
	spans := map[int]int{
		0: 4,
		1: 3,
		2: 3,
		3: 3,
		4: 3,
		5: 3,
	}
	var data [][]string

	params := model.QueryParams{
		FieldParamsList: []model.FieldParamsItem{
			{
				ConditionValue: "",
			},
			{
				ConditionValue: "",
			},
		},
	}
	start, end := time_plugin.GetFirstToTodayForMonth()
	params.FieldParamsList[0].ConditionValue = start.Format("2006-01-02 15:04:05")
	params.FieldParamsList[1].ConditionValue = end.Format("2006-01-02 15:04:05")

	list, err := AdminOrderSvc.OrderSummary(&params)
	if err != nil {
		return returnMsgForErr(update, err)
	}
	for _, v := range *list {
		data = append(data, []string{
			v.Date,
			fmt.Sprintf("%d", v.OrderTotal),
			fmt.Sprintf("%d", v.GeneralTotal),
			fmt.Sprintf("%d", v.RechargeTotal),
			fmt.Sprintf("%d", v.SubscribeTotal),
			fmt.Sprintf("%.2f", v.IncomeTotal),
		})
	}
	return returnMsgForPhoto(update, header, spans, data)
}

func returnMsgForErr(update *tgbotapi.Update, err error) tgbotapi.Chattable {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
	msg.ReplyToMessageID = update.Message.MessageID
	return msg
}
func returnMsgForNil(update *tgbotapi.Update) tgbotapi.Chattable {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "No data")
	msg.ReplyToMessageID = update.Message.MessageID
	return msg
}
func returnMsgForPhoto(update *tgbotapi.Update, header []string, spans map[int]int, data [][]string) tgbotapi.Chattable {
	p, err := charts.TableRender(
		header,
		data,
		spans,
	)
	if err != nil {
		return returnMsgForErr(update, err)
	}
	b, err := p.Bytes()
	if err != nil {
		return returnMsgForErr(update, err)
	}
	msg := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FileBytes{Name: "image.png", Bytes: b})
	msg.ReplyToMessageID = update.Message.MessageID
	return msg
}
