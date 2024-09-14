package routes

import (
	"net/http"

	"Interview_Hin_20240914/docs"
	"Interview_Hin_20240914/middlewares"
	"Interview_Hin_20240914/models"
	"Interview_Hin_20240914/services"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func New() *gin.Engine {
	r := gin.New()
	initRoute(r)

	r.Use(gin.LoggerWithWriter(middlewares.LogWriter()))
	r.Use(gin.CustomRecovery(middlewares.AppRecovery()))
	r.Use(middlewares.CORSMiddleware())

	v1 := r.Group("/v1")
	{
		PingRoute(v1)
		PlayerRoute(v1)
		LevelRoute(v1)
		RoomRoute(v1)
		ReservationRoute(v1)
	}

	docs.SwaggerInfo.BasePath = v1.BasePath() // adds /v1 to swagger base path

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return r
}

func initRoute(r *gin.Engine) {
	_ = r.SetTrustedProxies(nil)
	r.RedirectTrailingSlash = false
	r.HandleMethodNotAllowed = true

	r.NoRoute(func(c *gin.Context) {
		models.SendErrorResponse(c, http.StatusNotFound, c.Request.RequestURI+" not found")
	})

	r.NoMethod(func(c *gin.Context) {
		models.SendErrorResponse(c, http.StatusMethodNotAllowed, c.Request.Method+" is not allowed here")
	})
}

func InitGin() {
	gin.DisableConsoleColor()
	gin.SetMode(services.Config.Mode)
	// do some other initialization staff
}
