package controller

import (
	"net/http"
	"time"

	"github.com/RaymondCode/simple-demo/db"
	"github.com/gin-gonic/gin"
)

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment Comment `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	token := c.Query("token")
	actionType := c.Query("action_type")
	video_id := c.Query(("video_id"))
	comment_id := c.Query("comment_id")
	var user db.User
	var video db.Video
	var comment db.Comment
	db.Mysql.Where("name = ?", token).Find(&user)
	db.Mysql.Where("id = ?", video_id).Find(&video)

	if actionType == "1" {
		text := c.Query("comment_text")
		nowtime := time.Now()
		timeStr := nowtime.Format("2006-01-02 15:04:05")
		finaltimeStr := timeStr[5:10]

		cur_comment := db.Comment{
			UserID:     user.ID,
			VideoID:    video.ID,
			Content:    text,
			CreateTime: nowtime,
			Delete:     false,
		}

		db.Mysql.Create(&cur_comment)

		cur_user := User{
			Id:            int64(user.ID),
			Name:          user.Name,
			FollowCount:   int64(user.FollowCount),
			FollowerCount: int64(user.FollowerCount),
			IsFollow:      false,
		}

		c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 0},
			Comment: Comment{
				Id:         int64(cur_comment.ID),
				User:       cur_user,
				Content:    text,
				CreateDate: finaltimeStr,
			}})
		return
	} else {
		db.Mysql.Where("id = ?", comment_id).Find(&comment)
		db.Mysql.Model(&comment).Update("delete", 1)
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	video_id := c.Query("video_id")
	var comments []db.Comment
	var comment_list []Comment
	var user db.User

	db.Mysql.Where("video_id = ?", video_id).Find(&comments)

	for _, res := range comments {
		db.Mysql.Where("id = ?", res.UserID).Find(&user)
		newuser := User{
			Id:            int64(user.ID),
			Name:          user.Name,
			FollowCount:   int64(user.FollowCount),
			FollowerCount: int64(user.FollowerCount),
			IsFollow:      false,
		}
		comment_list = append(comment_list, Comment{
			Id:         int64(res.ID),
			User:       newuser,
			Content:    res.Content,
			CreateDate: res.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: comment_list,
	})
}
