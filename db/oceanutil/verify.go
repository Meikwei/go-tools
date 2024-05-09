/*
 * @Author: zhangkaiwei 1126763237@qq.com
 * @Date: 2024-05-06 21:39:37
 * @LastEditors: zhangkaiwei 1126763237@qq.com
 * @LastEditTime: 2024-05-08 20:12:46
 * @FilePath: \go-tools\db\oceanutil\verify.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package oceanutil

import (
	"github.com/Meikwei/go-tools/errs"
)

func (c *Config) ValidateAndSetDefaults() error {
	if c.Host == "" {
		return errs.New("either Host must be provided")
	}
	if c.SchemaName== "" {
		return errs.New("database is required")
	}

	if c.Dns == "" {
		dns,err:=buildOceanURI(c)
		if err != nil {
			return errs.New("dns is not fmt")
		}
		c.Dns = dns
	}
	return nil
}