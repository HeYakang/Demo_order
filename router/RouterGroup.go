package router

import (
	"Demo_order/db"
	."Demo_order/handler"
	"github.com/gin-gonic/gin"
)


func CreatRouterGroup(){
	r := gin.Default()
    db.GetDb()
	userGroup := r.Group("/userInfo")
	{
		userGroup.GET("/userInfo", OnSearchHandler)
		userGroup.POST("/userInfo", OnCreatHandler)
		userGroup.PUT("/userInfo", OnUpDataHandler)
		userGroup.DELETE("/userInfo", OnDeleteHandler)
		userGroup.GET("/userInfoList", OnSearchListHandler)
		userGroup.GET("/excel", OnDownloadExcelHandler)

	}
	shopGroup := r.Group("/service")
	{
		shopGroup.PUT("/file", OnUpDateFileHandler)
		shopGroup.GET("/file", OnDownloadFileHandler)
	}
	r.Run(":9090")
}