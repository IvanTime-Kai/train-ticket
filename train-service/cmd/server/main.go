package main

import (
	"fmt"

	_ "github.com/leminhthai/train-ticket/train-service/docs"
	"github.com/leminhthai/train-ticket/train-service/global"
	"github.com/leminhthai/train-ticket/train-service/internal/initialize"
)

// @title           Train Ticket Train Service API
// @version         1.0
// @description     API documentation for Train Service
// @host            localhost:8083
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	r := initialize.Run()

	r.Run(fmt.Sprintf(":%d", global.Config.Server.Port))
}
