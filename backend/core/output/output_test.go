package output

import "testing"

/*
----------------------------------------------------------------------------------------------------------
txtWriter
----------------------------------------------------------------------------------------------------------
*/

func TestFileWriter_WriteString(t *testing.T) {
	fw, _ := CreateTxtWriter("test.txt", true)
	defer fw.Close()

	fw.WriteString("[+] 检测成功 --> %s", "测试测试1")
	fw.WriteString("[+] 检测成功 --> %s", "测试测试1")
	fw.WriteString("[+] 检测成功 --> %s", "测试测试1")
}

func TestFileWriter_WriteStringLn(t *testing.T) {
	fw, _ := CreateTxtWriter("test.txt", true)
	defer fw.Close()

	fw.WriteStringLn("[+] 检测成功 --> %s", "测试测试2")
	fw.WriteStringLn("[+] 检测成功 --> %s", "测试测试2")
	fw.WriteStringLn("[+] 检测成功 --> %s", "测试测试2")
}

func TestFileWriter_WriteStrings(t *testing.T) {
	fw, _ := CreateTxtWriter("test.txt", true)
	defer fw.Close()

	var contents = []string{
		"[+] 检测成功 --> 测试测试3",
		"[+] 检测成功 --> 测试测试3",
		"[+] 检测成功 --> 测试测试3",
	}

	fw.WriteStrings(contents)
}

func TestFileWriter_WriteStringsLn(t *testing.T) {
	fw, _ := CreateTxtWriter("test.txt", false)
	defer fw.Close()

	var contents = []string{
		"[+] 检测成功 --> 测试测试4",
		"[+] 检测成功 --> 测试测试4",
		"[+] 检测成功 --> 测试测试4aa",
	}

	fw.WriteStringsLn(contents)
}

/*
----------------------------------------------------------------------------------------------------------
jsonWriter
----------------------------------------------------------------------------------------------------------
*/

func TestJsonWriter_Write(t *testing.T) {
	type student struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Age      int    `json:"age"`
	}

	students := []student{
		{Username: "AAA", Email: "a@qq.com", Age: 19},
		{Username: "BBB", Email: "b@qq.com", Age: 20},
		{Username: "CCC", Email: "c@qq.com", Age: 22},
	}

	jw, _ := CreateJsonWriter("test.json")
	defer jw.Close()

	jw.Write(students)
}

/*
----------------------------------------------------------------------------------------------------------
csvWriter
----------------------------------------------------------------------------------------------------------
*/

func TestCsvWriter_WriteRow(t *testing.T) {
	cw, _ := CreateCsvWriter("test.csv")
	defer cw.Close()

	cw.WriteRow([]string{"姓名", "邮箱", "年龄"})
	cw.WriteRow([]string{"张三", "zhang@qq.com", "20"})
	cw.WriteRow([]string{"李四", "lisi@gmail.com", "21"})
}

func TestCsvWriter_WriteRows(t *testing.T) {
	cw, _ := CreateCsvWriter("test.csv")
	defer cw.Close()

	cw.WriteRows([][]string{
		{"姓名", "邮箱", "年龄"},
		{"张三", "zhang@qq.com", "20"},
		{"李四", "lisi@gmail.com", "21"},
		{"李四2", "lisi2@gmail.com", "23"},
	})
}
