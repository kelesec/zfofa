package output

import (
	"bufio"
	"fmt"
	"os"
)

type txtWriter struct {
	file *os.File
}

func CreateTxtWriter(filename string, append bool) (*txtWriter, error) {
	flag := os.O_CREATE | os.O_WRONLY
	if append {
		flag |= os.O_APPEND
	}

	file, err := os.OpenFile(filename, flag, 0666)
	if err != nil {
		return nil, err
	}
	return &txtWriter{file: file}, nil
}

// WriteString 写入一个字符串
func (f *txtWriter) WriteString(format string, a ...any) (int, error) {
	return f.file.WriteString(fmt.Sprintf(format, a...))
}

// WriteStringLn 写入一个字符串
func (f *txtWriter) WriteStringLn(format string, a ...any) (int, error) {
	return f.WriteString(format+"\n", a...)
}

// WriteStrings 写入一个字符串数组
func (f *txtWriter) WriteStrings(contents []string) (int, error) {
	fw := bufio.NewWriter(f.file)
	index := 0

	for _, content := range contents {
		_, err := fw.WriteString(content)
		if err != nil {
			return index, err
		}

		fw.Flush()
		index++
	}

	return -1, nil
}

// WriteStringsLn 写入一个字符串数组
func (f *txtWriter) WriteStringsLn(contents []string) (int, error) {
	fw := bufio.NewWriter(f.file)
	index := 0

	for _, content := range contents {
		_, err := fw.WriteString(content + "\n")
		if err != nil {
			return index, err
		}

		fw.Flush()
		index++
	}

	return -1, nil
}

func (f *txtWriter) Close() {
	f.file.Close()
}
