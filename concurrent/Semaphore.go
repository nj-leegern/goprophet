package concurrent

import "time"

/*
	信号量
*/

type Semaphore struct {
	permits int      // 许可数量
	channel chan int // 通道
}

/* 创建信号量 */
func NewSemaphore(permits int) *Semaphore {
	return &Semaphore{channel: make(chan int, permits), permits: permits}
}

/* 获取许可 */
func (s *Semaphore) Acquire() {
	s.channel <- 0
}

/* 释放许可 */
func (s *Semaphore) Release() {
	<-s.channel
}

/* 尝试获取许可 */
func (s *Semaphore) TryAcquire() bool {
	select {
	case s.channel <- 0:
		return true
	default:
		return false
	}
}

/* 尝试指定时间内获取许可 */
func (s *Semaphore) TryAcquireOnTime(timeout time.Duration) bool {
	for i := 0; i < 2; i++ {
		select {
		case s.channel <- 0:
			return true
		default:
			if i == 0 {
				time.Sleep(timeout)
			} else {
				break
			}
		}
	}
	return false
}

/* 当前可用的许可数 */
func (s *Semaphore) AvailablePermits() int {
	return s.permits - len(s.channel)
}
