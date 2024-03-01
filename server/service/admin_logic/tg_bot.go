package admin_logic

import (
	"context"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ppoonk/AirGo/constant"
	"github.com/ppoonk/AirGo/global"
	"github.com/ppoonk/AirGo/model"
	"github.com/ppoonk/AirGo/utils/encrypt_plugin"
	"github.com/ppoonk/AirGo/utils/net_plugin"
	"gorm.io/gorm"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type TgBotService struct {
	bot    *tgbotapi.BotAPI
	cancel *context.CancelFunc
}

var EmailRegexp = regexp.MustCompile(`\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`)
var TgBotSvc TgBotService

func InitTgBotSvc() {
	TgBotSvc.StartTask()
	TgBotSvc.TGBotStart()
}

func (ts *TgBotService) TGBotStart() {
	if global.Server.Notice.BotToken == "" {
		return
	}
	bot, err := NewTGBot(global.Server.Notice.BotToken)
	if err != nil {
		return
	}
	ts.bot = bot
	ctx, cancel := context.WithCancel(context.Background())
	ts.cancel = &cancel
	go ts.TGBotListen(ctx)
}
func (ts *TgBotService) TGBotListen(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, _ := ts.bot.GetUpdatesChan(u)
	for update := range updates { //注：这里tg无消息时会阻塞，可能直到下一次来消息时退出协程，未验证，不影响使用
		select {
		case <-ctx.Done():
			ts.bot.StopReceivingUpdates()
			global.Logrus.Info("bot exit")
			return
		default:
			if update.Message != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "hhh")
				//fmt.Println("msg.From.ID:", update.Message.From.ID)
				ok := MessageAuth(update.Message)
				if ok { //管理员
					MessageHandlerForAdmin(&update, &msg)

				} else { //普通用户
					MessageHandlerForUser(&update, &msg)
				}
				//ts.bot.Send(msg)
				//消息入队
				global.Queue.Publish(constant.TG_BOT_SEND_MESSAGE, msg)
			}
		}
	}
}
func (ts *TgBotService) TGBotCloseListen() {
	if ts.cancel != nil {
		(*ts.cancel)()
	}
}
func (ts *TgBotService) TGBotSendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
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
			msg := v.(tgbotapi.MessageConfig)
			_, err = ts.bot.Send(msg)
			if err != nil {
				global.Logrus.Error("Tg bot send msg error:", err)
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
		bot, err = tgbotapi.NewBotAPIWithClient(token, c)
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

// tg bot message authentication
func MessageAuth(msg *tgbotapi.Message) bool {
	res := strings.Index(global.Server.Notice.TGAdmin, fmt.Sprintf("%d", msg.From.ID))
	if res == -1 {
		return false
	}
	return true

}
func MessageHandlerForUser(update *tgbotapi.Update, msg *tgbotapi.MessageConfig) {

	switch update.Message.Command() {
	case "start":
		ShowMenuForUser(update, msg)
		goto tomsg
	case "bind":
		CmdBind(update, msg)
	}
	switch update.Message.Text {
	case "打卡":
		CmdClockin(update, msg)
	case "绑定":
		msg.Text = "绑定格式：/bind xxx@qq.com|your_password"
	case "解绑":
		CmdUnbind(update, msg)
	case "TG ID":
		msg.Text = fmt.Sprintf("您的tg id：%d", update.Message.Chat.ID)
	case "官网":
		msg.Text = "官网：" + global.Server.Website.FrontendUrl
	case "刷新菜单":
		ShowMenuForUser(update, msg)
	}

tomsg:
	msg.ReplyToMessageID = update.Message.MessageID

}
func MessageHandlerForAdmin(update *tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	switch update.Message.Command() {
	case "start":
		ShowMenuForUser(update, msg)
		goto tomsg
	case "bind":
		CmdBind(update, msg)
	}
	switch update.Message.Text {
	case "绑定":
		msg.Text = "绑定格式：/bind xxx@qq.com|your_password"
	case "解绑":
		CmdUnbind(update, msg)
	case "打卡":
		CmdClockin(update, msg)
	case "TG ID":
		msg.Text = fmt.Sprintf("您的tg id：%d", update.Message.Chat.ID)
	case "查询用户":
		msg.Text = "查询用户格式：/find xxx@qq.com"
	case "收入概览":
		Income(update, msg)
	case "节点状态":
		NodeStatus(update, msg)

	case "官网":
		msg.Text = "官网：" + global.Server.Website.FrontendUrl

	case "刷新菜单":
		ShowMenuForAdmin(update, msg)
	}

tomsg:
	msg.ReplyToMessageID = update.Message.MessageID

}
func ShowMenuForUser(up *tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	msg.Text = "菜单"
	bt1 := tgbotapi.NewKeyboardButton("打卡")

	bt2 := tgbotapi.NewKeyboardButton("订阅")

	bt3 := tgbotapi.NewKeyboardButton("绑定")
	bt4 := tgbotapi.NewKeyboardButton("解绑")
	bt5 := tgbotapi.NewKeyboardButton("TG ID")

	bt6 := tgbotapi.NewKeyboardButton("官网")
	bt7 := tgbotapi.NewKeyboardButton("刷新菜单")

	row1 := tgbotapi.NewKeyboardButtonRow(bt1, bt2)
	row2 := tgbotapi.NewKeyboardButtonRow(bt3, bt4, bt5)
	row3 := tgbotapi.NewKeyboardButtonRow(bt6, bt7)

	keyboard := tgbotapi.NewReplyKeyboard(row1, row2, row3)
	keyboard.ResizeKeyboard = true

	msg.ReplyMarkup = keyboard
	msg.ParseMode = tgbotapi.ModeMarkdown
}
func ShowMenuForAdmin(up *tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	msg.Text = "菜单"

	bt1 := tgbotapi.NewKeyboardButton("打卡")
	bt2 := tgbotapi.NewKeyboardButton("订阅")

	bt3 := tgbotapi.NewKeyboardButton("绑定")
	bt4 := tgbotapi.NewKeyboardButton("解绑")
	bt5 := tgbotapi.NewKeyboardButton("TG ID")

	bt6 := tgbotapi.NewKeyboardButton("查询用户")

	bt7 := tgbotapi.NewKeyboardButton("用户分析")
	bt8 := tgbotapi.NewKeyboardButton("收入概览")
	bt9 := tgbotapi.NewKeyboardButton("节点状态")

	bt10 := tgbotapi.NewKeyboardButton("官网")
	bt11 := tgbotapi.NewKeyboardButton("刷新菜单")

	row1 := tgbotapi.NewKeyboardButtonRow(bt1, bt2)
	row2 := tgbotapi.NewKeyboardButtonRow(bt3, bt4, bt5)
	row3 := tgbotapi.NewKeyboardButtonRow(bt6)
	row4 := tgbotapi.NewKeyboardButtonRow(bt7, bt8, bt9)
	row5 := tgbotapi.NewKeyboardButtonRow(bt10, bt11)

	keyboard := tgbotapi.NewReplyKeyboard(row1, row2, row3, row4, row5)
	keyboard.ResizeKeyboard = true

	msg.ReplyMarkup = keyboard
	msg.ParseMode = tgbotapi.ModeMarkdown
}
func CmdClockin(up *tgbotapi.Update, msg *tgbotapi.MessageConfig) {

}

func CmdBind(up *tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	//查询用户
	user, _ := userService.FirstUser(&model.User{ID: int64(up.Message.From.ID)})
	if user != nil {
		if user.TgID != 0 {
			msg.Text = "已经绑定账户：" + user.UserName
			return
		}
	}
	userText := strings.TrimSpace(up.Message.Text[strings.LastIndex(up.Message.Text, "/bind")+6:])

	userName := strings.Split(userText, "|")[0]
	pwd := strings.Split(userText, "|")[1]

	ok := EmailRegexp.MatchString(userName)
	if !ok {
		msg.Text = "邮箱格式错误！绑定格式：/bind xxx@qq.com|your_password"
		return
	} else {

		user, err := userService.FirstUser(&model.User{UserName: userName})
		if errors.Is(err, gorm.ErrRecordNotFound) {
			msg.Text = "账户不存在！绑定格式：/bind xxx@qq.com|your_password"
			return
		}
		if err != nil {
			msg.Text = "账户查询错误！绑定格式：/bind xxx@qq.com|your_password"
			return
		}
		if user != nil {
			if err = encrypt_plugin.BcryptDecode(pwd, user.Password); err != nil {
				msg.Text = "密码错误！绑定格式：/bind xxx@qq.com|your_password"
				return

			}
			user.TgID = int64(up.Message.From.ID)
			userService.SaveUser(user)
			msg.Text = fmt.Sprintf("TG ID: %d\n绑定账户：%s", up.Message.From.ID, userName)
		}
	}

}

func CmdUnbind(up *tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	user, err := userService.FirstUser(&model.User{ID: int64(up.Message.From.ID)})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		msg.Text = "未绑定账户"
		return
	}
	if err != nil {
		msg.Text = "账户查询错误！"
		return
	}
	if user != nil {
		if user.TgID != 0 {
			user.TgID = 0
			userService.SaveUser(user)
			msg.Text = "已经解绑账户：" + user.UserName
			return
		}
	}
}

