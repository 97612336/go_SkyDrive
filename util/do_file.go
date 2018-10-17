package util

import (
	"os"
	"time"
)

//判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//创建当前文件夹，如果存在则不创建
func Create_dir_path(path string) {
	is_true, err := PathExists(path)
	CheckErr(err)
	if is_true == false {
		err := os.Mkdir(path, os.ModePerm)
		CheckErr(err)
	}
}

func Get_current_time_str() string {
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	return timeStr
}
