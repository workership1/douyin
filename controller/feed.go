package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/RaymondCode/simple-demo/db"
	"github.com/gin-gonic/gin"
	//"strconv"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	var video []db.Video
	var like, comm int64
	var user db.User
	var new_user User
	var fav db.Like
	var isfav bool
	db.Mysql.Limit(30).Order("create_time DESC").Find(&video)
	fmt.Println(video, len(video))
	var video_list []Video
	for _, res := range video {
		id := res.ID
		db.Mysql.Model(&db.Comment{}).Where("video_id = ?", id).Count(&comm)
		db.Mysql.Model(&db.Like{}).Where("video_id = ?", id).Count(&like)
		db.Mysql.Where("id = ?", res.UserID).Find(&user)
		db.Mysql.Where("user_id = ? AND video_id =?", user.ID, id).Find(&fav)
		if fav.UserID == 0 {
			isfav = false
		} else {
			isfav = true
		}
		fmt.Println(isfav)
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
			IsFavorite:    isfav,
		})
	}

	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: video_list,
		NextTime:  time.Now().Unix(),
	})
}
