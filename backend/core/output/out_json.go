package output

import (
	"encoding/json"
	"os"
)

type jsonWriter struct {
	file *os.File
}

func CreateJsonWriter(filename string) (*jsonWriter, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	return &jsonWriter{file: file}, nil
}

// Write 写入JSON文件
func (j *jsonWriter) Write(v any) (int, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return 0, err
	}

	return j.file.Write(bytes)
}

func (j *jsonWriter) Close() {
	j.file.Close()
}
