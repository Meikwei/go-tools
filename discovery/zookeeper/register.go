/*
 * @Author: zhangkaiwei 1126763237@qq.com
 * @Date: 2024-04-29 21:52:30
 * @LastEditors: zhangkaiwei 1126763237@qq.com
 * @LastEditTime: 2024-05-01 12:22:52
 * @FilePath: \tools\discovery\zookeeper\register.go
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
	"time"

	"github.com/Meikwei/go-tools/errs"
	"github.com/go-zookeeper/zk"
	"google.golang.org/grpc"
)

// CreateRpcRootNodes 在ZooKeeper中为指定的服务名称创建根节点。
// serviceNames: 需要创建根节点的服务名称列表。
// 返回值: 如果创建过程中出现错误（不包括节点已存在的情况），则返回错误信息；否则返回nil。
func (s *ZkClient) CreateRpcRootNodes(serviceNames []string) error {
	for _, serviceName := range serviceNames {
		if err := s.ensureName(serviceName); err != nil && err != zk.ErrNodeExists {
			return err
		}
	}
	return nil
}

// CreateTempNode 在ZooKeeper中创建一个临时有序节点。
// rpcRegisterName: 注册名称。
// addr: 服务的地址。
// 返回值: 创建成功的节点名称和可能出现的错误。
func (s *ZkClient) CreateTempNode(rpcRegisterName, addr string) (node string, err error) {
	node, err = s.conn.CreateProtectedEphemeralSequential(
		s.getPath(rpcRegisterName)+"/"+addr+"_",
		[]byte(addr),
		zk.WorldACL(zk.PermAll),
	)
	if err != nil {
		return "", errs.WrapMsg(err, "CreateProtectedEphemeralSequential failed", "path", s.getPath(rpcRegisterName)+"/"+addr+"_")
	}
	return node, nil
}

// Register 在ZooKeeper中注册服务，并建立与该服务的连接。
// rpcRegisterName: 注册名称。
// host: 服务主机地址。
// port: 服务端口。
// opts: grpc连接选项。
// 返回值: 如果注册过程中出现错误，则返回错误信息；否则返回nil。
func (s *ZkClient) Register(rpcRegisterName, host string, port int, opts ...grpc.DialOption) error {
	if err := s.ensureName(rpcRegisterName); err != nil {
		return err
	}
	addr := s.getAddr(host, port)
	_, err := grpc.Dial(addr, opts...)
	if err != nil {
		return errs.WrapMsg(err, "grpc dial error", "addr", addr)
	}
	node, err := s.CreateTempNode(rpcRegisterName, addr)
	if err != nil {
		return err
	}
	s.rpcRegisterName = rpcRegisterName
	s.rpcRegisterAddr = addr
	s.node = node
	s.isRegistered = true
	return nil
}

// UnRegister 在ZooKeeper中注销服务，断开与该服务的连接，并清除相关注册信息。
// 返回值: 如果注销过程中出现错误，则返回错误信息；否则返回nil。
func (s *ZkClient) UnRegister() error {
	s.lock.Lock()
	defer s.lock.Unlock()
	err := s.conn.Delete(s.node, -1)
	if err != nil {
		return errs.WrapMsg(err, "delete node error", "node", s.node)
	}
	time.Sleep(time.Second)
	s.node = ""
	s.rpcRegisterName = ""
	s.rpcRegisterAddr = ""
	s.isRegistered = false
	s.localConns = make(map[string][]*grpc.ClientConn)
	s.resolvers = make(map[string]*Resolver)
	return nil
}
