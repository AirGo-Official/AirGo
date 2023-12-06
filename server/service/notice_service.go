package service

import (
	"github.com/ppoonk/AirGo/global"
	"strconv"
	"strings"
)

func UnifiedPushMessage(msg string) {
	//推送tg消息
	if global.Server.Notice.TGAdmin != "" {
		tgIDs := strings.Fields(global.Server.Notice.TGAdmin)
		for _, v := range tgIDs {
			chatID, _ := strconv.ParseInt(v, 10, 64)
			TGBotSendMessage(chatID, msg)
		}
	}
}
