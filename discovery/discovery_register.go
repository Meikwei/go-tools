/*
 * @Author: zhangkaiwei 1126763237@qq.com
 * @Date: 2024-04-29 21:52:30
 * @LastEditors: zhangkaiwei 1126763237@qq.com
 * @LastEditTime: 2024-05-01 16:44:20
 * @FilePath: \tools\discovery\discovery_register.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
// Copyright © 2023 OpenIM. All rights reserved.
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

package discovery

import (
	"context"

	"google.golang.org/grpc"
)

// Conn 接口定义了与服务发现和连接管理相关的功能。
type Conn interface {
	  // GetConns 用于获取与指定服务名称相关联的所有客户端连接。
    // ctx: 上下文，用于控制请求的取消、超时等。
    // serviceName: 需要获取连接的服务名称。
    // opts: 用于 grpc.DialOption 的可选配置。
    // 返回值: 成功时返回 []*grpc.ClientConn 连接列表，失败时返回 error。
	GetConns(ctx context.Context, serviceName string, opts ...grpc.DialOption) ([]*grpc.ClientConn, error)
	    // GetConn 用于获取与指定服务名称相关联的一个客户端连接。
    // ctx: 上下文，用于控制请求的取消、超时等。
    // serviceName: 需要获取连接的服务名称。
    // opts: 用于 grpc.DialOption 的可选配置。
    // 返回值: 成功时返回 *grpc.ClientConn 单个连接，失败时返回 error。
	GetConn(ctx context.Context, serviceName string, opts ...grpc.DialOption) (*grpc.ClientConn, error)
	// GetSelfConnTarget 用于获取当前连接的目标地址。
	GetSelfConnTarget() string
	  // AddOption 用于为后续的连接操作添加额外的配置选项。
    // opts: 一个或多个 grpc.DialOption 配置项。
	AddOption(opts ...grpc.DialOption)
	  // CloseConn 用于关闭指定的客户端连接。
    // conn: 需要关闭的 *grpc.ClientConn。
	CloseConn(conn *grpc.ClientConn)
	// do not use this method for call rpc

	GetClientLocalConns() map[string][]*grpc.ClientConn //del

	GetUserIdHashGatewayHost(ctx context.Context, userId string) (string, error) //del
}
// SvcDiscoveryRegistry 接口扩展了 Conn 接口，添加了服务注册与发现的功能。
type SvcDiscoveryRegistry interface {
	Conn
	  // Register 用于将服务注册到服务发现系统。
    // serviceName: 要注册的服务名称。
    // host: 服务所在的主机地址。
    // port: 服务的端口号。
    // opts: 用于 grpc.DialOption 的可选配置。
    // 返回值: 成功时返回 nil，失败时返回 error。
	Register(serviceName, host string, port int, opts ...grpc.DialOption) error
	UnRegister() error                                   //del
	RegisterConf2Registry(key string, conf []byte) error //del
	GetConfFromRegistry(key string) ([]byte, error)      //del
	Close()
}
