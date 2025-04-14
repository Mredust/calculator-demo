package test

import (
	"calculator-go/core/calculator"
	"calculator-go/service"
	"connectrpc.com/connect"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Mock connect.NewError to simulate error responses in tests.
func mockNewError(code connect.Code, err error) *connect.Error {
	return connect.NewError(code, err)
}

// Test Case: Test additional basic operations and edge cases
func TestCalculatorServer_Calculate_AdditionalCases(t *testing.T) {
	tests := []struct {
		expression  string
		expected    string
		expectedErr error
	}{

		{"1 + 1", "2", nil},
		{"2 - 1", "1", nil},
		{"3 * 2", "6", nil},
		{"6 / 2", "3", nil},
		{"5 / 2", "2.5", nil},
		{"10 / 3", "3.3333333333", nil},

		{"-1 + 1", "0", nil},
		{"-3 - 1", "-4", nil},
		{"-2 * 3", "-6", nil},
		{"-6 / 2", "-3", nil},
		{"1 / 0", "", mockNewError(connect.CodeInvalidArgument, fmt.Errorf("不能除以0"))},
		{"0 / 1", "0", nil},
		{"-1 / (-1)", "1", nil},
		{"2 * 0", "0", nil},
		{"-5 + 10", "5", nil},
		{"2.5 + 3.4", "5.9", nil},
		{"0.1 + 0.2", "0.3", nil},
		{"5 + 0.0", "5", nil},
		{"2.999999 + 1.000001", "4", nil},

		{"2 + (-3)", "-1", nil},
		{"1 + 1", "2", nil},
		{"-1 + 1", "0", nil},
		{"-1 + (-1)", "-2", nil},
		{"1 - 1", "0", nil},
		{"-1 - 1", "-2", nil},
		{"-1 - (-1)", "0", nil},
		{"1 * 1", "1", nil},
		{"-1 * (-1)", "1", nil},
		{"-1 * 1", "-1", nil},
		{"1 / 1", "1", nil},
		{"1 / (-1)", "-1", nil},
		{"-1 / (-1)", "1", nil},
		{"(-1) / (-1)", "1", nil},

		{"1 +", "1", nil},
		{"1 -", "1", nil},
		{"1 *", "1", nil},
		{"1 /", "1", nil},
		{"-1 *", "-1", nil},
		{"-1 +", "-1", nil},
		{"-1 -", "-1", nil},
		{"-1 /", "-1", nil},
		{"1 +.1", "1.1", nil},
		{"1 - .1", "0.9", nil},
		{"1 / .1", "10", nil},
		{"1 * .1", "0.1", nil},

		// 三个操作符以上的运算
		{"1 + 1 + 1", "3", nil},
		{"1 + 2 + 3 + 4", "10", nil},
		{"5 * 2 - 3 + 4", "11", nil},
		{"10 / 2 * 3 - 1", "14", nil}, // 注意运算顺序：((10/2)*3)-1
		{"2 + 3 * 4 - 5", "9", nil},   // 3*4优先
		{"-1 + 2 * 3 - 4", "1", nil},
		{"0.5 * 2 + 1.5 / 3", "1.5", nil},

		// 混合运算（加减乘除）
		{"2 + 3 * 4", "14", nil},   // 乘法优先
		{"(2 + 3) * 4", "20", nil}, // 括号优先
		{"10 - 2 / 2", "9", nil},   // 除法优先
		{"10 / (2 + 3)", "2", nil},
		{"(1 + 2) * (3 - 4)", "-3", nil},
		{"2 * 3 + 4 / 2", "8", nil},   // 先乘除后加减
		{"2 * (3 + 4) / 2", "7", nil}, // 括号优先

		// 括号优先级测试
		{"(1 + 2) * 3", "9", nil},
		{"1 + (2 * 3)", "7", nil},
		{"((1 + 2) * 3) - 4", "5", nil}, // 嵌套括号
		{"(5 + 3) / (2 * 2)", "2", nil},
		{"-1 * (2 + 3)", "-5", nil},
		{"(2.5 + 1.5) * 2", "8", nil},
		{"1 / (1 + 1)", "0.5", nil},

		// 边界和特殊案例
		{"1 + 2 * (3 + 4 / 2) - 5", "6", nil},      // 复合嵌套
		{"(1 + (2 * (3 - (4 / 2))))", "3", nil},    // 多重嵌套
		{"0.1 * 0.2 + 0.3 - 0.4 / 2", "0.12", nil}, // 小数混合运算
		{"(1 + 2) * (3 + 4) / (5 - 2)", "7", nil},  // 多组括号
		{"1 + 2 - 3 * 4 / 5 + 6 - 7", "-0.4", nil}, // 长表达式
	}

	for _, tt := range tests {
		t.Run(tt.expression, func(t *testing.T) {
			req := &connect.Request[calculator.CalculationRequest]{
				Msg: &calculator.CalculationRequest{
					Expression: tt.expression,
				},
			}

			server := service.NewCalculatorServer()
			resp, err := server.Calculate(context.Background(), req)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, resp.Msg.Result)
			}
		})
	}
}
