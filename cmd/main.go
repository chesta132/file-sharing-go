package main

import (
	"file-sharing/config"
	"file-sharing/internal/routers"
	"file-sharing/internal/services/db"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect client
	client := db.Connect(filepath.Join(config.DB_PATH, "data.db"), true)
	defer client.Close()

	router := gin.Default()
	r := routers.New(client)

	r.RegisterFile(router)

	router.Run(":" + config.PORT)
}
