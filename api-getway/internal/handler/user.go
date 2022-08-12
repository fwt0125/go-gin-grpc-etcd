package handler

import (
	"api-getway/internal/service"
	"api-getway/pkg/e"
	"api-getway/pkg/res"
	"api-getway/pkg/util"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserRegister(ginCtx *gin.Context) {
	var userReq service.UserRequest
	PanicIfUserError(ginCtx.Bind(&userReq))

	//gin.key中取出实例
	userService := ginCtx.Keys["user"].(service.UserServiceClient)
	userRegister, err := userService.UserRegister(context.Background(), &userReq)
	PanicIfUserError(err)
	r := res.Response{
		Data:   userRegister,
		Status: uint(userRegister.Code),
		Msg:    e.GetMsg(uint(userRegister.Code)),
	}

	ginCtx.JSON(http.StatusOK, r)

}

func UserLogin(ginCtx *gin.Context) {
	var userReq service.UserRequest
	PanicIfUserError(ginCtx.Bind(&userReq))

	//gin.key中取出实例
	userService := ginCtx.Keys["user"].(service.UserServiceClient)
	userRegister, err := userService.UserLogin(context.Background(), &userReq)
	PanicIfUserError(err)
	token, err := util.GenerateToken(uint(userRegister.UserDetail.UserId))
	r := res.Response{
		Data: res.TokenData{
			User:  userRegister.UserDetail,
			Token: token,
		},
		Status: uint(userRegister.Code),
		Msg:    e.GetMsg(uint(userRegister.Code)),
	}

	ginCtx.JSON(http.StatusOK, r)

}
