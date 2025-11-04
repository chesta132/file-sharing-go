package config

import (
	"path/filepath"
)

const (
	PORT             = "3000"         // Server port
	DB_PATH          = "data/data.db" // Database path
	UPLOAD_PATH      = "uploads"      // Save uploaded file path
	SAVE_SPLIT       = 10             // Filtering upload files wheter greater or lower than this variable (MB)
	MAX_UPLOAD       = 50             // Max upload file (MB)
	TOKEN_LENGTH     = 10             // Token length for file token
	PAGINATION_LIMIT = 20							// Pagination limit for get many endpoints

	KB = 1 << 10 // 1024
	MB = 1 << 20 // 1048576
	GB = 1 << 30 // 1073741824
)

var (
	LARGE_PATH = filepath.Join(UPLOAD_PATH, "/large") // Save filtered file path
	SMALL_PATH = filepath.Join(UPLOAD_PATH, "/small") // Save filtered file path
)
