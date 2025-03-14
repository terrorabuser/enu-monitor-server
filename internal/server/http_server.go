package server

import (
	"golang_gpt/internal/handler"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	socketio "github.com/googollee/go-socket.io"
	"time"
)

// SetupRouter инициализирует маршруты с mTLS
func RunHTTPServer(
	monitorHandler *handler.MonitorHandler,
	// contentHandler *handler.ContentHandler,
	socketServer *socketio.Server,
	authMonitorHandler *handler.AuthMonitorHandler,
	apiHandler *handler.ApiHandler,
) *gin.Engine{
	r := gin.Default()

	// Включаем CORS
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"}, // Разрешить доступ с любых источников
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
        AllowHeaders:     []string{"Origin", "Content-Type"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

	

	r.GET("/socket.io/*any", gin.WrapH(socketServer))
	r.POST("/socket.io/*any", gin.WrapH(socketServer))

	

	MonitorRoutes := r.Group("/monitor")
	{
		MonitorRoutes.POST("/refresh", authMonitorHandler.RefreshToken)
		MonitorRoutes.POST("/login", authMonitorHandler.MonitorLogin)
		MonitorRoutes.GET("/", monitorHandler.GetAllMonitors)
	}

	



	apiGroup := r.Group("/api")
	// apiGroup.Use(authHandler.AuthMiddleware("user", "admin"))
	{
    apiGroup.GET("/buildings", apiHandler.GetBuildings)         // Получить список зданий
    apiGroup.GET("/floors/:building", apiHandler.GetFloors)     // Получить этажи по зданию
    apiGroup.GET("/notes/:building/:floor", apiHandler.GetNotes) // Получить примечания по зданию и этажу
	}



	return r
}
