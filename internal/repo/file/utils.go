package file

import (
	"bufio"
	"fmt"

	"github.com/spf13/afero"
)

func WriteLines(cont []string, fname string, fs afero.Fs) error {
	f, err := fs.Create(fname)
	if err != nil {
		return fmt.Errorf("failed create file %s: %w", fname, err)
	}
	defer f.Close()

	buf := bufio.NewWriter(f)
	for _, ln := range cont {
		_, err := buf.WriteString(ln + "\n")
		if err != nil {
			return fmt.Errorf("failed write to file %s: %w", fname, err)
		}
	}
	if err := buf.Flush(); err != nil {
		return fmt.Errorf("failed write to file %s: %w", fname, err)
	}
	return nil
}

func ReadLines(fname string, fs afero.Fs) ([]string, error) {
	file, err := fs.Open(fname)
	if err != nil {
		return nil, fmt.Errorf("failed open file %s: %w", fname, err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
