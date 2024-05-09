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
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/Meikwei/go-tools/errs"
	"github.com/Meikwei/go-tools/log"
	"github.com/go-zookeeper/zk"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

// 定义常量
const (
	defaultFreq = time.Minute * 30 // 默认刷新频率
	timeout     = 5               // 默认超时时间
)

// ZkClient 定义了ZooKeeper客户端的结构体
type ZkClient struct {
	ZkServers []string          // ZooKeeper服务器地址列表
	zkRoot    string            // ZooKeeper的根路径
	username  string            // ZooKeeper的用户名
	password  string            // ZooKeeper的密码

	rpcRegisterName string         // RPC注册名称
	rpcRegisterAddr string         // RPC注册地址
	isRegistered    bool           // 是否已注册
	scheme          string         // 方案

	timeout   int           // 超时时间
	conn      *zk.Conn      // ZooKeeper连接
	eventChan <-chan zk.Event // ZooKeeper事件通道
	node      string        // 节点
	ticker    *time.Ticker  // 定时器

	lock    sync.Locker    // 锁
	options []grpc.DialOption // gRPC拨号选项

	resolvers           map[string]*Resolver // 解析器映射
	localConns          map[string][]*grpc.ClientConn // 本地连接映射
	cancel              context.CancelFunc // 取消函数
	isStateDisconnected bool               // 是否断开状态
	balancerName        string             // 平衡器名称

	logger log.Logger // 日志记录器
}

// NewZkClient 初始化一个新的ZkClient实例并建立与ZooKeeper的连接
func NewZkClient(ZkServers []string, scheme string, options ...ZkOption) (*ZkClient, error) {
	// 初始化客户端实例并应用选项
	client := &ZkClient{
		ZkServers:  ZkServers,
		zkRoot:     "/",
		scheme:     scheme,
		timeout:    timeout,
		localConns: make(map[string][]*grpc.ClientConn),
		resolvers:  make(map[string]*Resolver),
		lock:       &sync.Mutex{},
		logger:     nilLog{},
	}
	for _, option := range options {
		option(client)
	}

	// 建立与ZooKeeper的连接并进行认证
	conn, eventChan, err := zk.Connect(ZkServers, time.Duration(client.timeout)*time.Second, zk.WithLogger(nilLog{}))
	if err != nil {
		return nil, errs.WrapMsg(err, "connect failed", "ZkServers", ZkServers)
	}

	ctx, cancel := context.WithCancel(context.Background())
	client.cancel = cancel
	client.ticker = time.NewTicker(defaultFreq)

	// 如果提供了用户名和密码，则进行认证
	if client.username != "" && client.password != "" {
		auth := []byte(client.username + ":" + client.password)
		if err := conn.AddAuth("digest", auth); err != nil {
			conn.Close()
			return nil, errs.WrapMsg(err, "AddAuth failed", "username", client.username, "password", client.password)
		}
	}

	client.zkRoot += scheme
	client.eventChan = eventChan
	client.conn = conn

	// 确保根节点的存在，如果不存在则创建
	if err := client.ensureRoot(); err != nil {
		conn.Close()
		return nil, err
	}

	resolver.Register(client)
	go client.refresh(ctx)
	go client.watch(ctx)

	return client, nil
}

// Close 关闭ZkClient实例，释放资源
func (s *ZkClient) Close() {
	s.logger.Info(context.Background(), "close zk called")
	s.cancel()
	s.ticker.Stop()
	s.conn.Close()
}

// ensureAndCreate 确保节点存在，如果不存在则创建
func (s *ZkClient) ensureAndCreate(node string) error {
	exists, _, err := s.conn.Exists(node)
	if err != nil {
		return errs.WrapMsg(err, "Exists failed", "node", node)
	}
	if !exists {
		_, err = s.conn.Create(node, []byte(""), 0, zk.WorldACL(zk.PermAll))
		if err != nil && err != zk.ErrNodeExists {
			return errs.WrapMsg(err, "Create failed", "node", node)
		}
	}
	return nil
}

// refresh 定期刷新本地连接和解析器状态
func (s *ZkClient) refresh(ctx context.Context) {
	for range s.ticker.C {
		s.logger.Debug(ctx, "zk refresh local conns")
		s.lock.Lock()
		for rpcName := range s.resolvers {
			s.flushResolver(rpcName)
		}
		for rpcName := range s.localConns {
			delete(s.localConns, rpcName)
		}
		s.lock.Unlock()
		s.logger.Debug(ctx, "zk refresh local conns success")
	}
}

// flushResolverAndDeleteLocal 清空指定服务的解析器并删除本地连接
func (s *ZkClient) flushResolverAndDeleteLocal(serviceName string) {
	s.logger.Debug(context.Background(), "zk start flush", "serviceName", serviceName)
	s.flushResolver(serviceName)
	delete(s.localConns, serviceName)
}

// flushResolver 立即刷新指定服务的解析器
func (s *ZkClient) flushResolver(serviceName string) {
	r, ok := s.resolvers[serviceName]
	if ok {
		r.ResolveNowZK(resolver.ResolveNowOptions{})
	}
}

// GetZkConn 返回ZooKeeper的连接实例
func (s *ZkClient) GetZkConn() *zk.Conn {
	return s.conn
}

// GetRootPath 返回ZooKeeper的根路径
func (s *ZkClient) GetRootPath() string {
	return s.zkRoot
}

// GetNode 返回当前节点
func (s *ZkClient) GetNode() string {
	return s.node
}

// ensureRoot 确保根节点的存在
func (s *ZkClient) ensureRoot() error {
	return s.ensureAndCreate(s.zkRoot)
}

// ensureName 确保指定名称的节点存在
func (s *ZkClient) ensureName(rpcRegisterName string) error {
	return s.ensureAndCreate(s.getPath(rpcRegisterName))
}

// getPath 返回指定名称的节点的完整路径
func (s *ZkClient) getPath(rpcRegisterName string) string {
	return s.zkRoot + "/" + rpcRegisterName
}

// getAddr 将主机和端口组合成一个地址字符串
func (s *ZkClient) getAddr(host string, port int) string {
	return net.JoinHostPort(host, strconv.Itoa(port))
}

// AddOption 添加gRPC拨号选项
func (s *ZkClient) AddOption(opts ...grpc.DialOption) {
	s.options = append(s.options, opts...)
}

// GetClientLocalConns 返回本地连接映射
func (s *ZkClient) GetClientLocalConns() map[string][]*grpc.ClientConn {
	return s.localConns
}
