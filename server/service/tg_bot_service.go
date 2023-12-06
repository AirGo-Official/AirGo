package service

import (
	"context"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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

var EmailRegexp = regexp.MustCompile(`\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`)

func TGBotListen(bot *tgbotapi.BotAPI, ctx *context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates { //注：这里tg无消息时会阻塞，可能直到下一次来消息时退出协程，未验证，不影响使用
		select {
		case <-(*ctx).Done():
			bot.StopReceivingUpdates()
			global.Logrus.Info("bot exit")
			return
		default:
			if update.Message != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "hhh")
				ok := MessageAuth(update.Message)
				if ok { //管理员
					MessageHandlerForAdmin(&update, &msg)

				} else { //普通用户
					MessageHandlerForUser(&update, &msg)
				}
				bot.Send(msg)
			}

		}

	}

}
func TGBotStartListen() {
	if global.Server.Notice.BotToken == "" {
		return
	}

	bot, err := NewTGBot(global.Server.Notice.BotToken)
	if err != nil {
		return
	}
	global.TGBot = bot

	ctx, cancel := context.WithCancel(context.Background())

	global.ContextGroup.MapLock.Lock()
	global.ContextGroup.CtxMap[global.TGBotCtx] = &ctx
	global.ContextGroup.CancelMap[global.TGBotCancel] = &cancel
	global.ContextGroup.MapLock.Unlock()

	go TGBotListen(global.TGBot, global.ContextGroup.CtxMap[global.TGBotCtx])
}

func TGBotCloseListen() {
	cancel, ok := global.ContextGroup.CancelMap[global.TGBotCancel]
	if ok {
		(*cancel)()
		delete(global.ContextGroup.CtxMap, global.TGBotCtx)
		delete(global.ContextGroup.CancelMap, global.TGBotCancel)
	}
}
func TGBotSendMessage(chatID int64, text string) {
	if global.Server.Notice.BotToken == "" {
		return
	}
	msg := tgbotapi.NewMessage(chatID, text)
	global.TGBot.Send(msg)
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
		CmdStart(update, msg)
		goto tomsg
	case "bind":
		CmdBind(update, msg)
	}
	switch update.Message.Text {
	case "打卡":
		CmdClockin(update, msg)
	case "订阅":
		CmdGetUser(update, msg)

	case "绑定":
		msg.Text = "绑定格式：/bind xxx@qq.com|your_password"
	case "解绑":
		CmdUnbind(update, msg)

	case "TG ID":
		msg.Text = fmt.Sprintf("您的tg id：%d", update.Message.Chat.ID)
	case "官网":
		msg.Text = "官网：" + global.Server.Subscribe.FrontendUrl
	case "刷新菜单":
		CmdStart(update, msg)
	}