func Income(up *tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	//今日订单 今日收入
	todayEnd := time.Now()
	todayStart := time.Date(todayEnd.Year(), todayEnd.Month(), todayEnd.Day(), 0, 0, 0, 0, todayEnd.Location())
	//本月订单 本月收入
	thisMonthStart := time.Date(todayEnd.Year(), todayEnd.Month(), 1, 0, 0, 0, 0, todayEnd.Location())
	thisMonthEnd := time.Date(todayEnd.Year(), todayEnd.Month()+1, 1, 0, 0, 0, 0, todayEnd.Location())
	//上月订单 上月收入
	lastMonthStart := time.Date(todayEnd.Year(), todayEnd.Month()-1, 1, 0, 0, 0, 0, todayEnd.Location())
	lastMonthEnd := time.Date(todayEnd.Year(), todayEnd.Month(), 1, 0, 0, 0, 0, todayEnd.Location())

	var (
		todayOrder     = &model.OrderStatistics{}
		thisMonthOrder = &model.OrderStatistics{}
		lastMonthOrder = &model.OrderStatistics{}
	)

	todayOrder, _ = orderService.GetOrderStatistics(todayStart, todayEnd)

	thisMonthOrder, _ = orderService.GetOrderStatistics(thisMonthStart, thisMonthEnd)

	lastMonthOrder, _ = orderService.GetOrderStatistics(lastMonthStart, lastMonthEnd)

	msg.Text = fmt.Sprintf("今日订单: %d, 今日收入: %.2f\n本月订单: %d, 本月收入: %.2f\n上月订单: %d, 上月订单: %.2f\n", todayOrder.Total, todayOrder.TotalAmount, thisMonthOrder.Total, thisMonthOrder.TotalAmount, lastMonthOrder.Total, lastMonthOrder.TotalAmount)

}
func NodeStatus(up *tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	text := GetNodeStatus()
	if text == "" {
		msg.Text = "查询错误"
		return
	}
	msg.Text = text

}
func GetNodeStatus() string {
	var NodeArr []model.Node
	err := global.DB.Where("enabled = ?", true).Find(&NodeArr).Error
	if err != nil {
		return ""
	}
	var msgArr []string
	for k, _ := range NodeArr {
		_, ok := global.LocalCache.Get(fmt.Sprintf("%s%d", constant.CACHE_NODE_STATUS_BY_NODEID, NodeArr[k].ID))
		if ok {
			msgArr = append(msgArr, fmt.Sprintf("节点: %s, 状态: %s ", NodeArr[k].Remarks, "✅"))
		} else {
			msgArr = append(msgArr, fmt.Sprintf("节点: %s, 状态: %s ", NodeArr[k].Remarks, "❌"))
		}
	}
	text := strings.Join(msgArr, "\n")
	return text
}
func GetOfflineNodeStatus() string {
	var NodeArr []model.Node
	err := global.DB.Find(&NodeArr).Error
	if err != nil {
		return ""
	}
	var msgArr []string
	for k, _ := range NodeArr {
		_, ok := global.LocalCache.Get(fmt.Sprintf("%s%d", constant.CACHE_NODE_STATUS_BY_NODEID, NodeArr[k].ID))
		if !ok {
			//获取离线节点的通知状态，防止频繁推送
			_, ok = global.LocalCache.Get(fmt.Sprintf("%s%d", constant.CACHE_NODE_STATUS_IS_NOTIFIED_BY_NODEID, NodeArr[k].ID))
			if ok {
				continue
			}
			//设置离线节点的通知状态，防止频繁推送
			global.LocalCache.SetNoExpire(fmt.Sprintf("%s%d", constant.CACHE_NODE_STATUS_IS_NOTIFIED_BY_NODEID, NodeArr[k].ID), constant.CACHE_NODE_STATUS_IS_NOTIFIED_BY_NODEID)
			msgArr = append(msgArr, fmt.Sprintf("节点: %s, 状态: %s ", NodeArr[k].Remarks, "❌"))
		}
	}
	text := strings.Join(msgArr, "\n")
	return text
}
