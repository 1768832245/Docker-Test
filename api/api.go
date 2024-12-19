package api

import (
	"RR7/models"
	"RR7/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func Register(c *gin.Context) {
	nickname := c.PostForm("nickname")
	username := c.PostForm("username")
	password := c.PostForm("password")
	user := models.User{
		Nickname: nickname,
		Username: username,
		Password: password,
	}

	models.DB.Create(&user)
	c.String(http.StatusOK, "success")
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	user := models.User{
		Username: username,
		Password: password,
	}
	err := models.DB.Where("username = ? AND password = ?", user.Username, user.Password).First(&user).Error
	if err != nil {
		c.String(http.StatusBadGateway, "用户名或密码错误")
		return
	}
	token, err := utils.GenerateToken(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "token generate error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "login successful",
		"token":   token,
	})
}

func PostLog(c *gin.Context) {
	content := c.PostForm("content")
	parentID := c.DefaultPostForm("parentid", "") // 如果没有parentid参数，默认为空字符串

	var user models.User
	username, exist := c.Get("username")
	if !exist {
		c.String(http.StatusUnauthorized, "未登录")
		return
	}

	models.DB.Where("username = ?", username).First(&user)

	var parentIDUint *int
	if parentID != "" {
		id, err := strconv.Atoi(parentID)
		if err != nil {
			c.String(http.StatusBadRequest, "ParentID 参数无效")
			return
		}
		parentIDUint = &id
	}

	message := models.Message{
		UserID:   user.ID,
		Content:  content,
		ParentID: parentIDUint, // 设置 ParentID
	}
	models.DB.Create(&message)

	c.String(http.StatusOK, "success")
}

func GetLogs(c *gin.Context) {
	var messages []models.Message
	var user models.User

	models.DB.Where("parent_id IS NULL AND is_deleted = ?", false).Find(&messages)

	// 遍历根留言，查找每条留言的回复（楼中楼）
	var result []map[string]interface{}
	for _, message := range messages {

		models.DB.Where("id = ?", message.UserID).First(&user)

		replies := getReplies(message.ID)

		result = append(result, map[string]interface{}{
			"Username": user.Username,
			"Content":  message.Content,
			"Replies":  replies,
		})
	}

	c.JSON(http.StatusOK, result)
}

func getReplies(parentID int) []map[string]interface{} {
	var replies []models.Message
	var user models.User
	var replyData []map[string]interface{}

	if err := models.DB.Where("parent_id = ? AND is_deleted = ?", parentID, false).Find(&replies).Error; err != nil {
		log.Println("获取回复时出现错误:", err)
		return nil
	}
	if len(replies) == 0 {
		return nil // 如果没有回复，返回空
	}

	for _, reply := range replies {
		childReplies := getReplies(reply.ID)

		replyData = append(replyData, map[string]interface{}{
			"Username": user.Username,
			"Content":  reply.Content,
			"Replies":  childReplies,
		})
	}

	return replyData
}

func DeleteLog(c *gin.Context) {
	messageid := c.PostForm("messageid")
	username, exist := c.Get("username")

	if !exist {
		c.String(http.StatusBadRequest, "unsuccessful")
		return
	}

	var user models.User
	err := models.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		c.String(http.StatusBadRequest, "用户不存在")
		return
	}

	err = models.DB.Where("id = ? AND user_id = ?", messageid, user.ID).Delete(&models.Message{}).Error
	if err != nil {
		c.String(http.StatusBadRequest, "你无权操作")
		return
	}

	c.String(http.StatusOK, "success")
}

func LikeMessage(c *gin.Context) {
	messageid := c.PostForm("messageid")
	username, exist := c.Get("username")
	if !exist {
		c.String(http.StatusUnauthorized, "no login")
		return
	}

	var user models.User
	var message models.Message
	var like models.Like

	err := models.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		c.String(http.StatusBadRequest, "no user")
		return
	}

	err = models.DB.Where("id = ?", messageid).First(&message).Error
	if err != nil {
		c.String(http.StatusBadRequest, "no message")
		return
	}

	err = models.DB.Where("user_id = ? AND message_id = ?", user.ID, message.ID).First(&like).Error
	if err == nil {
		c.String(http.StatusBadRequest, "已经点赞过！")
		return
	}

	like = models.Like{
		UserID:    user.ID,
		MessageID: message.ID,
	}
	models.DB.Create(&like)

	models.DB.Model(&message).Update("LikesCount", message.LikesCount+1)

	c.String(http.StatusOK, "点赞成功")
}

func CancelLikeMessage(c *gin.Context) {
	messageid := c.PostForm("messageid")
	username, exist := c.Get("username")
	if !exist {
		c.String(http.StatusUnauthorized, "no login")
		return
	}

	var user models.User
	var message models.Message
	var like models.Like

	err := models.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		c.String(http.StatusBadRequest, "no user")
		return
	}

	err = models.DB.Where("id = ?", messageid).First(&message).Error
	if err != nil {
		c.String(http.StatusBadRequest, "no message")
		return
	}

	//上面保证了点赞唯一性
	err = models.DB.Where("user_id = ? AND message_id = ?", user.ID, message.ID).First(&like).Error
	//与上面区别，一个找的到一个找不到
	if err != nil {
		c.String(http.StatusBadRequest, "未曾点赞过！")
		return
	}

	models.DB.Delete(&like)

	models.DB.Model(&message).Update("LikesCount", message.LikesCount-1)

	c.String(http.StatusOK, "取消点赞成功")
}

func GetMessageLikes(c *gin.Context) {
	messageID := c.PostForm("messageid")

	var message models.Message
	err := models.DB.Where("id = ?", messageID).First(&message).Error
	if err != nil {
		c.String(http.StatusBadRequest, "no exist")
		return
	}

	// 返回点赞数
	c.JSON(http.StatusOK, gin.H{
		"messageid":  message.ID,
		"LikesCount": message.LikesCount,
	})
}
