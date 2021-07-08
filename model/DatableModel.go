package model

import "gorm.io/gorm"

type Demo_order struct{
	gorm.Model   //内嵌4个字段
	//Id int `gorm:"column:id"`
	Order_no string `gorm:"column:order_no""type:varchar(120)"`//订单号
	User_name string `gorm:"column:user_name""type:varchar(120)"`//用户名
	Amount float64 `gorm:"column:amount""default:0""type:float"`//金额
	Status string `gorm:"column:status""default:null""type:varchar(120)"`//状态
	File_url string `gorm:"column:file_url""index:addr""default:null""type:varchar(120)"`//文件地址

}