package db

import (
	. "Demo_order/model"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)
var Db *gorm.DB
// 错误检测函数
func CheckErr(err error) {
	if err != nil {
		fmt.Println("err:"+err.Error())
		//panic(err)
	}
}

func DbInit() (db *gorm.DB ,err error){
	//配置MySQL连接参数
	username := "root"  //账号
	password := "123" //密码
	host := "127.0.0.1" //数据库地址，可以是Ip或者域名
	port := 3306 //数据库端口
	Dbname := "test_db" //数据库名

	//通过前面的数据库参数，拼接MYSQL DSN， 其实就是数据库连接串（数据源名称）
	//MYSQL dsn格式： {username}:{password}@tcp({host}:{port})/{Dbname}?charset=utf8&parseTime=True&loc=Local
	//类似{username}使用花括号包着的名字都是需要替换的参数
	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, Dbname)
	//db,err := gorm.Open("mysql","root:123@(127.0.0.1)/test?charset=utf8mb4&loc=Local")
	Db, err = gorm.Open("mysql" , dsn)
	CheckErr(err)
	//禁用复数结构
	Db.SingularTable(true)
	//如果表不存在则创建
	Db.AutoMigrate(&Demo_order{})
	return Db,err
}

func GetDb() (db *gorm.DB){
	if Db == nil {
		Db,_ =DbInit();
	}
	return Db;
}
//创建
func CreatMode(d *Demo_order){
	err := GetDb().Model(&Demo_order{}).Create(d).Error
	CheckErr(err)
}

//查询


//修改

//删除
