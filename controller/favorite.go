package controller

import (
	"fmt"
	"net/http"

	"github.com/RaymondCode/simple-demo/db"
	"github.com/gin-gonic/gin"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	//userId := c.Query("user_id")
	videoId := c.Query("video_id")
	action := c.Query("action_type") //1-like;2-cancle like
	var user db.User
	var video db.Video
	var like db.Like

	res := db.Mysql.Find(&user, "name = ?", token)
	if res.RowsAffected == 0 {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}
	db.Mysql.Where("id = ?", videoId).Find(&video)
	if action == "1" {
		db.Mysql.Create(&db.Like{
			UserID:  user.ID,
			VideoID: video.ID,
		})
		c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "点赞成功"})
	} else {
		db.Mysql.Where("user_id = ? AND video_id = ?", user.ID, video.ID).Find(&like)
		fmt.Println(like)
		db.Mysql.Unscoped().Where("user_id = ? AND video_id = ?", user.ID, video.ID).Delete(&like)
		//db.Mysql.Unscoped().Delete(&like)
		c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "取消点赞成功"})
	}

}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {

	id := c.Query("user_id")
	var like_video []db.Like
	var video db.Video
	var like, comm int64
	var user db.User
	var new_user User
	db.Mysql.Where("user_id = ?", id).Find(&like_video)
	fmt.Println(like_video, len(like_video))
	var video_list []Video
	for _, res := range like_video {
		id := res.VideoID
		db.Mysql.Model(&db.Comment{}).Where("video_id = ?", id).Count(&comm)
		db.Mysql.Model(&db.Like{}).Where("video_id = ?", id).Count(&like)
		db.Mysql.Where("user_id = ?", id).Find(&video)
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
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: like,
			CommentCount:  comm,
			IsFavorite:    true,
		})
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: video_list,
	})

}