tomsg:
	msg.ReplyToMessageID = update.Message.MessageID

}
func MessageHandlerForAdmin(update *tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	switch update.Message.Command() {
	case "start":
		CmdStart(update, msg)
		goto tomsg
	case "bind":
		CmdBind(update, msg)
	case "find":
		CmdFindUser(update, msg)
	}
	switch update.Message.Text {
	case "绑定":
		msg.Text = "绑定格式：/bind xxx@qq.com|your_password"
	case "解绑":
		CmdUnbind(update, msg)
	case "打卡":
		CmdClockin(update, msg)

	case "订阅":
		CmdGetUser(update, msg)

	case "TG ID":
		msg.Text = fmt.Sprintf("您的tg id：%d", update.Message.Chat.ID)

	case "查询用户":
		msg.Text = "查询用户格式：/find xxx@qq.com"
	case "用户分析":
		UserAnalysis(update, msg)
	case "收入概览":
		Income(update, msg)
	case "节点状态":
		NodeStatus(update, msg)

	case "官网":
		msg.Text = "官网：" + global.Server.Subscribe.FrontendUrl

	case "刷新菜单":
		CmdStartForAdmin(update, msg)
	}

tomsg:
	msg.ReplyToMessageID = update.Message.MessageID

}
func CmdStart(up *tgbotapi.Update, msg *tgbotapi.MessageConfig) {
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
func CmdStartForAdmin(up *tgbotapi.Update, msg *tgbotapi.MessageConfig) {
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
	//查询用户
	user, err := FindUserByTgID(int64(up.Message.From.ID))
	if err != nil || user == nil {
		msg.Text = "未绑定账户"
		return
	}
	//判断是否已打卡
	_, ok := global.LocalCache.Get(fmt.Sprintf("%d", user.ID) + "clockin")
	if ok { //已打卡
		msg.Text = "已打卡"
		return
	}
	t, day, err := ClockIn(user.ID)
	if err != nil {
		msg.Text = "打卡错误"
		return
	}
	now := time.Now()
	zeroTime := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	global.LocalCache.Set(fmt.Sprintf("%d", user.ID)+"clockin", nil, zeroTime.Sub(now))

	msg.Text = fmt.Sprintf("天数: +%d，流量：+%dMB", day, t)
	return

}

func CmdBind(up *tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	//查询用户
	user, _ := FindUserByTgID(int64(up.Message.From.ID))
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

		user, err := FindUserByUserName(userName)
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
			SaveUser(user)
			msg.Text = fmt.Sprintf("TG ID: %d\n绑定账户：%s", up.Message.From.ID, userName)
		}
	}

}

func CmdUnbind(up *tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	//查询用户
	user, err := FindUserByTgID(int64(up.Message.From.ID))
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
			SaveUser(user)
			msg.Text = "已经解绑账户：" + user.UserName
			return
		}

	}

}

func CmdGetUser(up *tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	//查询用户
	user, _ := FindUserByTgID(int64(up.Message.From.ID))
	if user != nil {
		if user.TgID != 0 {

			expiredTime := user.SubscribeInfo.ExpiredAt.Format("2006-01-02 15-04-05")
			t := float64(user.SubscribeInfo.T) / 1024 / 1024 / 1024 //B->G
			u := float64(user.SubscribeInfo.U) / 1024 / 1024 / 1024 //B->G
			d := float64(user.SubscribeInfo.D) / 1024 / 1024 / 1024 //B->G
			subUrl := global.Server.Subscribe.BackendUrl + global.Server.Subscribe.ApiPrefix + "/user/getSub/?link=" + user.SubscribeInfo.SubscribeUrl + "&type=singbox"

			msg.Text = fmt.Sprintf("当前账户：%s\n订购套餐：%s\n到期时间：%s\n总流量：%s GB\n剩余流量：%s GB\n余额：%s\n通用订阅(singbox): %s", user.UserName, user.SubscribeInfo.GoodsSubject, expiredTime, fmt.Sprintf("%.2f", t), fmt.Sprintf("%.2f", t-u-d), fmt.Sprintf("%.2f", user.Remain), subUrl)

			return
		}

	}
	msg.Text = "未绑定账户"
}

func CmdFindUser(up *tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	userName := up.Message.Text[strings.LastIndex(up.Message.Text, "/find")+6:]
	userName = strings.TrimSpace(userName)
	ok := EmailRegexp.MatchString(userName)
	if !ok {
		msg.Text = "邮箱格式错误！查询格式：/find xxx@qq.com"
		return
	} else {
		user, err := FindUserByUserName(userName)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			msg.Text = "账户不存在！"
			return
		}
		if err != nil {
			msg.Text = "账户查询错误！"
			return
		}
		if user != nil {
			expiredTime := user.SubscribeInfo.ExpiredAt.Format("2006-01-02 15-04-05")
			t := float64(user.SubscribeInfo.T) / 1024 / 1024 / 1024 //B->G
			u := float64(user.SubscribeInfo.U) / 1024 / 1024 / 1024 //B->G
			d := float64(user.SubscribeInfo.D) / 1024 / 1024 / 1024 //B->G

			subUrl := global.Server.Subscribe.BackendUrl + global.Server.Subscribe.ApiPrefix + "/user/getSub/?link=" + user.SubscribeInfo.SubscribeUrl + "&type=singbox"

			msg.Text = fmt.Sprintf("账户：%s\n订购套餐：%s\n到期时间：%s\n总流量：%s GB\n剩余流量：%s GB\n余额：%s\n通用订阅(singbox): %s", user.UserName, user.SubscribeInfo.GoodsSubject, expiredTime, fmt.Sprintf("%.2f", t), fmt.Sprintf("%.2f", t-u-d), fmt.Sprintf("%.2f", user.Remain), subUrl)

			return
		}

	}

}

