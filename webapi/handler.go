package webapi

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"

	"bookq.xyz/mercariWatchdog/utils/analysisdata"
	"bookq.xyz/mercariWatchdog/utils/analysistask"
	"bookq.xyz/mercariWatchdog/utils/fetchdata"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
			Message: "没有用户数据，请确定从机器人处添加任务",
			Auth:    auth,
		})
		return
	}
	c.JSON(http.StatusOK, struct {
		status string
		data   fetchdata.TaskAddFetchData
	}{
		status: "ok",
		data:   data,
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
