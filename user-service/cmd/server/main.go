package main

import (
	"fmt"

	"github.com/leminhthai/train-ticket/user-service/global"
	"github.com/leminhthai/train-ticket/user-service/internal/initialize"
)

func main() {
	r := initialize.Run()

	r.Run(fmt.Sprintf(":%d", global.Config.Server.Port))
}