func UserAnalysis(up *tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	//总用户，已订阅用户，今日新增用户
	var userArr []model.User
	err := global.DB.Select("created_at, id, sub_status").Find(&userArr).Error
	if err != nil {
		msg.Text = "查询错误"
		return
	}
	now := time.Now()
	todayZeroTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	newUsers := 0
	alreadyPurchasedUsers := 0
	for _, v := range userArr {
		if v.CreatedAt.After(todayZeroTime) {
			newUsers += 1
		}
		if v.SubscribeInfo.SubStatus {
			alreadyPurchasedUsers += 1
		}
	}
	msg.Text = fmt.Sprintf("总用户: %d\n已订阅用户: %d\n今日新增用户: %d", len(userArr), alreadyPurchasedUsers, newUsers)

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

	global.GoroutinePool.Submit(func() {
		todayOrder, _ = GetOrderStatistics(todayStart, todayEnd)
	})
	global.GoroutinePool.Submit(func() {
		thisMonthOrder, _ = GetOrderStatistics(thisMonthStart, thisMonthEnd)

	})
	global.GoroutinePool.Submit(func() {
		lastMonthOrder, _ = GetOrderStatistics(lastMonthStart, lastMonthEnd)

	})

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
	err := global.DB.Find(&NodeArr).Error
	if err != nil {
		return ""
	}
	var msgArr []string
	for _, v := range NodeArr {
		_, ok := global.LocalCache.Get(fmt.Sprintf("%d%s", v.ID, global.NodeStatus))
		if ok {
			msgArr = append(msgArr, fmt.Sprintf("节点: %s, 状态: %s ", v.Remarks, "✅"))
		} else {
			msgArr = append(msgArr, fmt.Sprintf("节点: %s, 状态: %s ", v.Remarks, "❌"))
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
	for _, v := range NodeArr {
		_, ok := global.LocalCache.Get(fmt.Sprintf("%d%s", v.ID, global.NodeStatus))
		if !ok {
			//获取离线节点的通知状态，防止频繁推送
			_, ok = global.LocalCache.Get(fmt.Sprintf("%d%s", v.ID, global.NodeStatusIsNotified))
			if ok {
				continue
			}
			//设置离线节点的通知状态，防止频繁推送
			global.LocalCache.SetNoExpire(fmt.Sprintf("%d%s", v.ID, global.NodeStatusIsNotified), global.NodeStatusIsNotified)
			msgArr = append(msgArr, fmt.Sprintf("节点: %s, 状态: %s ", v.Remarks, "❌"))
		}
	}
	text := strings.Join(msgArr, "\n")
	return text
}
func CmdInfo(up *tgbotapi.Update) {
	//msg := tgbotapi.NewMessage(up.Message.Chat.ID, "请选择")
	//btn1 := tgbotapi.NewInlineKeyboardButtonData("查看信息", "")
	//row11 := tgbotapi.NewInlineKeyboardRow(btn1)
	//
	//keyboard := tgbotapi.NewInlineKeyboardMarkup(row11)
	//
	//msg.ReplyMarkup = keyboard
	//msg.ParseMode = tgbotapi.ModeMarkdown
	//return &msg

}
