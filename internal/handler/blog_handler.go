package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/go-mogu/hz-framework/global"
	"github.com/go-mogu/hz-framework/pkg/consts/Constants"
	"github.com/go-mogu/hz-framework/pkg/consts/RedisConf"
	"github.com/go-mogu/hz-framework/pkg/mq"
	"github.com/go-mogu/hz-framework/pkg/util/gconv"
)

type blogHandler struct{}

var BlogHandler = &blogHandler{}

func (h *blogHandler) BlogConsumerHandler(mq *mq.RabbitMQ, data map[string]interface{}) error {
	ctx := context.Background()
	chErrors := make(chan error)
	go func() {
		deliveries, err := mq.Consume()
		if err != nil {
			chErrors <- err
		}
		for d := range deliveries {
			// 消费mq
			data = gconv.Map(d.Body)
			//从Redis清空对应的数据
			_, err = global.Redis.Del(ctx, RedisConf.BLOG_LEVEL+Constants.SYMBOL_COLON+gconv.String(Constants.NUM_ONE)).Result()
			_, err = global.Redis.Del(ctx, RedisConf.BLOG_LEVEL+Constants.SYMBOL_COLON+gconv.String(Constants.NUM_TWO)).Result()
			_, err = global.Redis.Del(ctx, RedisConf.BLOG_LEVEL+Constants.SYMBOL_COLON+gconv.String(Constants.NUM_THREE)).Result()
			_, err = global.Redis.Del(ctx, RedisConf.BLOG_LEVEL+Constants.SYMBOL_COLON+gconv.String(Constants.NUM_FOUR)).Result()
			_, err = global.Redis.Del(ctx, RedisConf.HOT_BLOG).Result()
			_, err = global.Redis.Del(ctx, RedisConf.NEW_BLOG).Result()
			_, err = global.Redis.Del(ctx, RedisConf.DASHBOARD+Constants.SYMBOL_COLON+RedisConf.BLOG_CONTRIBUTE_COUNT).Result()
			_, err = global.Redis.Del(ctx, RedisConf.DASHBOARD+Constants.SYMBOL_COLON+RedisConf.BLOG_COUNT_BY_SORT).Result()
			_, err = global.Redis.Del(ctx, RedisConf.DASHBOARD+Constants.SYMBOL_COLON+RedisConf.BLOG_COUNT_BY_TAG).Result()
			//comment := data[SysConf.COMMAND]
			//uid := data[SysConf.BLOG_UID]
			if err != nil {
				chErrors <- err
				return
			}
			err = d.Ack(true)
			if err != nil {
				return
			}
		}
	}()
	select {
	case err := <-chErrors:
		close(chErrors)
		hlog.Errorf("Consumer failed: %s\n", err)
		return err
	}
	return nil
}
