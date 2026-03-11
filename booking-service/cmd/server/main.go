package main

import (
	"fmt"

	_ "github.com/leminhthai/train-ticket/booking-service/docs"
	"github.com/leminhthai/train-ticket/booking-service/global"
	"github.com/leminhthai/train-ticket/booking-service/internal/initialize"
)

// @title           Train Ticket Booking Service API
// @version         1.0
// @description     API documentation for Booking Service
// @host            localhost:8084
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	r := initialize.Run()

	r.Run(fmt.Sprintf(":%d", global.Config.Server.Port))
}
