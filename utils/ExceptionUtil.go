package utils

import (
	"fmt"
	"runtime"
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

/* 打印调用栈 */
func printStack() string {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	return string(buf[:n])
}
