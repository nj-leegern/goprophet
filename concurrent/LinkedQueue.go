package concurrent

import (
	"container/list"
	"sync"
)

/*
	安全的链表队列
*/

type LinkedQueue struct {
	store *list.List
	mutex *sync.Mutex
}

/* 创建队列 */
func NewLinkedQueue() *LinkedQueue {
	return &LinkedQueue{store: list.New()}
}

/* 压入队列 */
func (q *LinkedQueue) Push(e interface{}) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.store.PushBack(e)
	return nil
}

/* 弹出元素 */
func (q *LinkedQueue) Pop() interface{} {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.store.Front().Value
}

/* 队列大小 */
func (q *LinkedQueue) Size() int64 {
	return int64(q.store.Len())
}
