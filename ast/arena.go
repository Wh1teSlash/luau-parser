package ast

const defaultChunkSize = 1024

type TypedArena[T any] struct {
	chunks   [][]T
	chunkIdx int
	itemIdx  int
}

func NewTypedArena[T any]() *TypedArena[T] {
	return &TypedArena[T]{
		chunks: [][]T{make([]T, defaultChunkSize)},
	}
}

func (a *TypedArena[T]) Alloc() *T {
	if a.itemIdx >= defaultChunkSize {
		a.chunks = append(a.chunks, make([]T, defaultChunkSize))
		a.chunkIdx++
		a.itemIdx = 0
	}

	ptr := &a.chunks[a.chunkIdx][a.itemIdx]
	a.itemIdx++
	return ptr
}

func (a *TypedArena[T]) Reset() {
	a.chunks = a.chunks[:1]
	a.chunkIdx = 0
	a.itemIdx = 0
}
