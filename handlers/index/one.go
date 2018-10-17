package index

import (
	"go_SkyDrive/config"
	"go_SkyDrive/util"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

//具体的添加文件的操作
func Upload_file_v1(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1024 * 1024 * 10)
	if r.Method == "POST" {
		var data = make(map[string]interface{})
		var err error
		//获取验证码
		code := util.Get_argument(r, "code")
		current_code := os.Getenv("code")
		if code != current_code {
			data["code"] = 400
			data["msg"] = "Code is not right!"
			util.Return_json(w, data)
			return
		}
		// 获取当前时间天数,然后转换成字符串,重新拼接路径
		time_str := util.Get_current_time_str()
		time_str_arr := strings.Split(time_str, " ")
		day_str := time_str_arr[0]
		file_path := config.Static_Path + "upload_file/" + day_str
		// 查询目录是否存在,如果不存在就创建
		util.Create_dir_path(file_path)
		//获取表单文件
		file, header, err := r.FormFile("file")
		if err != nil {
			log.Println("接收表单文件出错")
			log.Println(err)
			return
		}
		defer file.Close()
		//创建写入到本地的文件名
		file_name := util.Get_md5str(header.Filename + time_str)
		file_name_str_arr := strings.Split(header.Filename, ".")
		ext := file_name_str_arr[len(file_name_str_arr)-1]
		new_file_name := file_name + "." + ext
		//创建文件名和后缀名组合
		file_path_and_name := file_path + "/" + new_file_name
		f, err := os.Create(file_path_and_name)
		if err != nil {
			log.Println("创建本地文件的时候出错")
			log.Println(err)
		}
		defer f.Close()
		//执行写入操作
		io.Copy(f, file)

		data["code"] = 200
		data["url"] = config.Static_href + "upload_file/" + day_str + "/" + new_file_name
		util.Return_json(w, data)
		return
	}

}
