package main

import (
	"github.com/go-redis/redis/v7"
	"github.com/muhammednagy/PR-519-Software-development-project-backend/api"
	"github.com/muhammednagy/PR-519-Software-development-project-backend/api/models"
)

func init() {
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	defer rdb.Close()
	rdb.SAdd(models.ChannelsKey, "general", "random")
}

func main() {
	api.Run()
}
