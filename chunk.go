package xgo

// Index is chunked index that contains from and to
type ChunkedIndex struct {
	From, To int
}

// ChunkedIndex divides index into chunks
func ChunkIndex(length int, chunkSize int) <-chan ChunkedIndex {
	if chunkSize == 0 {
		chunkSize = 1
	}

	ch := make(chan ChunkedIndex)

	go func() {
		defer close(ch)

		for i := 0; i < length; i += chunkSize {
			idx := ChunkedIndex{i, i + chunkSize}
			if length < idx.To {
				idx.To = length
			}
			ch <- idx
		}
	}()

	return ch
}

func Chunk[T any](list []T, chunkSize int) [][]T {
	var chunks [][]T

	for chunkSize < len(list) {
		list, chunks = list[chunkSize:], append(chunks, list[0:chunkSize:chunkSize])
	}
	chunks = append(chunks, list)

	return chunks
}
