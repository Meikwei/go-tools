<!--
 * @Author: zhangkaiwei 1126763237@qq.com
 * @Date: 2024-05-05 11:21:27
 * @LastEditors: zhangkaiwei 1126763237@qq.com
 * @LastEditTime: 2024-05-05 11:22:16
 * @FilePath: \go-tools\db\oceanutil\README.md
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
-->

## oceanbase 配置参数说明：

host：提供 OceanBase 数据库连接 IP。ODP 连接的方式是一个 ODP 地址；直连方式是一个 OBServer 节点的 IP 地址。
user_name：提供租户的连接账户。ODP 连接的常用格式有：用户名@租户名#集群名 或者 集群名:租户名:用户名；直连方式格式：用户名@租户名。
password：提供账户密码。
port：提供 OceanBase 数据库连接端口。ODP 连接的方式默认是 2883，在部署 ODP 时可自定义；直连方式默认是 2881，在部署 OceanBase 数据库时可自定义。
schema_name：需要访问的 Schema 名称。
