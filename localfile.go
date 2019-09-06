package zlfile

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

// 获取路径文件及其子目录文件列表
func GetDirAndSubDirFileList(pathname string) ([]string, error) {
	var fileList []string
	err := filepath.Walk(pathname,
		func(path string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			if f.IsDir() {
				fileList = append(fileList, path)
				return nil
			}

			return nil
		})
	return fileList, err
}

// 获取路径下文件及文件夹名称列表
func GetDirFileList(path string) ([]string, error) {
	var fileList []string
	rd, err := ioutil.ReadDir(path)
	if err != nil {
		return fileList, err
	}
	for _, fi := range rd {
		fileList = append(fileList, fi.Name())
	}
	return fileList, nil
}

// 判断本地路径是否为目录
func IsDir(path string) (bool, error) {
	f, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return f.IsDir(), nil
}

// 将path文件权限为filePermission
func Chmod(path string, modeFile *ChmodFile) error {
	if modeFile.Mode == 0 {
		if len(modeFile.FilePermission) > 0 {
			if modeFile.FilePermission[0] != '0' {
				modeFile.FilePermission = "0" + modeFile.FilePermission
			}
			permNum, err := strconv.ParseInt(modeFile.FilePermission, 8, 10)
			if err != nil {
				return err
			}
			modeFile.Mode = os.FileMode(permNum)
		}
	}
	return os.Chmod(path, modeFile.Mode)
}
