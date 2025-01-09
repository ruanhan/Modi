package controllers

import (
	"github.com/bigartists/Modi/src/handler"
	"github.com/bigartists/Modi/src/result"
	"github.com/bigartists/Modi/src/service"
	"github.com/gin-gonic/gin"
)

type UserController struct {
}

func NewUserHandler() *UserController {
	return &UserController{}
}

// GET /users/123

// Build Build方法
func (this *UserController) Build(r *gin.RouterGroup) {
	r.GET("/users", UserList)
	r.GET("/user/:id", UserDetail)
	r.POST("/test", Test)
}

func UserList(c *gin.Context) {
	ret := ResultWrapper(c)(service.UserServiceGetter.GetUserList(), "")(OK)
	c.JSON(200, ret)
}

func UserDetail(c *gin.Context) {
	id := &struct {
		Id int64 `uri:"id" binding:"required"`
	}{}
	result.Result(c.ShouldBindUri(id)).Unwrap()
	ret := ResultWrapper(c)(service.UserServiceGetter.GetUserDetail(id.Id).Unwrap(), "")(OK)
	c.JSON(200, ret)
}

type TestUserReq struct {
	ID          string `validate:"required" json:"id"`
	UserID      string `json:"-"`
	CanDelete   bool   `json:"-"`
	CaptchaID   string `json:"captcha_id"`
	CaptchaCode string `json:"captcha_code"`
}

func Test(c *gin.Context) {
	req := &TestUserReq{}
	if handler.BindAndCheck(c, req) {
		return
	}
	handler.HandleResponse(c, nil, "success")
}
