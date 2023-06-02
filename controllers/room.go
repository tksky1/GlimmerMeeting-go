package controllers

import (
	"GlimmerMeeting/models"
	"GlimmerMeeting/repositories"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

// GetRoomList 允许所有用户使用
func GetRoomList(c *gin.Context) {
	token := c.Query("token")
	if TokenMap[token] == "" {
		set400(c, "权限不足！"+token)
		return
	}

	repo := repositories.NewRoomRepository()
	rooms, err := repo.List()
	if err != nil {
		set500(c, err)
		return
	}
	c.JSON(200, rooms)
}

func PostRoom(c *gin.Context) {
	token := c.PostForm("token")
	if TokenMap[token] != "admin" {
		set400(c, "权限不足！")
		return
	}
	name := c.PostForm("name")
	info := c.PostForm("info")
	location := c.PostForm("location")
	if name == "" || location == "" || info == "" {
		set400(c, "必填字段不能为空")
		return
	}
	repo := repositories.NewRoomRepository()
	err := repo.Create(&models.Room{
		Name:     name,
		Info:     info,
		Location: location,
	})
	if err != nil {
		set500(c, err)
		return
	}
	set200(c, "添加会议室"+name+"成功！")
}

func PutRoom(c *gin.Context) {
	token := c.PostForm("token")
	if TokenMap[token] != "admin" {
		set400(c, "权限不足！")
		return
	}

	id := c.PostForm("id")
	repo := repositories.NewRoomRepository()
	idInt, err := strconv.Atoi(id)
	if err != nil {
		set400(c, "id格式错误")
		return
	}
	room, err := repo.GetById(uint(idInt))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			set400(c, "会议室不存在")
		} else {
			set500(c, err)
		}
		return
	}

	room.Name = c.PostForm("name")
	room.Info = c.PostForm("info")
	room.Location = c.PostForm("location")
	if room.Name == "" || room.Location == "" || room.Info == "" {
		set400(c, "必填字段不能为空")
		return
	}
	room.ID = uint(idInt)
	err = repo.Update(room)
	if err != nil {
		set500(c, err)
		return
	}
	set200(c, "修改成功！")
}

func DeleteRoom(c *gin.Context) {
	token := c.Query("token")
	if TokenMap[token] != "admin" {
		set400(c, "权限不足！")
		return
	}

	id := c.Param("id")
	repo := repositories.NewRoomRepository()
	idInt, err := strconv.Atoi(id)
	if err != nil {
		set400(c, "id格式错误")
		return
	}
	_, err = repo.GetById(uint(idInt))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			set400(c, "房间不存在！")
		} else {
			set500(c, err)
		}
		return
	}
	err = repo.DeleteById(uint(idInt))
	if err != nil {
		set500(c, err)
		return
	}
	set200(c, "删除成功！")
}
