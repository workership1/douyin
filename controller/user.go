package controller

import (
	//"strings"
	"fmt"
	"net/http"
	"time"

	"github.com/RaymondCode/simple-demo/db"
	"github.com/gin-gonic/gin"
	//"sync/atomic"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

//var userIdSequence = int64(1)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

// code 0：成功， 1：失败， 2：用户或密码为空

func Register(c *gin.Context) {

	var user db.User
	username := c.Query("username")
	password := c.Query("password")

	fmt.Println("success1")

	//需要判断该用户名是否被占用
	db.Mysql.Where("Name = ?", username).Find(&user)
	if user.ID != 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "该用户名已被占用",
		})
		return
	}
	fmt.Println("success2")
	//将用户信息插入数据库
	db.Mysql.Create(&db.User{
		Name:          username,
		Password:      password,
		FollowCount:   0,
		FollowerCount: 0,
		RegisterTime:  time.Now(),
	})

	fmt.Println("success3")

	/*c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
	})*/

	db.Mysql.Where("Name = ?", username).Find(&user)
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0},
		UserId:   int64(user.ID),
		Token:    username,
	})
	return

}

func Login(c *gin.Context) {

	username := c.Query("username")
	password := c.Query("password")

	//先查看是否存在该用户
	var user db.User
	//db.Mysql.AutoMigrate(&user)
	res := db.Mysql.Find(&user, "Name = ? AND Password = ?", username, password)
	// select * from user where
	if res.RowsAffected == 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
		return
	}
	token := username

	db.Mysql.Where("Name = ?", username).Find(&user)
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0},
		UserId:   int64(user.ID),
		Token:    token,
	})
}

func UserInfo(c *gin.Context) {
	username := c.Query("user_id")
	var user db.User
	fmt.Println(username)
	res := db.Mysql.Find(&user, "Id = ?", username)
	fmt.Println(res)
	if res.RowsAffected == 0 {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "——User doesn't exist"},
		})
		return
	}
	db.Mysql.Where("Name = ?", username).Find(&user)
	newuser := User{
		Id:            int64(user.ID),
		Name:          user.Name,
		FollowCount:   int64(user.FollowCount),
		FollowerCount: int64(user.FollowerCount),
		IsFollow:      false,
	}
	c.JSON(http.StatusOK, UserResponse{
		Response: Response{StatusCode: 0},
		User:     newuser,
	})
	return
}
