package utils

type Queue struct {
	Items []interface{}
}

func (q *Queue) Push(item interface{}) {
	q.Items = append(q.Items, item)
}

func (q *Queue) Pop() interface{} {
	if len(q.Items) == 0 {
		return nil
	}

	item := q.Items[0]
	q.Items = q.Items[1:]
	return item
}

func (q *Queue) NextItem() interface{} {
	if len(q.Items) == 0 {
		return nil
	}

	return q.Items[0]
}
