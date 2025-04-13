package main

import (
	"testing"
)

// https://github.com/labstack/echo/blob/de44c53a5b16f7dca451f337f7221a1448c92007/middleware/util_test.go#L101

func TestRandomString(t *testing.T) {
	testCases := []struct {
		name   string
		input  uint8
		expect string
	}{
		{
			name:  "ok, 16",
			input: 16,
		},
		{
			name:  "ok, 32",
			input: 32,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			uid := randomString(tc.input)
			got := len(uid)
			want := int(tc.input)
			if got != want {
				t.Errorf("want %d, got %d", want, got)
			}
		})
	}
}

func TestRandomStringBias(t *testing.T) {
	t.Parallel()
	const slen = 33
	const loop = 100000

	counts := make(map[rune]int)
	var count int64

	for range loop {
		s := randomString(slen)
		got := len(s)
		if got != slen {
			t.Fatalf("want %d, got %d", slen, got)
		}
		for _, b := range s {
			counts[b]++
			count++
		}
	}

	got := len(counts)
	if got != randomStringCharsetLen {
		t.Fatalf("want %d, got %d", randomStringCharsetLen, got)
	}

	avg := float64(count) / float64(len(counts))
	for k, n := range counts {
		diff := float64(n) / avg
		if diff < 0.95 || diff > 1.05 {
			t.Errorf("Bias on '%c': expected average %f, got %d", k, avg, n)
		}
	}
}

func TestGenerateID(t *testing.T) {
	id := generateID()
	got := len(id)
	want := 32
	if got != want {
		t.Errorf("want %d, got %d", want, got)
	}
}
