package service

import (
	"AirGo/global"
	"AirGo/utils/net_plugin"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gorm.io/gorm"
	"log"
	"regexp"
	"runtime"
	"strings"
	"time"
)

var EmailRegexp = regexp.MustCompile(`\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`)

func TGBot(token string) {
	var bot *tgbotapi.BotAPI
	var err error

	if runtime.GOOS == "darwin" { //本机开发使用代理
		c := net_plugin.ClientWithSocks5("127.0.0.1", 1080, 10*time.Second)
		bot, err = tgbotapi.NewBotAPIWithClient(token, c)
	} else {
		bot, err = tgbotapi.NewBotAPI(token)
	}

	if err != nil {
		global.Logrus.Error("TGBot error:", err)
		return
	}
	//bot.Debug = true
	//log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "hhh")

			switch update.Message.Command() {
			case "start":
				CmdStart(&update, &msg)
				goto tomsg
			case "bind":
				CmdBindConfirm(&update, &msg)

			}
			switch update.Message.Text {
			case "打卡":
				CmdClockin(&update, &msg)
			case "订阅":
				CmdSub(&update, &msg)

			case "绑定":
				CmdBind(&update, &msg)

			case "解绑":
				CmdUnbind(&update, &msg)

			case "TG ID":
				msg.Text = fmt.Sprintf("您的tg id：%d", update.Message.Chat.ID)
			case "官网":
				msg.Text = "官网：" + global.Server.Subscribe.FrontendUrl
			case "刷新菜单":
				CmdStart(&update, &msg)
			}

		tomsg:
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
		if update.CallbackQuery != nil {
			fmt.Println("CallbackQuery:", update.CallbackQuery)
		}
	}
}

func CmdStart(up *tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	msg.Text = "菜单"
	//msg.Text = fmt.Sprintf("")
	fmt.Println("start")
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
	fmt.Println("msg:", msg)
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

func CmdBind(up *tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	//查询用户
	user, _ := FindUserByTgID(int64(up.Message.From.ID))
	if user != nil {
		if user.TgID != 0 {
			msg.Text = "已经绑定账户：" + user.UserName
			return

		}

	}
	msg.Text = "绑定格式：/bind xxx@qq.com"

}
func CmdBindConfirm(up *tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	userName := up.Message.Text[strings.LastIndex(up.Message.Text, "/bind")+6:]
	userName = strings.TrimSpace(userName)
	ok := EmailRegexp.MatchString(userName)
	if !ok {
		msg.Text = "邮箱格式错误！绑定格式：/bind xxx@qq.com"
		return
	} else {
		user, err := FindUserByUserName(userName)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			msg.Text = "账户不存在！绑定格式：/bind xxx@qq.com"
			return
		}
		if err != nil {
			msg.Text = "账户查询错误！绑定格式：/bind xxx@qq.com"
			return
		}
		if user != nil {
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

func CmdSub(up *tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	//查询用户
	user, _ := FindUserByTgID(int64(up.Message.From.ID))
	if user != nil {
		if user.TgID != 0 {

			expiredTime := user.SubscribeInfo.ExpiredAt.Format("2006-01-02 15-04-05")
			t := float64(user.SubscribeInfo.T) / 1024 / 1024 / 1024 //B->G
			u := float64(user.SubscribeInfo.U) / 1024 / 1024 / 1024 //B->G
			d := float64(user.SubscribeInfo.D) / 1024 / 1024 / 1024 //B->G

			msg.Text = fmt.Sprintf("当前账户：%s\n订购套餐：%s\n到期时间：%s\n总流量：%s GB\n剩余流量：%s GB\n余额：%s", user.UserName, user.SubscribeInfo.GoodsSubject, expiredTime, fmt.Sprintf("%.2f", t), fmt.Sprintf("%.2f", t-u-d), fmt.Sprintf("%.2f", user.Remain))

			return
		}

	}
	msg.Text = "未绑定账户"
}
