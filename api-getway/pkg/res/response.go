package res

import (
	"api-getway/pkg/e"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Status uint        `json:"Status"`
	Data   interface{} `json:"Data"`
	Msg    string      `json:"Msg"`
	Error  string      `json:"Error"`
}

type DataList struct {
	Item  interface{} `json:"Item"`
	Total uint        `json:"Total"`
}

type TokenData struct {
	User  interface{} `json:"User"`
	Token string      `json:"Token"`
}

func ginH(msgCode int, data interface{}) gin.H {
	return gin.H{
		"code": msgCode,
		"msg":  e.GetMsg(uint(msgCode)),
		"data": data,
	}
}
