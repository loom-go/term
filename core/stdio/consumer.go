package stdio

type consumer func([]byte) (consumed int, complete bool)

type BufferedConsumer struct {
	buf      []byte
	consumer consumer
}

func NewBufferedConsumer(consumer consumer) *BufferedConsumer {
	return &BufferedConsumer{
		consumer: consumer,
	}
}

func (b *BufferedConsumer) Feed(data []byte) {
	b.buf = append(b.buf, data...)

	for len(b.buf) > 0 {
		consumed, complete := b.consumer(b.buf)
		if !complete {
			return // wait for more
		}

		if consumed == 0 {
			consumed = 1 // force progress to avoid infinite loop, just in case
		}

		b.buf = b.buf[consumed:]
	}
}
