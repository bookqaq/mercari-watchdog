package webapi

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"bookq.xyz/mercari-watchdog/datatype/analysisdata"
	"bookq.xyz/mercari-watchdog/datatype/analysistask"
	"bookq.xyz/mercari-watchdog/datatype/fetchdata"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type fetchedSettingsReply struct {
	expire   int64
	settings fetchdata.FetchedSettings
	override fetchdata.FetchOverride
}

var config = struct {
	Settings fetchdata.FetchedSettings
	Expire   int64
}{
	Settings: fetchdata.FetchedSettings{
		Interval: []fetchdata.Interval{
			{Time: 300, Text: "5分钟"},
			{Time: 600, Text: "10分钟"},
			{Time: 3600, Text: "1小时"},
		},
		PageRange: [2]int{1, 5},
	},
	Expire: 600,
}

func getAllRouters(router *gin.Engine) {
	tasks := router.Group("/task")
	{
		tasks.POST("/fetch", postTaskAddFetch)
		tasks.POST("/submit", postTaskAddSubmit)
	}
}

func postTaskAddFetch(c *gin.Context) {
	auth := c.PostForm("auth")
	if auth == "" {
		c.JSON(http.StatusOK, genericPostReply{
			Status:  "failed",
			Message: "没有用户数据，请确定从机器人处添加任务",
		})
		return
	}
	data, err := fetchdata.GetOne(auth)
	if err != nil {
		c.JSON(http.StatusOK, genericPostReply{
			Status:  "failed",
			Message: "没有查到用户数据，请确定从机器人处添加任务",
			Auth:    auth,
		})
		return
	}
	c.JSON(http.StatusOK, struct {
		status string
		data   fetchedSettingsReply
		auth   string
	}{
		status: "ok",
		data: fetchedSettingsReply{
			expire:   time.Now().Unix() + config.Expire,
			settings: config.Settings,
			override: data.Override,
		},
		auth: data.Auth,
	})
}

func postTaskAddSubmit(c *gin.Context) {
	data, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusOK, genericPostReply{
			Status:  "failed",
			Message: "数据提交失败了，请联系我并提供该条信息",
		})
		return
	}

	var parsed taskAddSubmitData
	err = json.Unmarshal(data, &parsed)
	if err != nil {
		c.JSON(http.StatusOK, genericPostReply{
			Status:  "failed",
			Message: "数据转换失败了，请联系我并提供该条信息",
		})
		return
	}

	ok := fetchdata.IfExist(parsed.Auth)
	if !ok {
		c.JSON(http.StatusOK, genericPostReply{
			Status:  "failed",
			Message: "验证数据来源出了问题，请联系我并提供该条信息",
		})
		return
	}

	var tid int32
	tid = rand.Int31()
	for analysistask.IfExist(tid) {
		tid = rand.Int31()
	}
	atask := analysistask.AnalysisTask{
		ID:          primitive.NewObjectID(),
		TaskID:      tid,
		Owner:       parsed.Data.Owner,
		Group:       parsed.Data.Group,
		Keywords:    parsed.Data.Keywords,
		MustMatch:   parsed.Data.MustMatch,
		Interval:    parsed.Data.Interval,
		TargetPrice: parsed.Data.TargetPrice,
		MaxPage:     parsed.Data.MaxPage,
		Sort:        "created_time",
		Order:       "desc",
	}
	adata := analysisdata.AnalysisData{
		ID:       primitive.NewObjectID(),
		Keywords: parsed.Data.Keywords,
		TaskID:   tid,
		Length:   0,
		Data:     nil,
	}
	err = analysisdata.Insert(adata)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, genericPostReply{
			Status:  "failed",
			Message: "数据添加程序出现问题，请联系我",
		})
		return
	}
	analysistask.AddTaskChannel <- atask

	c.JSON(http.StatusOK, genericPostReply{
		Status:  "success",
		Message: "任务添加请求提交成功，结果请通过查询进行查看",
	})
}
