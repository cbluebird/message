package ws

import "message/app/service/messageService"

func InitClientHub() {
	messageService.ClientHub = messageService.NewHub()
	go messageService.ClientHub.Run()
}
