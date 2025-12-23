package controller

import (
	"github.com/gin-gonic/gin"
	"search-system/service"
)

func Search(c *gin.Context) {
	//query := c.Param("query")
	c.JSON(200, service.Search("test query"))
}

func GetDataSourceLenth(c *gin.Context) {
	c.JSON(200, service.GetDataSourceLenth())
}
