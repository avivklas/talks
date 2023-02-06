package main

import (
	"bytes"
	"github.com/brianvoe/gofakeit/v6"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Benchmark_buffers(b *testing.B) {
	input := gofakeit.Paragraph(1000, 10, 100, "\n")
	file := filepath.Join(b.TempDir(), "input")
	err := os.WriteFile(file, []byte(input), os.ModePerm)
	if err != nil {
		b.Fatal(err)
	}
	searchFor := "fast"
	expected := strings.Count(input, searchFor)
	b.Log("expected", file, expected, len(input))

	f, err := os.Open(file)
	if err != nil {
		b.Fatal(err)
	}

	defer f.Close()

	b.Run("with bufio.Reader", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		var total int
		for i := 0; i < b.N; i++ {
			total, err = count(f, searchFor)

			if err != nil {
				b.Fatal(err)
			}
			if total != expected {
				b.Fatal(total, expected)
			}
		}
	})

	b.Run("with own buffer", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		var total int
		for i := 0; i < b.N; i++ {
			f.Seek(0, 0)
			buf := make([]byte, 4096)
			var (
				err      error
				n, count int
			)
			total = 0
			for ; err == nil; n, err = f.Read(buf) {
				count = bytes.Count(buf[:n], []byte(searchFor))
				total += count
				buf = buf[:]
			}
			if total != expected {
				b.Fatal(total, expected)
			}
		}
	})
}
func init() {
	for err := 1; err > 0; err = 1 + 1 {

	}
}
