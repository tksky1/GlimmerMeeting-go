package controllers

import (
	"GlimmerMeeting/models"
	"GlimmerMeeting/repositories"
	"encoding/json"
	"github.com/bitly/go-simplejson"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"sort"
	"strconv"
)

func GetMeetingList(c *gin.Context) {
	token := c.Query("token")
	if TokenMap[token] == "" {
		set400(c, "权限不足！")
		return
	}

	repo := repositories.NewMeetingRepository()
	meetingRecords, err := repo.List()
	if err != nil {
		set500(c, err)
		return
	}

	meetings := make([]*models.Meeting, len(meetingRecords))
	for i, meeting := range meetingRecords {
		meetings[i] = meeting.ToMeeting()
	}

	c.JSON(200, meetings)
}

// GetMeetingByBooker 查看个人会议记录
func GetMeetingByBooker(c *gin.Context) {
	token := c.PostForm("token")
	if TokenMap[token] == "" {
		set400(c, "权限不足！")
		return
	}

	username := TokenMap[token]
	repoUser := repositories.NewUserRepository()
	user, err := repoUser.FindByUsername(username)
	repo := repositories.NewMeetingRepository()
	meetingRecords, err := repo.GetByBookerID(int(user.ID))
	meetings := make([]*models.Meeting, len(meetingRecords))
	for i, mr := range meetingRecords {
		meetings[i] = mr.ToMeeting()
	}
	if err != nil {
		set500(c, err)
		return
	}
	c.JSON(200, meetings)
}

// GetOccupiedTime 返回当天会议室被占用时段
func GetOccupiedTime(c *gin.Context) {
	token := c.PostForm("token")
	roomIDString := c.PostForm("roomid")
	date := c.PostForm("date")
	if TokenMap[token] == "" {
		set400(c, "权限不足！")
		return
	}

	repo := repositories.NewMeetingRepository()
	roomID, err := strconv.Atoi(roomIDString)
	if err != nil {
		set500(c, err)
		return
	}
	meetingRecords, err := repo.GetByDayAndRoomID(date, roomID)
	if err != nil && err != gorm.ErrRecordNotFound {
		set500(c, err)
		return
	}

	if err == gorm.ErrRecordNotFound {
		c.JSON(200, nil)
		return
	}
	var records models.MeetingRecordList = meetingRecords
	sort.Sort(records)

	var timepieces []models.Timepiece
	for i, record := range records {
		if i > 0 && record.BeginHour == records[i-1].EndHour && record.BeginMinute == records[i-1].EndMinute {
			timepieces[len(timepieces)-1].EndHour = record.EndHour
			timepieces[len(timepieces)-1].EndMinute = record.EndMinute
		} else {
			timepieces = append(timepieces, models.Timepiece{
				BeginHour:   record.BeginHour,
				BeginMinute: record.BeginMinute,
				EndHour:     record.EndHour,
				EndMinute:   record.EndMinute,
			})
		}
	}

	c.JSON(200, timepieces)

}

// PutMeeting 这里PutPost通用
func PutMeeting(c *gin.Context) {
	byteJson, err := io.ReadAll(c.Request.Body)
	if err != nil {
		set400(c, "请求参数错误！"+err.Error())
		return
	}
	bodyJson, _ := simplejson.NewJson(byteJson)
	tokenJson, success := bodyJson.CheckGet("token")
	if tokenJson == nil || !success {
		set400(c, "token无法解析："+string(byteJson))
		return
	}
	token, err := tokenJson.String()
	if err != nil {
		set400(c, "token无法解析："+string(byteJson))
		return
	}
	if TokenMap[token] == "" {
		set400(c, "权限不足！")
		return
	}

	// 解析请求体
	var meeting models.Meeting
	err = json.Unmarshal(byteJson, &meeting)
	if err != nil {
		set400(c, "请求参数错误！"+err.Error())
		return
	}

	repoUser := repositories.NewUserRepository()
	username := TokenMap[token]
	user, err := repoUser.FindByUsername(username)
	if err != nil {
		set500(c, err)
		return
	}
	meeting.BookerID = int(user.ID)
	repoRoom := repositories.NewRoomRepository()
	room, err := repoRoom.GetById(uint(meeting.RoomID))
	if err != nil {
		set500(c, err)
		return
	}
	meeting.RoomLocation = room.Location

	// 冲突检查
	repo1 := repositories.NewMeetingRepository()
	meetingRecords, err := repo1.GetByDayAndRoomID(meeting.Day, meeting.RoomID)
	if err != nil && err != gorm.ErrRecordNotFound {
		set500(c, err)
		return
	}
	var records models.MeetingRecordList = meetingRecords
	sort.Sort(records)

	var timepieces []models.Timepiece
	for i, record := range records {
		if i > 0 && record.BeginHour == records[i-1].EndHour && record.BeginMinute == records[i-1].EndMinute {
			timepieces[len(timepieces)-1].EndHour = record.EndHour
			timepieces[len(timepieces)-1].EndMinute = record.EndMinute
		} else {
			timepieces = append(timepieces, models.Timepiece{
				BeginHour:   record.BeginHour,
				BeginMinute: record.BeginMinute,
				EndHour:     record.EndHour,
				EndMinute:   record.EndMinute,
			})
		}
	}

	for _, record := range timepieces {
		if models.TimePieceConflict(meeting.Duration, record) {
			c.JSON(401, gin.H{
				"message": "时间段已有占用！",
			})
			return
		}
	}

	// 将Meeting转换为MeetingRecord
	meetingRecord := meeting.ToMeetingRecord()

	// 创建Meeting
	repo := repositories.NewMeetingRepository()
	err = repo.Update(meetingRecord)
	if err != nil {
		set500(c, err)
		return
	}

	set200(c, "添加会议成功！")
}

func DeleteMeeting(c *gin.Context) {
	token := c.Query("token")
	if TokenMap[token] != "admin" {
		set400(c, "权限不足！")
		return
	}

	// 获取请求参数
	id := c.Query("meetingid")
	repo := repositories.NewMeetingRepository()
	idInt, err := strconv.Atoi(id)
	if err != nil {
		set400(c, "id格式错误")
		return
	}
	_, err = repo.GetById(uint(idInt))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			set400(c, "会议不存在！")
		} else {
			set500(c, err)
		}
		return
	}

	// 删除Meeting
	err = repo.DeleteById(uint(idInt))
	if err != nil {
		set500(c, err)
		return
	}

	set200(c, "删除成功！")
}
