/*
	文件操作相关
*/
package mutil

import (
	"io/ioutil"
	"os"
	"strings"
)

// pathExists 文件存在性验证，存在true
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// CreateDir 创建文件目录
func CreateDir(_dir string) (err error) {
	err = os.MkdirAll(_dir, os.ModePerm)
	if err != nil {
		return
	}

	return
}

// ReadFileByIoUtil 使用IoUtil读取文件内容
func ReadFileByIoUtil(filePath string) (contents []byte, err error) {
	contents, err = ioutil.ReadFile(filePath)
	if err != nil {
		return
	}

	return
}

// WriteWithIoUtil 使用IoUtil写入文件内容
func WriteWithIoUtil(name, content string) error {
	data := []byte(content)
	if err := ioutil.WriteFile(name, data, 0644); err != nil {
		return err
	}
	return nil
}

// GetAllFile 返回目录下的全部文件
func GetAllFile(path string, justName bool) (fileName []string, err error) {

	path = strings.TrimRight(path, "/") + "/"
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}
	if len(files) > 0 {
		for _, v := range files {

			//只要文件名,不要路径
			if v.IsDir() {
				dir := path + v.Name()
				fileName, err = GetAllFile(dir, justName)
				if err != nil {
					return
				}
			} else {
				fn := ""

				//不要路径
				if justName {
					fn = v.Name()
				} else {
					fn = path + v.Name()
				}

				fileName = append(fileName, fn)
			}
		}
	}
	return
}
