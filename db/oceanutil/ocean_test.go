/*
 * @Author: zhangkaiwei 1126763237@qq.com
 * @Date: 2024-05-06 18:46:47
 * @LastEditors: zhangkaiwei 1126763237@qq.com
 * @LastEditTime: 2024-05-09 19:39:45
 * @FilePath: \go-tools\db\oceanutil\ocean_test.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package oceanutil

import (
	"context"
	"testing"
)

func Test_connectOcean(t *testing.T) {
	ctx:=context.Background();
	config:=&Config{
		Dns:"",// 连接服务器的URI
		Host:"127.0.0.1", // 服务器主机地址
		Username:"root", // 连接服务器的用户名
		TenantName:"sys", // 租户名称
		ClusterName:"",// 集群名称
		Password:"",// 连接服务器的密码
		Port:2881,  // 服务器端口
		SchemaName:"test", // 使用的数据库模式名称
		DataType:"mysql",// 数据类型，mysql、oracle
		Charset:"utf8mb4", // 字符集名称
		MaxRetry:3,   // 最大重试次数
	}
	cli,err:=NewOceanbase(ctx,config)
	if err != nil {
		t.Errorf("connect oceanbase error:%v",err)
		return
	}
	t.Logf("connect oceanbase success:%v",cli)

}