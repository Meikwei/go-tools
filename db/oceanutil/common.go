/*
 * @Author: zhangkaiwei 1126763237@qq.com
 * @Date: 2024-05-05 11:11:23
 * @LastEditors: zhangkaiwei 1126763237@qq.com
 * @LastEditTime: 2024-05-08 20:09:52
 * @FilePath: \go-tools\db\oceanutil\common.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
/*
 * @Author: zhangkaiwei 1126763237@qq.com
 * @Date: 2024-05-05 11:11:23
 * @LastEditors: zhangkaiwei 1126763237@qq.com
 * @LastEditTime: 2024-05-06 23:23:36
 * @FilePath: \go-tools\db\oceanutil\common.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package oceanutil

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Meikwei/go-tools/errs"
	"golang.org/x/text/transform"
	"gorm.io/gorm"
)

// 将字符串转换为指定的编码格式
func transformString(str string, encoder transform.Transformer) string {
	result, _, err := transform.String(encoder, str)
	if err != nil {
			return str
	}
	return result
}
/**
conn := "{username}:{password}@tcp({hostname}:{port})/{dbname}"
conn := "root:@tcp(127.0.0.1:2881)/test"
参数说明：
username：取自 -u 参数，租户的连接用户名，格式为 用户@租户#集群名称，集群的默认租户是 'sys'，租户的默认管理员用户是 'root'。直连数据库时不填写集群名称，通过 ODP 连接时需要填写。
password：取自 -p 参数，用户密码。
hostname：取自 -h 参数，OceanBase 数据库连接地址，有时候是 ODP 地址。
port：取自 -P 参数，OceanBase 数据库连接端口，也是 ODP 的监听端口。
dbname：取自 -D 参数，需要访问的数据库名称。
*/
const oceanBaseURIFormat = "%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local"
func buildOceanURI(config *Config) (string,error) {
	if config==nil{
		return "",fmt.Errorf("config is null")
	}
	var credentialsBuilder strings.Builder
	if config.Username != "" {
		if config.TenantName!=""&&config.ClusterName!=""{
			credentialsBuilder.WriteString(fmt.Sprintf("%s@%s#%s", config.Username,config.TenantName,config.ClusterName))
		}else if config.TenantName!=""{
			credentialsBuilder.WriteString(fmt.Sprintf("%s@%s", config.Username, config.TenantName))
		}else{
			credentialsBuilder.WriteString(config.Username)
		}
	}else {
		// 如果 Username 为空，直接返回错误，避免构建 credentials
		return "", fmt.Errorf("username cannot be empty")
	}
 	// 构造 OceanBase URI
	return fmt.Sprintf(oceanBaseURIFormat, credentialsBuilder.String(), config.Password, config.Host, config.Port, config.SchemaName, config.Charset), nil



}

// shouldRetry 判断一个错误是否应该触发重试。
// 参数:
//   ctx context.Context: 上下文，用于控制函数的生命周期。
//   err error: 执行操作过程中发生的错误。
// 返回值:
//   bool: 如果错误应该触发重试则返回 true，否则返回 false。
func shouldRetry(ctx context.Context, err error) bool {
	select {
	case <-ctx.Done():
		// 如果上下文被取消或超时，则不重试
		return false
	default:
		// 检查错误是否为 mongo.CommandError 类型
		// if cmdErr, ok := err.(mongo.CommandError); ok {
		// 	// 不重试特定的错误代码
		// 	return cmdErr.Code != 13 && cmdErr.Code != 18
		// }
		// 对于非 mongo.CommandError 类型的错误，默认重试
		return true
	}
}

// initializeDBWithRetry 尝试初始化数据库连接，包含重试逻辑
func initializeDBWithRetry(ctx context.Context, config *Config) (*gorm.DB, error) {

	var (
		db *gorm.DB
		err error
	)
	for i := 0; i < config.MaxRetry; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			db, err = connectOcean(config)
			if err == nil {
				return db, nil
			}
			errs.WrapMsg(err,"Database connection attempt %d",i+1)
			time.Sleep(time.Second << uint(i)) // 使用指数退避策略
		}
	}

	return nil, err
}