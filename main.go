package main

import (
	"GlimmerMeeting/controllers"
	"GlimmerMeeting/repositories"
	"GlimmerMeeting/route"
	"fmt"
	log "github.com/sirupsen/logrus"
)

func main() {
	r := route.InitRouter()
	controllers.TokenMap = make(map[string]string)

	fmt.Println("\n   _____ _ _                               __  __           _   _             \n" +
		"  / ____| (_)                             |  \\/  |         | | (_)            \n " +
		"| |  __| |_ _ __ ___  _ __ ___   ___ _ __| \\  / | ___  ___| |_ _ _ __   __ _ \n " +
		"| | |_ | | | '_ ` _ \\| '_ ` _ \\ / _ \\ '__| |\\/| |/ _ \\/ _ \\ __| | '_ \\ / _` |\n" +
		" | |__| | | | | | | | | | | | | |  __/ |  | |  | |  __/  __/ |_| | | | | (_| |\n " +
		" \\_____|_|_|_| |_| |_|_| |_| |_|\\___|_|  |_|  |_|\\___|\\___|\\__|_|_| |_|\\__, |\n" +
		"                                                                         __/ |\n" +
		"                                                                        |___/ ")

	// 初始化数据库
	repositories.Init()

	err := r.Run(":2023")
	if err != nil {
		log.Errorf("backend exit accidentally: %v", err)
	}
}
