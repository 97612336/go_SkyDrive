package util

import (
	"time"
	"math/rand"
	"crypto/md5"
	"io"
	"encoding/hex"
	"strconv"
	"feidu/util"
)

//获取一组随机数
func Get_random_arr(count int, max_num int) []int {
	//定义随机数返回数组
	var nums_arr []int
	//定义一个中间值int64数字
	var one_tmp_num int64 = 1
	//如果数组大小小于规定的长度,则执行遍历
	for len(nums_arr) < count {
		//定义随机数种子
		time_int := time.Now().Unix() + one_tmp_num
		rand.Seed(time_int)
		//生成随机数
		one_rand_num := rand.Intn(max_num)
		one_tmp_num = time_int
		nums_arr = append(nums_arr, one_rand_num)
	}
	return nums_arr
}

//得到md5字符串
func Get_md5str(str string) string {
	h := md5.New()
	io.WriteString(h, str)
	return hex.EncodeToString(h.Sum(nil))
}

// 把字符串转变为数字
func String_to_int(one_str string) int {
	one_int, err := strconv.Atoi(one_str)
	util.CheckErr(err)
	return one_int
}

//　把数字转变为字符串
func Int_to_string(one_int int) string {
	one_string := strconv.Itoa(one_int)
	return one_string
}
