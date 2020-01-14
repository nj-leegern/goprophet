package utils

import (
	"fmt"
	"runtime"
)

/*
	异常处理
*/

/* 异常处理 */
func HandleException() string {
	if err := recover(); err != nil {
		return fmt.Sprintf("error: %s, stack trace: %s", err, printStack())
	}
	return ""
}

/* 打印调用栈 */
func printStack() string {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	return string(buf[:n])
}
