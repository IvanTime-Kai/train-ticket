package main

import (
	"fmt"

	"github.com/leminhthai/train-ticket/user-service/global"
	"github.com/leminhthai/train-ticket/user-service/internal/initialize"
)

// @title           Train Ticket User Service API
// @version         1.0
// @description     API documentation for User Service
// @host            localhost:8082
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token
func main() {
	r := initialize.Run()

	r.Run(fmt.Sprintf(":%d", global.Config.Server.Port))
}
