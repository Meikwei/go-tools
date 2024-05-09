/*
 * @Author: zhangkaiwei 1126763237@qq.com
 * @Date: 2024-04-29 21:52:30
 * @LastEditors: zhangkaiwei 1126763237@qq.com
 * @LastEditTime: 2024-05-01 12:17:07
 * @FilePath: \tools\discovery\zookeeper\options.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
// Copyright © 2024 OpenIM open source community. All rights reserved.
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

// 定义ZooKeeper客户端配置选项类型
package zookeeper

import (
	"time"

	"github.com/Meikwei/go-tools/log"
	"google.golang.org/grpc"
)

// ZkOption 为ZooKeeper客户端定制选项的函数类型
type ZkOption func(*ZkClient)

// WithRoundRobin 配置负载均衡策略为轮询
func WithRoundRobin() ZkOption {
    return func(client *ZkClient) {
        client.balancerName = "round_robin"
    }
}

// WithUserNameAndPassword 配置用户名和密码
func WithUserNameAndPassword(userName, password string) ZkOption {
    return func(client *ZkClient) {
        client.username = userName
        client.password = password
    }
}

// WithOptions 配置gRPC拨号选项
func WithOptions(opts ...grpc.DialOption) ZkOption {
    return func(client *ZkClient) {
        client.options = opts
    }
}

// WithFreq 配置频率
func WithFreq(freq time.Duration) ZkOption {
    return func(client *ZkClient) {
        client.ticker = time.NewTicker(freq)
    }
}

// WithTimeout 配置超时时间
func WithTimeout(timeout int) ZkOption {
    return func(client *ZkClient) {
        client.timeout = timeout
    }
}

// WithLogger 配置日志记录器
func WithLogger(logger log.Logger) ZkOption {
    return func(client *ZkClient) {
        client.logger = logger
    }
}
