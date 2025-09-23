package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Initialize() {

	// Crie um roteador Gin com os middlewares padrão (logger e recovery)
	router := gin.Default()

	// Defina um endpoint GET simples
	router.GET("/ping", func(c *gin.Context) {
		// Retorna resposta JSON
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	initializeRoutes(router)

	// Inicie o servidor na porta 8080 (padrão)
	// O servidor escutará em 0.0.0.0:8080 (localhost:8080 no Windows)
	router.Run()
}
