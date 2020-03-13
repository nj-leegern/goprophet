package utils

import (
	"fmt"
	"runtime"
	"strings"
)

/*
	异常处理
*/

/* 异常处理 */
func HandleException() (string, error) {
	if err := recover(); err != nil {
		stack := fmt.Sprintf("error: %s, stack trace: %s", err, printStack())
		if er, ok := err.(error); ok {
			return stack, er
		} else {
			return stack, fmt.Errorf("panic other exception")
		}
	}
	return "", nil
}

/* 输出错误信息 */
func GetStackTrace(e interface{}) string {
	items := make([]string, 0)
	if msg, ok := e.(string); ok {
		items = append(items, fmt.Sprintf("error: %s", msg))
	}
	items = append(items, fmt.Sprintf("stack trace: %s", printStack()))
	return strings.Join(items, "\n")
}

/* 打印调用栈 */
func printStack() string {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	return string(buf[:n])
}
