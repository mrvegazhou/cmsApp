package api

import (
	"cmsApp/internal/errorx"
	gvalidator "cmsApp/pkg/validator"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	perrors "github.com/pkg/errors"
	"net/http"
)

type BaseController struct {
}

type SuccessResponse struct {
	DefaultResponse
	Data interface{} `json:"data" swaggertype:"object"` //接口返回的业务数据
}

type DefaultResponse struct {
	Status  int         `json:"status"`  //code 为1表示正常 0表示业务请求错误
	Message string      `json:"message"` //错误提示信息
	Data    interface{} `json:"data" swaggertype:"object"`
}

func (Base BaseController) Success(c *gin.Context, obj interface{}) {
	var res SuccessResponse
	res.Status = 200
	res.Message = "success"
	res.Data = obj

	c.JSON(http.StatusOK, res)
}

func (Base BaseController) Error(c *gin.Context, err error, data interface{}) {

	var res DefaultResponse
	//返回包装错误对应的最原始错误
	sourceErr := perrors.Cause(err)
	customErr, ok := sourceErr.(*errorx.CustomError)
	if ok {
		res.Status = customErr.ErrStatus
		res.Message = customErr.ErrMsg
		// 保存日志
		//if customErr.Err != nil {
		//	ctx, _ := c.Get("ctx")
		//	loggers.LogError(ctx.(context.Context), "api-custom-error", "err msg", map[string]string{
		//		"err message": err.Error(),
		//		"stack":   fmt.Sprintf("%+v", err),
		//	})
		//}
	} else {
		res.Status = errorx.HTTP_UNKNOW_ERR
		res.Message = err.Error()
		if data == nil {
			data = err.Error()
		}
		res.Data = data
	}
	c.JSON(http.StatusOK, res)
}

func (Base BaseController) FormBind(c *gin.Context, obj interface{}) error {

	trans, err := gvalidator.InitTrans("zh")
	if err != nil {
		return err
	}

	if err := c.ShouldBind(obj); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok && errs != nil {
			return errs
		}
		for _, v := range errs.Translate(trans) {
			return errors.New(v)
		}
		return err
	}
	return nil
}
