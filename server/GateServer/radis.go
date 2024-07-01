package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

// 配置模块中的Redis配置信息

var ctx = context.Background()

// 创建Redis客户端实例
var redisCli *redis.Client

// 初始化函数，设置连接和错误处理
func init() {
	// 监听连接错误事件
	redisCli = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", RedisHost, RedisPort),
		Password: RedisPasswd,
		DB:       0, // 使用默认的DB
	})

	pong, err := redisCli.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	fmt.Println("Redis connection successful:", pong)
	if redisCli != nil {
		go sendHeartbeat()
		redisCli.AddHook(redisHook{})
	}

}

// 自定义钩子函数，用于处理连接错误和重连
type redisHook struct{}

func (h redisHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (h redisHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	return nil
}

func (h redisHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (h redisHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	return nil
}

func sendHeartbeat() {
	for {
		// 发送心跳消息，比如向一个特定的 key 写入当前时间戳
		err := redisCli.Set(ctx, "heartbeat", time.Now().Unix(), 0).Err()
		if err != nil {
			log.Println("Failed to send heartbeat:", err)
		}
		time.Sleep(60 * time.Second) // 每 60 秒发送一次心跳消息
	}
}

// 根据key获取value
func getRedis(key string) (string, error) {
	result, err := redisCli.Get(ctx, key).Result()
	if err == redis.Nil {
		log.Printf("Result: <%v> This key cannot be found...\n", result)
		return "", nil
	} else if err != nil {
		log.Println("GetRedis error:", err)
		return "", err
	}
	log.Printf("Result: <%v> Get key success!\n", result)
	return result, nil
}

// 根据key查询redis中是否存在key
func queryRedis(key string) (bool, error) {
	result, err := redisCli.Exists(ctx, key).Result()
	if err != nil {
		log.Println("QueryRedis error:", err)
		return false, err
	}
	if result == 0 {
		log.Printf("Result: <%v> This key is null...\n", result)
		return false, nil
	}
	log.Printf("Result: <%v> With this value!\n", result)
	return true, nil
}

// 设置key和value，并过期时间
func setRedisExpire(key string, value interface{}, exptime time.Duration) error {
	err := redisCli.Set(ctx, key, value, exptime).Err()
	if err != nil {
		log.Println("SetRedisExpire error:", err)
		return err
	}
	// err = redisCli.Expire(ctx, key, exptime).Err()
	// if err != nil {
	// 	log.Println("SetRedisExpire error:", err)
	// 	return err
	// }
	return nil
}

// func main() {
// 	go sendHeartbeat()

// 	// 示例：设置Redis键值对
// 	err := setRedisExpire("exampleKey", "exampleValue", 60*time.Second)
// 	if err != nil {
// 		log.Println("Failed to set Redis key:", err)
// 	}

// 	// 示例：获取Redis键值对
// 	value, err := getRedis("exampleKey")
// 	if err != nil {
// 		log.Println("Failed to get Redis key:", err)
// 	} else {
// 		fmt.Println("Value:", value)
// 	}

// 	// 示例：查询Redis键是否存在
// 	exists, err := queryRedis("exampleKey")
// 	if err != nil {
// 		log.Println("Failed to query Redis key:", err)
// 	} else {
// 		fmt.Println("Exists:", exists)
// 	}

// 	// 阻止main函数退出
// 	select {}
// }
