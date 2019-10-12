package utils

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"path"
)

/*
	文件工具类
*/

/* 判断文件路径是否存在 */
func IsExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

/* 判断指定路径是否是文件 */
func IsFile(f string) bool {
	fi, e := os.Stat(f)
	if e != nil {
		return false
	}
	return !fi.IsDir()
}

/* 删除文件或文件夹 */
func RemoveAll(filePath string) bool {
	if err := os.RemoveAll(filePath); err != nil {
		return false
	}
	return true
}

/* 创建指定目录文件 */
func MakeFile(filePath string, filename string) (*os.File, error) {
	if !IsExist(filePath) {
		err := os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}
	return os.Create(path.Join(filePath, filename))
}

/* 文件拷贝 */
func CopyFile(srcName, destName string) (bool, error) {
	srcFile, err := os.Open(srcName)
	if err != nil {
		return false, err
	}
	defer srcFile.Close()
	destFile, err := os.OpenFile(destName, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return false, err
	}
	defer destFile.Close()
	buf := make([]byte, 2048)
	_, er := io.CopyBuffer(destFile, srcFile, buf)
	if er != nil {
		return false, er
	}
	return true, nil
}

/* 读取文件内容 */
func ReadFile(filePath string) ([]byte, error) {
	return ioutil.ReadFile(filePath)
}

/* 内容写入目标文件 */
func WriteFile(b []byte, destPath string) (bool, error) {
	file, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return false, err
	}
	defer file.Close()
	bw := bufio.NewWriter(file)
	_, err = bw.Write(b)
	if err != nil {
		return false, err
	}
	err = bw.Flush()
	if err != nil {
		return false, err
	}
	return true, nil
}
