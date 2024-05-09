/*
 * @Author: zhangkaiwei 1126763237@qq.com
 * @Date: 2024-05-08 21:28:57
 * @LastEditors: zhangkaiwei 1126763237@qq.com
 * @LastEditTime: 2024-05-08 21:33:33
 * @FilePath: \go-tools\utils\mageutil\usage-guide\magefile_unix.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
//go:build mage && !windows
// +build mage,!windows

package main

import (
	"github.com/Meikwei/go-tools/utils/mageutil"
	"syscall"
)

func setMaxOpenFiles() error {
	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		return err
	}
	rLimit.Max = uint64(mageutil.MaxFileDescriptors)
	rLimit.Cur = uint64(mageutil.MaxFileDescriptors)
	return syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
}
