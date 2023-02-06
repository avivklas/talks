package main

import (
	"bufio"
	"io"
	"strings"
)

func main() {

}

func count(f io.ReadSeeker, word string) (total int, err error) {
	f.Seek(0, 0)
	scanner := bufio.NewScanner(f)
	total = 0
	for scanner.Scan() {
		total += strings.Count(scanner.Text(), word)
	}

	err = scanner.Err()

	return
}
