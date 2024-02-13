package admin_logic

import (
	"github.com/ppoonk/AirGo/global"
	"strconv"
	"strings"
)

type PushMessageService struct {
}

func (pm *PushMessageService) UnifiedPushMessage(msg string) {
	//推送tg消息
	if global.Server.Notice.TGAdmin != "" {
		tgIDs := strings.Fields(global.Server.Notice.TGAdmin)
		for _, v := range tgIDs {
			chatID, _ := strconv.ParseInt(v, 10, 64)
			TgBotSvc.TGBotSendMessage(chatID, msg)
		}
	}
}
func (pm *PushMessageService) StartTask() { //todo

}
