package route

import (
	"GlimmerMeeting/controllers"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	r.Static("/static", "./static")
	r.POST("/login", controllers.UserLogin)

	adminGrp := r.Group("/admin")
	{
		adminGrp.POST("/people", controllers.PostUser)
		adminGrp.GET("/people", controllers.GetUserList)
		adminGrp.PUT("/people", controllers.PutUser)
		adminGrp.DELETE("/people", controllers.DeleteUser)

		adminGrp.GET("/rooms", controllers.GetRoomList)
		adminGrp.POST("/rooms", controllers.PostRoom)
		adminGrp.PUT("/rooms", controllers.PutRoom)
		adminGrp.DELETE("/rooms", controllers.DeleteRoom)
	}

	return r
}
