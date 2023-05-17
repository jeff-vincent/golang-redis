package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var (
	redis_host = "127.0.0.1"
	redis_port = "6379"
	redis_uri  = fmt.Sprintf("redis://%s:%s/0", redis_host, redis_port)
)

func write_cache(ctx *gin.Context) {
	key := ctx.Query("key")
	value := ctx.Query("value")

	opt, err := redis.ParseURL(redis_uri)
	if err != nil {
		panic(err)
	}
	rdb := redis.NewClient(opt)

	_, err = rdb.Set(ctx, key, value, 0).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"key":   key,
		"value": value,
	})
}

func read_cache(ctx *gin.Context) {
	key := ctx.Query("key")

	opt, err := redis.ParseURL(redis_uri)
	if err != nil {
		panic(err)
	}
	rdb := redis.NewClient(opt)

	value, err := rdb.Get(ctx, key).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"key":   key,
		"value": value,
	})
}

func main() {
	opt, err := redis.ParseURL(redis_uri)
	if err != nil {
		panic(err)
	}
	rdb := redis.NewClient(opt)

	router := gin.Default()
	router.GET("/", func(ctx *gin.Context) {
		ping_response := (rdb.Ping(rdb.Context()).Val())
		data := map[string]interface{}{
			"message": ping_response,
		}
		ctx.JSONP(http.StatusOK, data)
	})
	router.GET("/write", write_cache)
	router.GET("/read", read_cache)
	router.Run()
}
