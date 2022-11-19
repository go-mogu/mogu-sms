package router

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/app/server/registry"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/go-mogu/hz-framework/global"
	"github.com/go-mogu/hz-framework/pkg/response"
	"github.com/go-mogu/hz-framework/pkg/util"
	"github.com/go-mogu/mogu-registry/nacos"
	"github.com/hertz-contrib/requestid"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

func Register(port string) *server.Hertz {
	//获取本机ip
	addr := util.GetIpAddr()
	//nacos服务发现客户端
	nacosCli, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &global.Cfg.Nacos.Client,
			ServerConfigs: global.Cfg.Nacos.Server,
		})
	if err != nil {
		panic(err)
	}
	if global.Cfg.Server.Port == "" {
		global.Cfg.Server.Port = port
	}
	addr = addr + ":" + port
	//注册服务
	r := nacos.NewNacosRegistry(nacosCli)
	h := server.Default(
		server.WithHostPorts("0.0.0.0"+":"+port),
		server.WithRegistry(r, &registry.Info{
			ServiceName: global.Cfg.Server.Name,
			Addr:        utils.NewNetAddr("tcp", addr),
			Weight:      1,
			Tags:        global.Cfg.Nacos.Discovery.Metadata,
		}),
	)

	// header add X-Request-Id
	h.Use(requestid.New())
	// 404 not found
	h.NoRoute(func(c context.Context, ctx *app.RequestContext) {
		path := ctx.Request.URI().Path()
		method := ctx.Request.Method()
		response.NotFoundException(ctx, fmt.Sprintf("%s %s not found", method, path))
	})
	return h
}
