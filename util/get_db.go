package util

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"go_SkyDrive/models"
	"io/ioutil"
	"log"
	"os/user"
	"strings"
	"time"
)

var DB *sql.DB

//得到家目录的路径
//得到家目录的路径
//得到家目录的路径
func Get_home_path() string {
	current_user, err := user.Current()
	if nil != err {
		fmt.Println("get user current dir err:", current_user.HomeDir)
		return ""
	}
	user_home := current_user.HomeDir
	return user_home
}

// 获取验证码code的值
func Get_code() string {
	user_home := Get_home_path()
	config_file := user_home + "/conf/code_conf"
	data, _ := ioutil.ReadFile(config_file)
	//将读取到的文件转化为字符串
	str_data := string(data)
	fmt.Println(str_data)
	fmt.Println(strings.Count(str_data, "") - 1)
	return str_data
}

//获取mysql对象
//获取mysql对象
//获取mysql对象
func Get_sql_db() *sql.DB {
	sqlconf := Get_conf_info()
	//打开数据库
	db, err := sql.Open("mysql",
		sqlconf.SqlUser+":"+sqlconf.SqlPassword+
			"@tcp("+sqlconf.SqlHost+":"+sqlconf.SqlPort+")/bigbiy_web?charset=utf8")
	if err != nil {
		log.Println("打开数据库出错")
	}
	//设置最大连接数
	db.SetMaxOpenConns(100)
	//设置连接池最大数
	db.SetMaxIdleConns(50)
	//设置每个链接的存活的时长
	db.SetConnMaxLifetime(time.Second * 50)
	return db
}

//获取mysql配置文件信息
func Get_conf_info() models.SqlConf {
	user_home := Get_home_path()
	config_file := user_home + "/conf/sqlconf"
	//	读取数据库配置文件
	data, _ := ioutil.ReadFile(config_file)
	//转化为字符串格式
	str_data := string(data)
	//实例化数据库配置类型
	var sqlconf models.SqlConf
	//得到json字符串数据
	var sql_json = []byte(str_data)
	//把json数据赋值给实例化的数据配置对象
	json.Unmarshal(sql_json, &sqlconf)
	return sqlconf
}

//获取上传文件的验证文件
func Get_img_account() models.Upload_account {
	user_home := Get_home_path()
	config_file := user_home + "/conf/upload_account"
	data, _ := ioutil.ReadFile(config_file)
	//将读取到的文件转化为字符串
	str_data := string(data)
	var account models.Upload_account
	//将读取到的配置文件赋值给类型
	var account_json = []byte(str_data)
	json.Unmarshal(account_json, &account)
	return account
}

//获取redis操作对象
//获取redis操作对象
//获取redis操作对象
func Get_redis_conf() models.Redis_conf {
	//获取reids相关的配置文件
	home_path := Get_home_path()
	config_file := home_path + "/conf/redis_conf"
	data, _ := ioutil.ReadFile(config_file)
	//将读取到的数据赋值给对象
	var redis_conf models.Redis_conf
	json.Unmarshal(data, &redis_conf)
	//通过对象信息,链接到redis数据库
	return redis_conf
}

//设置redis字符串的值
func Set_redis(key string, value string, times ...string) {
	//获取配置信息,连接到redis
	redis_conf := Get_redis_conf()
	//连接到redis数据库
	red, err := redis.Dial("tcp", string(redis_conf.Ip_addr)+":"+string(redis_conf.Port))
	CheckErr(err, "链接redis数据库出错:")
	defer red.Close()
	//判断是否有时间参数,如果有的话就设置过期时间,没有的话就不设置
	if len(times) < 1 {
		_, err := red.Do("set", key, value)
		CheckErr(err)
	} else {
		_, err := red.Do("set", key, value, "EX", times[0])
		CheckErr(err)
	}
}

//获取redis的值
func Get_redis(key string) string {
	//获取配置信息,连接到redis
	redis_conf := Get_redis_conf()
	//连接到redis数据库
	red, err := redis.Dial("tcp", string(redis_conf.Ip_addr)+":"+string(redis_conf.Port))
	CheckErr(err, "链接redis数据库出错:")
	defer red.Close()
	res, err := redis.String(red.Do("get", key))
	CheckErr(err)
	return res
}

//banner_novel的配置
func Get_banner_novel_id() []int {
	home_path := Get_home_path()
	config_path := home_path + "/conf/banner_novel"
	data, err := ioutil.ReadFile(config_path)
	CheckErr(err)
	var banner_id_list []int
	json.Unmarshal(data, &banner_id_list)
	return banner_id_list
}
