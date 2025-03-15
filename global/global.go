package global

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	Rdb *redis.Client
	Mdb *gorm.DB
)
