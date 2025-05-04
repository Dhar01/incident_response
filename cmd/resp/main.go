package main

import (
	"fmt"

	"github.com/Dhar01/incident_resp/config"
	"github.com/Dhar01/incident_resp/internal/database"
	"github.com/Dhar01/incident_resp/router"
)

func main() {
	err := config.Config()
	if err != nil {
		fmt.Println(err)
		return
	}
	configure := config.GetConfig()
	if config.IsRDBMS() {
		if err := database.InitDB().Error; err != nil {
			fmt.Println(err)
			return
		}
	}
	if config.IsRedis() {
		if _, err := database.InitRedis(); err != nil {
			fmt.Println(err)
			return
		}
	}
	if config.IsMongo() {
		if _, err := database.InitMongo(); err != nil {
			fmt.Println(err)
		}
	}

	r, err := router.SetUpRouter(*configure)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := r.Run(configure.Server.ServerHost + ":" + configure.Server.ServerPort); err != nil {
		fmt.Println(err)
		return
	}
}
