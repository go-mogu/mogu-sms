package handler

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/go-mogu/mgu-sms/global"
	"github.com/go-mogu/mgu-sms/internal/client"
	"github.com/go-mogu/mgu-sms/internal/consts"
	"github.com/go-mogu/mgu-sms/pkg/consts/Constants"
	"github.com/go-mogu/mgu-sms/pkg/consts/RedisConf"
	"github.com/go-mogu/mgu-sms/pkg/consts/SysConf"
	"github.com/go-mogu/mgu-sms/pkg/mq"
	"github.com/go-mogu/mgu-sms/pkg/util/gconv"
	"time"
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

			searchModel := consts.SQL
			result, err := client.WebClient.GetSearchModel(ctx)
			if err != nil {
				return
			}
			resultMap := gconv.MapStrStr(result)
			if SysConf.SUCCESS == resultMap[SysConf.CODE] {
				searchModel = resultMap[SysConf.DATA]
			}

			comment := data[SysConf.COMMAND]
			uid := data[SysConf.BLOG_UID]
			switch gconv.String(comment) {
			case SysConf.DELETE_BATCH:
				hlog.Info("mogu-sms处理批量删除博客")
				global.Redis.Set(ctx, RedisConf.BLOG_SORT_BY_MONTH+Constants.SYMBOL_COLON, "", -1)
				global.Redis.Set(ctx, RedisConf.MONTH_SET, "", -1)
				if searchModel == consts.ZINC {
					// 删除zinc博客索引
					_, err = client.SearchClient.DeleteElasticSearchByUidStr(ctx, gconv.String(uid))
					if err != nil {
						return
					}
				}
			case SysConf.EDIT_BATCH:
				global.Redis.Set(ctx, RedisConf.BLOG_SORT_BY_MONTH+Constants.SYMBOL_COLON, "", -1)
				global.Redis.Set(ctx, RedisConf.MONTH_SET, "", -1)
			case SysConf.ADD, SysConf.EDIT:
				updateSearch(ctx, data)
				// 增加ES索引
				_, err = client.SearchClient.AddElasticSearchByUid(ctx, gconv.String(uid))
				if err != nil {
					hlog.Error(err)
					return
				}
			case SysConf.DELETE:
				updateSearch(ctx, data)
				// 删除ES索引
				_, err = client.SearchClient.DeleteElasticSearchByUidStr(ctx, gconv.String(uid))
				if err != nil {
					hlog.Error(err)
					return
				}
			}

			if err != nil {
				chErrors <- err
				return
			}
			err = d.Ack(false)
			if err != nil {
				hlog.Error(err)
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
}

func updateSearch(ctx context.Context, data map[string]interface{}) {
	createTime, _ := time.ParseInLocation("2006-01-02", gconv.String(data[SysConf.CREATE_TIME]), time.Local)
	key := fmt.Sprintf("%d年%d月", createTime.Year(), createTime.Month())
	global.Redis.Del(ctx, RedisConf.BLOG_SORT_BY_MONTH+Constants.SYMBOL_COLON+key)
	jsonResult := global.Redis.Get(ctx, RedisConf.MONTH_SET)
	monthSet := make([]string, 0)
	err := jsonResult.Scan(&monthSet)
	if err != nil {
		return
	}
	haveMonth := false
	for _, item := range monthSet {
		if item == key {
			haveMonth = true
			break
		}
	}

	if !haveMonth {
		monthSet = append(monthSet, key)
		_, err = global.Redis.Set(ctx, RedisConf.MONTH_SET, monthSet, -1).Result()
		if err != nil {
			hlog.Error("更新Redis失败", err)
			return
		}
	}

}
