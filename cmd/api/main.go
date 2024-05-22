package main

import "notification-service/internal/app"

const configsDir = "configs"
const envDir = ".env"

func main() {

	app.Run(configsDir, envDir)

}
