package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func initializeRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1")

	v1.GET("/opening", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "GET oppening",
		})
	})

	v1.POST("/opening", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "POST oppening",
		})
	})

	v1.DELETE("/opening", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "DELETE oppening",
		})
	})

	v1.PUT("/opening", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "PUT oppening",
		})
	})

	v1.GET("/openings", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "GET oppenings",
		})
	})
}
