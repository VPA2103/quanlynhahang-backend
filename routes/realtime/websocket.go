package realtime

import (
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/realtime"
)

func WebSocketRoutes(r *gin.Engine) {
	realtime.HubInstance = realtime.NewHub()
	go realtime.HubInstance.Run()

	r.GET("/ws", realtime.WsHandler(realtime.HubInstance))
}
