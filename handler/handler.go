package handler

import (
	. "Demo_order/db"
	. "Demo_order/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

//func HtmlErr(err *gin.Error){
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{
//			"message": err.Error(),
//		})
//		return
//	}
//}

//创建记录
func OnCreatHandler(c *gin.Context){
   //处理数据
	order_no := c.Query("order_no")
	user_name := c.Query("user_name")
	amount,_ := strconv.ParseFloat(c.Query("amount"), 64)
	status := c.Query("status")

	order := Demo_order{Order_no : order_no ,User_name : user_name, Amount : amount, Status:status}
	err := GetDb().Model(&Demo_order{}).Create(&order).Error
	CheckErr(err)

	c.JSON(200, gin.H{
		//"err"   : err.Error(),
		"order_no":  order_no,
		"user_name": user_name,
		"amount":    amount,
		"status":status,
	})
}

//查询
func OnSearchHandler(c *gin.Context){
  //按照ID进行查询
	//根据ID进行删除数据
	id := c.Query("id")
	order := Demo_order{}
	err := GetDb().Model(&Demo_order{}).Where("id=?",id).Find(&order).Error
	CheckErr(err)
	c.JSON(200, gin.H{
		//"err"   : err.Error(),
		"id":id,
		"order_no":  order.Order_no,
		"user_name": order.User_name,
		"amount":    order.Amount,
		"status": order.Status,
		"file_url": order.File_url,
	})
}

//列表查询 （根据USERNAME模糊查询并且按照时间排序 按时间升序）
func OnSearchListHandler(c *gin.Context){
	//按照ID进行查询
	//根据ID进行删除数据
	name := c.Query("name")
	orders := []Demo_order{}
	//find := "user_name LIKE" +"%"+name+"%"
	err := GetDb().Model(&Demo_order{}).Where("user_name LIKE ?","%"+name+"%").Find(&orders).Error
	CheckErr(err)

	//判断是否为空
	if len(orders) == 0{
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "not find",
		})
		return
	}
	var result string
    var result2 map[string]string
	result2 = make(map[string]string)
	//格式化输出
	for k,v := range orders{
		result +=fmt.Sprint("%d:%s",k,v)
		result2 [fmt.Sprintf("%d",k)] = fmt.Sprintf("%s",v)

		//b, err := json.Marshal(v)
		//if err != nil {
		//	fmt.Println("Umarshal failed:", err)
		//	return
		//}

	}
	c.JSON(200,result2)
}


//更新数据
func OnUpDataHandler(c *gin.Context){
  //根据ID进行更新数据
	id := c.Query("id")
	order_no := c.Query("order_no")
	user_name := c.Query("user_name")
	amount,_ := strconv.ParseFloat(c.Query("amount"), 64)
	status := c.Query("status")
	order := Demo_order{Order_no : order_no ,User_name : user_name, Amount : amount, Status:status}

	err := GetDb().Model(&Demo_order{}).Where("id=?",id).Update(&order).Error
	CheckErr(err)

	c.JSON(200, gin.H{
		//"err"   : err.Error(),
		"id":id,
		"order_no":  order_no,
		"user_name": user_name,
		"amount":    amount,
		"status":status,
	})
}

func OnDeleteHandler(c *gin.Context){
	//根据ID进行删除数据
	id := c.Query("id")
	//此处Unscoped是物理删除
	err := GetDb().Model(&Demo_order{}).Where("id=?",id).Unscoped().Delete(&Demo_order{}).Error
	CheckErr(err)

	c.JSON(200, gin.H{
		//"err"   : err.Error(),
		"id":id,
	})
}

//下载数据库内数据生成excel表
func OnDownloadExcelHandler(c *gin.Context){
	//表头
	titleList := []string{"id", "created_at", "updated_at", "deleted_at", "order_no", "user_name", "amount", "status", "file_url"}
	//获取所有数据
	orders := []Demo_order{};
	err := GetDb().Model(&Demo_order{}).Find(&orders).Error
	CheckErr(err)

	//生成一个新文件
	file := xlsx.NewFile()
	//添加sheet页
	sheet,_ := file.AddSheet("demo_order")
	//插入表头
	titleRow := sheet.AddRow()
	for _, v := range titleList {
		cell := titleRow.AddCell()
		cell.Value = v
		//表头字体颜色
		cell.GetStyle().Font.Color = "00FF0000"
		//居中显示
		cell.GetStyle().Alignment.Horizontal = "center"
		cell.GetStyle().Alignment.Vertical = "center"
	}
	// 插入内容
	for _, v := range orders {
		row := sheet.AddRow()

		cell := row.AddCell()
		cell.Value = fmt.Sprintf("%d",v.ID)
		cell = row.AddCell()
		cell.Value = v.CreatedAt.String();
		cell = row.AddCell()
		cell.Value = v.UpdatedAt.String();
		cell = row.AddCell()
		cell.Value = v.DeletedAt.Time.String();
		row.WriteStruct(&v, -1)
	}

	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	disposition := fmt.Sprintf("attachment; filename=\"%s-%s.xlsx\"", "odmo_order", time.Now().Format("2006-01-02 15:04:05"))
	c.Writer.Header().Set("Content-Disposition", disposition)
	_ = file.Write(c.Writer)

}



func OnUpDateFileHandler(c *gin.Context){
	id := c.Query("id")
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	//检查ID是否存在
	order := Demo_order{}
	err = GetDb().Model(&Demo_order{}).Where("id=?",id).Find(&order).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	//此处需要开启事务
	tx := GetDb().Begin()
	//先保存地址，如果上传失败就回滚
	dst := fmt.Sprintf("./Data/%s", file.Filename)
	//保存地址
	order.File_url = dst
	err = tx.Model(&Demo_order{}).Where("id=?",id).Update(&order).Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	log.Println(file.Filename)

	// 上传文件到指定的目录
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		//此处需要先进行数据库回滚
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),

		})
		return
	}else{
		//数据上传成功则提交数据库
		tx.Commit()
		c.JSON(http.StatusOK, gin.H{
			"iD" : order.ID,
			"dst": dst,
			"message": fmt.Sprintf("'%s' uploaded!", file.Filename),
		})
	}
}

func OnDownloadFileHandler(c *gin.Context){
	//获取ID验证是否存在
	id := c.Query("id")
	order := Demo_order{}
	err := GetDb().Model(&Demo_order{}).Where("id=?",id).Find(&order).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	//获取URL检验文件是否存在
	if _, err := os.Stat(order.File_url); os.IsNotExist(err) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"file not exist" : true,
		})
		return
	}

	//下载文件
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", order.File_url))//fmt.Sprintf("attachment; filename=%s", filename)对下载的文件重命名
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.File(order.File_url)
}

//判断文件是否存在  存在返回 true 不存在返回false
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}