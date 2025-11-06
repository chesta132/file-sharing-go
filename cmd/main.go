package main

import (
	"file-sharing/config"
	"file-sharing/internal/routers"
	"file-sharing/internal/services/db"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect client
	client := db.Connect(config.DB_PATH, true)
	defer client.Close()

	router := gin.Default()
	r := routers.New(client)

	r.RegisterFile(router)

	router.Run(":" + config.PORT)
}
