package xgo

// Chunk is chunked index that contains from and to
type Chunk struct {
	From, To int
}

// SplitChunks splits chunks from slice length
func SplitChunks(length int, chunkSize int) <-chan Chunk {
	if chunkSize == 0 {
		chunkSize = 1
	}

	ch := make(chan Chunk)

	go func() {
		defer close(ch)

		for i := 0; i < length; i += chunkSize {
			idx := Chunk{i, i + chunkSize}
			if length < idx.To {
				idx.To = length
			}
			ch <- idx
		}
	}()

	return ch
}
