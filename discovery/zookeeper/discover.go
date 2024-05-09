/*
 * @Author: zhangkaiwei 1126763237@qq.com
 * @Date: 2024-05-08 22:18:08
 * @LastEditors: zhangkaiwei 1126763237@qq.com
 * @LastEditTime: 2024-05-08 22:21:21
 * @FilePath: \go-tools\discovery\zookeeper\discover.go
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

package zookeeper

import (
	"context"
	"fmt"
	"strings"

	"github.com/Meikwei/go-tools/errs"
	"github.com/go-zookeeper/zk"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

// 定义ZkClient相关的错误类型。
var (
	ErrConnIsNil               = errs.New("conn is nil")
	ErrConnIsNilButLocalNotNil = errs.New("conn is nil, but local is not nil")
)

// watch方法用于监听zk事件。
func (s *ZkClient) watch(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			// 监听上下文结束，返回。
			s.logger.Info(ctx, "zk watch ctx done")
			return
		case event := <-s.eventChan:
			// 接收到zk事件，根据事件类型进行处理。
			s.logger.Debug(ctx, "zk eventChan recv new event", "event", event)
			switch event.Type {
			case zk.EventSession:
				// 处理会话事件。
				switch event.State {
				case zk.StateHasSession:
					// 会话建立，进行注册操作。
					if s.isRegistered && !s.isStateDisconnected {
						s.logger.Debug(ctx, "zk session event stateHasSession, client prepare to create new temp node", "event", event)
						node, err := s.CreateTempNode(s.rpcRegisterName, s.rpcRegisterAddr)
						if err != nil {
							s.logger.Error(ctx, "zk session event stateHasSession, create temp node error", err, "event", event)
						} else {
							s.node = node
						}
					}
				case zk.StateDisconnected:
					// 会话断开。
					s.isStateDisconnected = true
				case zk.StateConnected:
					// 会话重新连接。
					s.isStateDisconnected = false
				default:
					// 其他会话状态。
					s.logger.Debug(ctx, "zk session event", "event", event)
				}
			case zk.EventNodeChildrenChanged:
				// 子节点变化事件，进行服务解析和本地连接删除。
				s.logger.Debug(ctx, "zk event", "event", event)
				l := strings.Split(event.Path, "/")
				if len(l) > 1 {
					serviceName := l[len(l)-1]
					s.lock.Lock()
					s.flushResolverAndDeleteLocal(serviceName)
					s.lock.Unlock()
				}
				s.logger.Debug(ctx, "zk event handle success", "path", event.Path)
			case zk.EventNodeDataChanged:
			case zk.EventNodeCreated:
				// 节点创建事件，记录日志。
				s.logger.Debug(ctx, "zk node create event", "event", event)
			case zk.EventNodeDeleted:
			case zk.EventNotWatching:
			}
		}
	}
}

// GetConnsRemote方法用于获取指定服务的远程连接地址。
//
// 参数:
// ctx - 上下文，用于控制请求的取消、超时等。
// serviceName - 需要获取连接的服务名称。
//
// 返回值:
// conns - 远程服务地址列表。
// err - 获取过程中发生的错误。
func (s *ZkClient) GetConnsRemote(ctx context.Context, serviceName string) (conns []resolver.Address, err error) {
	err = s.ensureName(serviceName)
	if err != nil {
		return nil, err
	}

	path := s.getPath(serviceName)
	_, _, _, err = s.conn.ChildrenW(path)
	if err != nil {
		return nil, errs.WrapMsg(err, "children watch error", "path", path)
	}
	childNodes, _, err := s.conn.Children(path)
	if err != nil {
		return nil, errs.WrapMsg(err, "get children error", "path", path)
	} else {
		for _, child := range childNodes {
			fullPath := path + "/" + child
			data, _, err := s.conn.Get(fullPath)
			if err != nil {
				return nil, errs.WrapMsg(err, "get children error", "fullPath", fullPath)
			}
			s.logger.Debug(ctx, "get addr from remote", "conn", string(data))
			conns = append(conns, resolver.Address{Addr: string(data), ServerName: serviceName})
		}
	}
	return conns, nil
}

// GetUserIdHashGatewayHost方法用于获取用户ID的哈希值和网关主机信息。当前方法未实现。
//
// 参数:
// ctx - 上下文。
// userId - 用户ID。
//
// 返回值:
// string - 用户ID的哈希值和网关主机信息。
// error - 方法未实现错误。
func (s *ZkClient) GetUserIdHashGatewayHost(ctx context.Context, userId string) (string, error) {
	s.logger.Warn(ctx, "not implement", errs.New("zkclinet not implement GetUserIdHashGatewayHost method"))
	return "", nil
}

// GetConns方法用于获取指定服务的连接。首先尝试从本地缓存获取，若不存在则从远程获取并缓存。
//
// 参数:
// ctx - 上下文。
// serviceName - 需要获取连接的服务名称。
// opts - gRPC拨号选项。
//
// 返回值:
// []*grpc.ClientConn - 服务连接列表。
// error - 获取连接过程中发生的错误。
func (s *ZkClient) GetConns(ctx context.Context, serviceName string, opts ...grpc.DialOption) ([]*grpc.ClientConn, error) {
	s.logger.Debug(ctx, "get conns from client", "serviceName", serviceName)
	s.lock.Lock()
	defer s.lock.Unlock()
	conns := s.localConns[serviceName]
	if len(conns) == 0 {
		s.logger.Debug(ctx, "get conns from zk remote", "serviceName", serviceName)
		addrs, err := s.GetConnsRemote(ctx, serviceName)
		if err != nil {
			return nil, err
		}
		if len(addrs) == 0 {
			return nil, errs.New("addr is empty").WrapMsg("no conn for service", "serviceName",
				serviceName, "local conn", s.localConns, "ZkServers", s.ZkServers, "zkRoot", s.zkRoot)
		}
		for _, addr := range addrs {
			cc, err := grpc.DialContext(ctx, addr.Addr, append(s.options, opts...)...)
			if err != nil {
				s.logger.Error(context.Background(), "dialContext failed", err, "addr", addr.Addr, "opts", append(s.options, opts...))
				return nil, errs.WrapMsg(err, "DialContext failed", "addr.Addr", addr.Addr)
			}
			conns = append(conns, cc)
		}
		s.localConns[serviceName] = conns
	}
	return conns, nil
}

// GetConn方法用于获取指定服务的一个连接。使用的服务名称和选项来创建连接。
//
// 参数:
// ctx - 上下文。
// serviceName - 需要获取连接的服务名称。
// opts - gRPC拨号选项。
//
// 返回值:
// *grpc.ClientConn - 服务的连接。
// error - 获取连接过程中发生的错误。
func (s *ZkClient) GetConn(ctx context.Context, serviceName string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	newOpts := append(s.options, grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, s.balancerName)))
	s.logger.Debug(context.Background(), "get conn from client", "serviceName", serviceName)
	return grpc.DialContext(ctx, fmt.Sprintf("%s:///%s", s.scheme, serviceName), append(newOpts, opts...)...)
}

// GetSelfConnTarget方法用于获取当前客户端的连接目标地址。
//
// 返回值:
// string - 当前客户端的连接目标地址。
func (s *ZkClient) GetSelfConnTarget() string {
	return s.rpcRegisterAddr
}

// CloseConn方法用于关闭指定的gRPC客户端连接。
//
// 参数:
// conn - 需要关闭的gRPC客户端连接。
func (s *ZkClient) CloseConn(conn *grpc.ClientConn) {
	conn.Close()
}