/*
 * @Author: zhangkaiwei 1126763237@qq.com
 * @Date: 2024-05-08 21:43:05
 * @LastEditors: zhangkaiwei 1126763237@qq.com
 * @LastEditTime: 2024-05-08 21:45:21
 * @FilePath: \go-tools\checker\check_test.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package checker_test

import (
	"testing"

	"github.com/Meikwei/go-tools/checker"
	"github.com/Meikwei/go-tools/errs"
	"github.com/stretchr/testify/assert"
)

type mockChecker struct {
	err error
}

func (m mockChecker) Check() error {
	return m.err
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name      string
		arg       any
		wantError error
	}{
		{
			name:      "non-checker argument",
			arg:       "non-checker",
			wantError: nil,
		},
		{
			name:      "checker with no error",
			arg:       mockChecker{nil},
			wantError: nil,
		},
		{
			name:      "checker with generic error",
			arg:       mockChecker{errs.New("generic error")},
			wantError: errs.ErrArgs,
		},
		{
			name:      "checker with CodeError",
			arg:       mockChecker{errs.NewCodeError(400, "bad request")},
			wantError: errs.NewCodeError(400, "bad request"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checker.Validate(tt.arg)
			if tt.wantError != nil {
				assert.ErrorIs(t, err, tt.wantError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
