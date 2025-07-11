// Copyright 2024 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"net"

	"github.com/cloudwego/biz-demo/gomall/app/payment/biz/dal"
	"github.com/cloudwego/biz-demo/gomall/app/payment/conf"
	"github.com/cloudwego/biz-demo/gomall/app/payment/middleware"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/payment/paymentservice"
	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
)

func main() {
	_ = godotenv.Load()
	dal.Init()
	opts := kitexInit()

	svr := paymentservice.NewServer(new(PaymentServiceImpl), opts...)

	err := svr.Run()
	if err != nil {
		klog.Error(err.Error())
	}
}

func kitexInit() (opts []server.Option) {
	// address
	addr, err := net.ResolveTCPAddr("tcp", conf.GetConf().Kitex.Address)
	if err != nil {
		panic(err)
	}

	rateLimitMiddleware := func(next endpoint.Endpoint) endpoint.Endpoint {
		rateLimiter := rate.NewLimiter(rate.Limit(conf.GetConf().RateLimiter.Rate), conf.GetConf().RateLimiter.BucketSize)
		return func(ctx context.Context, req, resp interface{}) (err error) {
			if conf.GetConf().RateLimiter.Enabled {
				err = rateLimiter.Wait(ctx)
				if err != nil {
					return err
				}
			}
			return next(ctx, req, resp)
		}
	}

	opts = append(opts, server.WithServiceAddr(addr), server.WithMiddleware(rateLimitMiddleware))

	opts = append(opts,
		server.WithMiddleware(middleware.ServerMiddleware),
	)

	serviceName := conf.GetConf().Kitex.Service

	opts = append(opts,
		server.WithMetaHandler(transmeta.ServerHTTP2Handler),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: serviceName}),
	)

	return
}
