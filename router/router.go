package router

import (
	"github.com/gin-gonic/gin"
	"search-system/controller"
)

func InitRouter(r *gin.Engine) {
	search := r.Group("/search")
	{
		search.GET("/datasource", controller.GetDataSourceLenth)
		search.GET("", controller.Search)
	}
}
