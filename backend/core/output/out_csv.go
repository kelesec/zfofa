package output

import (
	"encoding/csv"
	"os"
)

type csvWriter struct {
	file *os.File
}

func CreateCsvWriter(filename string) (*csvWriter, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	// 解决中文乱码问题
	_, err = file.WriteString("\xEF\xBB\xBF")
	if err != nil {
		return nil, err
	}

	return &csvWriter{file: file}, nil
}

// WriteRow 写入一行数据
func (c *csvWriter) WriteRow(row []string) error {
	w := csv.NewWriter(c.file)
	err := w.Write(row)
	if err != nil {
		return err
	}

	w.Flush()
	return nil
}

// WriteRows 写入多行数据
func (c *csvWriter) WriteRows(rows [][]string) error {
	for _, row := range rows {
		err := c.WriteRow(row)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *csvWriter) Close() {
	c.file.Close()
}
