package main

import (
	"bufio"
	"crypto/rand"
	"io"
	"sync"
)

// https://github.com/labstack/echo/blob/de44c53a5b16f7dca451f337f7221a1448c92007/middleware/util.go#L72

// https://tip.golang.org/doc/go1.19#:~:text=Read%20no%20longer%20buffers%20random%20data%20obtained%20from%20the%20operating%20system%20between%20calls
var randomReaderPool = sync.Pool{New: func() any {
	return bufio.NewReader(rand.Reader)
}}

const (
	randomStringCharset    = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	randomStringCharsetLen = 52 // len(randomStringCharset)
	randomStringMaxByte    = 255 - (256 % randomStringCharsetLen)
)

func randomString(length uint8) string {
	reader := randomReaderPool.Get().(*bufio.Reader)
	defer randomReaderPool.Put(reader)

	b := make([]byte, length)
	r := make([]byte, length+(length/4)) // perf: avoid read from rand.Reader many times
	var i uint8 = 0

	// security note:
	// we can't just simply do b[i]=randomStringCharset[rb%len(randomStringCharset)],
	// len(len(randomStringCharset)) is 52, and rb is [0, 255], 256 = 52 * 4 + 48.
	// make the first 48 characters more possibly to be generated then others.
	// So we have to skip bytes when rb > randomStringMaxByte

	for {
		_, err := io.ReadFull(reader, r)
		if err != nil {
			panic("unexpected error happened when reading from bufio.NewReader(crypto/rand.Reader)")
		}
		for _, rb := range r {
			if rb > randomStringMaxByte {
				// Skip this number to avoid bias.
				continue
			}
			b[i] = randomStringCharset[rb%randomStringCharsetLen]
			i++
			if i == length {
				return string(b)
			}
		}
	}
}

func generateID() string {
	return randomString(32)
}
