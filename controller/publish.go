package controller

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/RaymondCode/simple-demo/db"
	"github.com/gin-gonic/gin"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := c.PostForm("token")
	title := c.PostForm("title")
	var user db.User
	/*if _, exist := usersLoginInfo[token]; !exist {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}*/
	res := db.Mysql.Find(&user, "name = ?", token)
	if res.RowsAffected == 0 {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}
	fmt.Println("success1")
	data, err := c.FormFile("data")

	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	fmt.Println("success2")
	filename := filepath.Base(data.Filename)
	//user := usersLoginInfo[token]
	finalName := fmt.Sprintf("%d_%s", user.ID, filename)
	fmt.Println(finalName) //url:http://127.0.0.1:8080/static/1_VID_20220601_165824.mp4
	saveFile := filepath.Join("./public/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	fmt.Println("success3")

	db.Mysql.Create(&db.Video{
		Title:      title,
		UserID:     user.ID,
		PlayUrl:    "http://10.0.2.2:8080/static/" + finalName,
		CoverUrl:   "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg", //自动生成封面待完成
		CreateTime: time.Now(),
	})
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	id := c.Query("user_id")
	var video []db.Video
	var like, comm int64
	var user db.User
	var new_user User
	db.Mysql.Where("user_id = ?", id).Find(&video)
	fmt.Println(video, len(video))
	var video_list []Video
	for _, res := range video {
		id := res.ID
		db.Mysql.Model(&db.Comment{}).Where("video_id = ?", id).Count(&comm)
		db.Mysql.Model(&db.Like{}).Where("video_id = ?", id).Count(&like)
		db.Mysql.Where("id = ?", res.UserID).Find(&user)
		new_user = User{
			Id:            int64(user.ID),
			Name:          user.Name,
			FollowCount:   int64(user.FollowCount),
			FollowerCount: int64(user.FollowerCount),
			IsFollow:      false,
		}
		video_list = append(video_list, Video{
			Id:            int64(id),
			Author:        new_user,
			PlayUrl:       res.PlayUrl,
			CoverUrl:      res.CoverUrl,
			FavoriteCount: like,
			CommentCount:  comm,
			IsFavorite:    false,
		})
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: video_list,
	})
}
