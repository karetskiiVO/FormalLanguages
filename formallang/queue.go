package formallang

type queue []any

func (q *queue) Push(x any) {
	*q = append(*q, x)
}

func (q queue) Size() int {
	return len(q)
}

func (q *queue) Pop() {
	if q.Size() <= 0 {
		return
	}

	*q = (*q)[1:]
}

func (q *queue) Top() any {
	return (*q)[0]
}
