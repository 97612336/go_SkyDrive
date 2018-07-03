package handler

import (
	"net/http"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"html/template"
	"fmt"
	"strconv"
	"net/url"
	"github.com/garyburd/redigo/redis"
	"encoding/json"
	"crypto/md5"
	"encoding/hex"
	"go_SkyDrive/models"
)
var DB *sql.DB
var User models.Users
var UserInfo models.UserInfo
var Info models.Info

func md5Str(password string) string {
	h := md5.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}

func IndexHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == "GET"{
		t, _ := template.ParseFiles("static/html/index.html")
		t.Execute(w,nil)
	}
}
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("static/html/login.html")
		log.Println(t.Execute(w, nil))
		return
	}
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := md5Str(r.Form.Get("password"))

		query_users_sql := `select * from users where name = ?`
		if stmt, err := DB.Prepare(query_users_sql); err != nil {
			defer stmt.Close()
			fmt.Println("DB prepare users error:", err)
			fmt.Fprintf(w, `{"code":500, "msg":"服务器异常!"}`)
		}else{
			defer stmt.Close()
			rows, err := stmt.Query(username)
			fmt.Println(rows)
			defer rows.Close()
			for rows.Next() {
				err = rows.Scan(&User.Id, &User.Name, &User.Password )
				if err != nil {
					fmt.Println(err.Error())
					continue
				}
				fmt.Println(222, User.Password)
				if User.Password == password {
					t, _ := template.
						ParseFiles("static/html/loginsuccess.html")
					query_sql := `select infoid from user_info where userid=? `
					if stmt, err = DB.Prepare(query_sql);err!=nil{
						fmt.Println("query error!",err)
					}
					rows, _ = stmt.Query(User.Id)
					defer rows.Close()
					for rows.Next(){
						err = rows.Scan(&UserInfo.InfoId)
						if err != nil{
							fmt.Println(err.Error())
							continue
						}
					}
					t.Execute(w,UserInfo.InfoId)

				}else{
					fmt.Println(444)
					t, _ := template.
						ParseFiles("static/html/loginerr.html")
					log.Println(t.Execute(w, nil))
					return
					}
				}
			}
		}
	}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("static/html/register.html")
		t.Execute(w, nil)
	}
	if r.Method == "POST" {
		var info_id int64
		var user_id int64
		r.ParseForm()
		username := r.Form.Get("username")
		password := md5Str(r.Form.Get("password"))
		id_card := r.Form.Get("idcard")
		age, err := strconv.Atoi(r.Form.Get("age"))
		if err != nil {
			fmt.Println("age error!")
		}
		charsex := r.Form.Get("sex")
		var sex int
		if charsex == "男" {
			sex = 1
		} else if charsex == "女" {
			sex = 0
		} else {
			fmt.Println("sex error")
			return
		}
		address := r.Form.Get("address")
		phone, err := strconv.Atoi(r.Form.Get("phone"))
		//fmt.Println(id_card)
		//fmt.Println(age)
		//fmt.Println(sex)
		//fmt.Println(address)
		//fmt.Println(phone)
		if err != nil {
			fmt.Println("phone error!")
		}
		insert_sql := `insert into info(idcard, age, sex, address, phone)
 						values(?,?,?,?,?)`

		stmt, err := DB.Prepare(insert_sql)
		if err != nil {
			defer stmt.Close()
			fmt.Println("DB insert into userinfo error:", err)
		} else {
			defer stmt.Close()
			res, err := stmt.Exec(id_card, age, sex, address, phone)
			if err != nil {
				fmt.Println(err)
				fmt.Println("db insert exec error!")
			}
			affect, _ := res.RowsAffected()
			fmt.Println(affect)
			info_id, _ = res.LastInsertId()
			fmt.Println("infor_id", info_id)
		}
		insert_sql = `insert into users(name,password) values(?,?)`
		stmt, err = DB.Prepare(insert_sql)
		if err != nil {
			defer stmt.Close()
			fmt.Println("DB insert into userinfo error:", err)
		} else {
			defer stmt.Close()
			res, err := stmt.Exec(username, password)
			if err != nil {
				fmt.Println(err)
				fmt.Println("db user insert exec error!")
			}
			user_id, _ = res.LastInsertId()
			fmt.Println("uid", user_id)
		}
		insert_sql = `insert into user_info(userid,infoid) values(?,?)`
		if stmt, err = DB.Prepare(insert_sql);err != nil{
			fmt.Println("DB prepare error!", err)
			return
		}
		defer stmt.Close()
		res, _ := stmt.Exec(user_id,info_id)
		fmt.Println(res.RowsAffected())
		t, _ := template.ParseFiles("static/html/index.html")
		t.Execute(w, nil)
	}
}
func InfoHandler(w http.ResponseWriter, r *http.Request){
	c,err := redis.Dial("tcp", "localhost:6379")
	if err != nil{
		log.Fatal("connect to redis failed!",err)
	}
	defer c.Close()
	if r.Method == "GET"{
		queryForm, err := url.ParseQuery(r.URL.RawQuery)
		if err == nil && len(queryForm["infoid"]) > 0 {
			infoid := queryForm["infoid"][0]
			fmt.Println("infoid:",infoid)
			fmt.Printf("%T",infoid)

			key := "userinfo" + infoid
			is_key_exists, err := redis.Bool(c.Do("exists",key))
			if err!=nil {
				log.Fatal("error", err)
			}
			if is_key_exists{
				valueGet, err:= redis.Bytes(c.Do("get",key))
				if err!=nil{
					log.Fatal(err)
				}
				err = json.Unmarshal(valueGet,&Info)
				if err !=nil{
					log.Fatal(err)
				}
				fmt.Println("11111111", Info)
			}else{
				query_sql := `select * from info where id=?`
				stmt, err := DB.Prepare(query_sql)
				if err != nil{
					log.Fatal("query error",err)
				}else{
					rows,err:= stmt.Query(infoid)
					defer rows.Close()
					if err != nil{
						fmt.Println("exec query error",err)
					}
					fmt.Println("rows:", rows)
					for rows.Next(){
						fmt.Println("here!")
						err = rows.Scan(&Info.Id, &Info.IdCard,&Info.Age,
							&Info.Sex,&Info.Address,&Info.Phone)
						if err!= nil{
							fmt.Println("rows error",err)
						}
					}
					fmt.Println("idcard:",Info.IdCard)
					fmt.Println("age:",Info.Age)
					fmt.Println("sex:",Info.Sex)
					fmt.Println("address:",Info.Address)
					fmt.Println("phone:",Info.Phone)

					fmt.Println(err)
					if err !=nil {
						fmt.Println(err)
					}
					key := "userinfo" + strconv.Itoa(Info.Id)
					value, err := json.Marshal(Info)
					fmt.Println(string(value))
					if err!=nil{
						fmt.Println("json Marshal error:",err)
					}
					n, err := c.Do("setnx",key,value)
					if err != nil{
						fmt.Println(err)
					}
					if n == int64(1){
						fmt.Println("success!")
					}
				}
			}
			t, _ := template.ParseFiles("static/html/userinfo.html")
			t.Execute(w, Info)
		}
	}
}