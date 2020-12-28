/*
	队列
*/

package mutil

// the base node
type Node struct {
	Data interface{}
	Next *Node
}

// the queue struct
type Queue struct {
	Head *Node
	End  *Node
}

// get last push data
func (q *Queue) GetEnd() interface{}{
	if q.End !=nil{
		return q.End.Data
	}
	return nil
}

// push data
func (q *Queue) Push(data interface{}) {
	n := &Node{
		Data: data,
		Next: nil,
	}

	if q.End == nil {
		q.Head = n
		q.End = n
	} else {
		q.End.Next = n
		q.End = n
	}

	return
}

// Pop data
func (q *Queue) Pop() (interface{}, bool) {
	if q.Head == nil {
		return nil, false
	}

	data := q.Head.Data
	q.Head = q.Head.Next
	if q.Head == nil {
		q.End = nil
	}
	return data, true
}

// GetSize get queue size
func (q *Queue) GetSize() (total int) {

	head := q.Head
	for ; head != nil; head = head.Next {
		total++
	}
	return total
}

