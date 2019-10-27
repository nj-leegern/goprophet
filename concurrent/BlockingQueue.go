package concurrent

/*
	阻塞队列
*/

type BlockingQueue struct {
	queue chan interface{}
}

/* 创建队列 */
func NewBlockingQueue(capacity int64) *BlockingQueue {
	if capacity <= 0 {
		capacity = 16
	}
	return &BlockingQueue{queue: make(chan interface{}, capacity)}
}

/* 存放元素 */
func (q *BlockingQueue) Put(e interface{}) {
	q.queue <- e
}

/* 取出元素 */
func (q *BlockingQueue) Take() interface{} {
	return <-q.queue
}
