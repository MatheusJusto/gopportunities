package main

import (
	"github.com/MatheusJusto/gopportunities/config"
	"github.com/MatheusJusto/gopportunities/router"
)

var (
	logger *config.Logger
)

func main() {
	logger = config.GetLooger("main")
	//initialize connfigs
	err := config.Init()
	if err != nil {
		logger.Errf("config initialization erro: %v", err)
		return
	}

	router.Initialize()
}
