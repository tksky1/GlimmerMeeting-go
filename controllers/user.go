package controllers

import (
	"GlimmerMeeting/models"
	"GlimmerMeeting/repositories"
	"crypto/sha256"
	"encoding/hex"
	json "github.com/bitly/go-simplejson"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"io"
	"os"
	"strconv"
)

var TokenMap map[string]string

func UserLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	hash := sha256.Sum256([]byte(password))

	repo := repositories.NewUserRepository()
	user, err := repo.FindByUsername(username)
	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
	} else {
		if user == nil || user.Password != hex.EncodeToString(hash[:]) {
			c.JSON(400, gin.H{
				"message": "用户名或密码错误！",
			})
		} else {
			tokenUUID, _ := uuid.NewRandom()
			token := tokenUUID.String()
			TokenMap[token] = username
			c.JSON(200, gin.H{
				"message": "登录成功！",
				"token":   token,
			})
		}
	}
}

func PostUser(c *gin.Context) {

	token := c.Query("token")
	if TokenMap[token] != "admin" {
		set400(c, "权限不足！")
		return
	}
	facePicsJSON, _ := io.ReadAll(c.Request.Body)
	picsJson, err := json.NewJson(facePicsJSON)
	if err != nil {
		set400(c, "请求格式不正确！")
	}
	username := picsJson.Get("username").MustString()
	password := picsJson.Get("password").MustString()
	realName := picsJson.Get("real_name").MustString()
	repo := repositories.NewUserRepository()
	user, err := repo.FindByUsername(username)
	if err != nil && err != gorm.ErrRecordNotFound {
		set500(c, err)
		return
	}
	if username == "" || (user != nil && user.Username == username) {
		set400(c, "用户已存在")
		return
	}
	if realName == "" {
		set400(c, "")
		return
	}
	codedPsw := sha256.Sum256([]byte(password))
	err = repo.Create(&models.User{
		Username: username,
		RealName: realName,
		Password: hex.EncodeToString(codedPsw[:]),
	})
	if err != nil {
		set500(c, err)
		return
	}

	//处理图片
	err = os.Mkdir("./static/"+username, 0777)
	if err != nil {
		log.Error(err)
	}

	picsJson = picsJson.Get("facepic")
	picsNum := len(picsJson.MustArray())
	for i := 0; i < picsNum; i++ {
		picJson := picsJson.GetIndex(i)
		tokenUUID, _ := uuid.NewRandom()
		filename := tokenUUID.String()
		dataString := picJson.Get("filedata").MustString()
		err := os.WriteFile("./static/"+username+"/"+filename, []byte(dataString), 0777)
		if err != nil {
			set500(c, err)
			return
		}
	}
	set200(c, "用户创建成功")
}

func PutUser(c *gin.Context) {
	token := c.PostForm("token")
	if TokenMap[token] != "admin" {
		set400(c, "权限不足！")
		return
	}

	id := c.PostForm("id")
	repo := repositories.NewUserRepository()
	idInt, err := strconv.Atoi(id)
	if err != nil {
		set400(c, "id格式错误")
		return
	}
	user, err := repo.FindByID(uint(idInt))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			set400(c, "用户不存在")
		} else {
			set500(c, err)
		}
		return
	}

	username := c.PostForm("username")
	realName := c.PostForm("real_name")
	user.ID = uint(idInt)
	user.Username = username
	user.RealName = realName
	err = repo.Update(user)
	if err != nil {
		set500(c, err)
		return
	}
	set200(c, "修改成功！")
}

func GetUserList(c *gin.Context) {
	token := c.Query("token")
	if TokenMap[token] != "admin" {
		set400(c, "权限不足！")
		return
	}

	repo := repositories.NewUserRepository()
	users, err := repo.List()
	if err != nil {
		set500(c, err)
		return
	}
	c.JSON(200, gin.H{
		"users": users,
	})
}

func DeleteUser(c *gin.Context) {
	token := c.Query("token")
	if TokenMap[token] != "admin" {
		set400(c, "权限不足！")
		return
	}

	id := c.Param("id")
	repo := repositories.NewUserRepository()
	idInt, err := strconv.Atoi(id)
	if err != nil {
		set400(c, "id格式错误")
		return
	}
	theUser, err := repo.FindByID(uint(idInt))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			set400(c, "用户不存在！")
		} else {
			set500(c, err)
		}
		return
	}
	err = repo.Delete(theUser)
	if err != nil {
		set500(c, err)
		return
	}
	set200(c, "删除成功！")
}
