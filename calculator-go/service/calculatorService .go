package service

import (
	"calculator-go/core/calculator"
	"connectrpc.com/connect"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

import (
	"context"
)

// CalculatorServer 实现计算器服务
type CalculatorServer struct{}

// NewCalculatorServer 创建服务实例
func NewCalculatorServer() *CalculatorServer {
	return &CalculatorServer{}
}

func (s *CalculatorServer) Calculate(
	ctx context.Context,
	req *connect.Request[calculator.CalculationRequest],
) (*connect.Response[calculator.CalculationResponse], error) {
	expression := req.Msg.GetExpression()
	// 清除所有空格
	expression = strings.ReplaceAll(expression, " ", "")

	// 预处理表达式：处理负数和末尾操作符
	expression, err := preprocessExpression(expression)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			fmt.Errorf("表达式无效"),
		)
	}

	// 使用正则表达式来匹配数字、操作符和括号
	re := regexp.MustCompile(`(\d+\.?\d*|\.\d+)|([+\-*/])|([()])`)
	matches := re.FindAllStringSubmatch(expression, -1)

	// 用于存储数值和操作符的栈
	var numStack []float64
	var opStack []string

	// 定义操作符优先级
	priority := map[string]int{
		"+": 1,
		"-": 1,
		"*": 2,
		"/": 2,
	}

	// 逐个处理匹配到的数字、操作符和括号
	for _, match := range matches {
		if match[1] != "" { // 如果是数字
			num, err := strconv.ParseFloat(match[1], 64)
			if err != nil {
				return nil, connect.NewError(
					connect.CodeInvalidArgument,
					fmt.Errorf("无法解析数字"),
				)
			}
			numStack = append(numStack, num)
			continue
		}

		if match[3] == "(" { // 如果是左括号
			opStack = append(opStack, match[3])
			continue
		}

		if match[3] == ")" { // 如果是右括号
			// 弹出栈顶操作符并进行计算，直到遇到左括号
			for len(opStack) > 0 && opStack[len(opStack)-1] != "(" {
				err := processOperation(&numStack, &opStack, priority)
				if err != nil {
					return nil, err
				}
			}
			// 弹出左括号
			if len(opStack) == 0 || opStack[len(opStack)-1] != "(" {
				return nil, connect.NewError(
					connect.CodeInvalidArgument,
					fmt.Errorf("括号不匹配"),
				)
			}
			opStack = opStack[:len(opStack)-1]
			continue
		}

		// 如果是操作符
		if match[2] != "" {
			// 操作符压栈
			for len(opStack) > 0 && priority[opStack[len(opStack)-1]] >= priority[match[2]] {
				err := processOperation(&numStack, &opStack, priority)
				if err != nil {
					return nil, err
				}
			}
			opStack = append(opStack, match[2])
		}
	}
	for len(opStack) > 0 {
		err := processOperation(&numStack, &opStack, priority)
		if err != nil {
			return nil, err
		}
	}

	// 最终栈中应该只有一个结果
	if len(numStack) != 1 {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			fmt.Errorf("表达式无效"),
		)
	}

	// 格式化结果
	result := numStack[0]
	formattedResult := formatResult(result)

	// 返回结果
	return connect.NewResponse(&calculator.CalculationResponse{
		Result: formattedResult,
	}), nil
}

// 预处理表达式
func preprocessExpression(expr string) (string, error) {
	// 首先检查括号是否匹配
	balance := 0
	for _, ch := range expr {
		if ch == '(' {
			balance++
		} else if ch == ')' {
			balance--
			if balance < 0 {
				return "", fmt.Errorf("括号不匹配")
			}
		}
	}
	if balance != 0 {
		return "", fmt.Errorf("括号不匹配")
	}

	// 处理负数情况
	expr = strings.ReplaceAll(expr, "(-", "(0-")
	if strings.HasPrefix(expr, "-") {
		expr = "0" + expr
	}

	// 处理多个连续负号的情况
	expr = strings.ReplaceAll(expr, "--", "+")
	expr = strings.ReplaceAll(expr, "+-", "-")
	expr = strings.ReplaceAll(expr, "*-", "*-1*")
	expr = strings.ReplaceAll(expr, "/-", "/-1*")

	// 处理末尾操作符
	for len(expr) > 0 && strings.ContainsAny(string(expr[len(expr)-1]), "+-*/") {
		expr = expr[:len(expr)-1]
	}

	return expr, nil
}

// 处理单个运算操作
func processOperation(numStack *[]float64, opStack *[]string, priority map[string]int) *connect.Error {
	if len(*opStack) == 0 || len(*numStack) < 2 {
		return connect.NewError(
			connect.CodeInvalidArgument,
			fmt.Errorf("表达式无效"),
		)
	}

	op := (*opStack)[len(*opStack)-1]
	*opStack = (*opStack)[:len(*opStack)-1]

	num2 := (*numStack)[len(*numStack)-1]
	*numStack = (*numStack)[:len(*numStack)-1]
	num1 := (*numStack)[len(*numStack)-1]
	*numStack = (*numStack)[:len(*numStack)-1]

	result, err := performOperation(num1, num2, op)
	if err != nil {
		return err
	}

	*numStack = append(*numStack, result)
	return nil
}

// 执行运算
func performOperation(num1, num2 float64, op string) (float64, *connect.Error) {
	switch op {
	case "+":
		return num1 + num2, nil
	case "-":
		return num1 - num2, nil
	case "*":
		return num1 * num2, nil
	case "/":
		if num2 == 0 {
			return 0, connect.NewError(
				connect.CodeInvalidArgument,
				fmt.Errorf("不能除以0"),
			)
		}
		return num1 / num2, nil
	default:
		return 0, connect.NewError(
			connect.CodeInvalidArgument,
			fmt.Errorf("未知操作符"),
		)
	}
}

// 格式化结果
func formatResult(result float64) string {
	// 处理浮点数精度问题
	result = roundToPrecision(result, 10)

	// 如果结果是整数，显示为整数；否则显示必要的小数位
	if result == float64(int(result)) {
		return fmt.Sprintf("%d", int(result))
	}

	// 去除末尾的0
	str := fmt.Sprintf("%.10f", result)
	str = strings.TrimRight(str, "0")
	str = strings.TrimRight(str, ".")

	return str
}

// 四舍五入到指定精度
func roundToPrecision(num float64, precision int) float64 {
	shift := math.Pow(10, float64(precision))
	return math.Round(num*shift) / shift
}
