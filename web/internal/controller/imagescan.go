package controller

import (
	"asyvy/pkg/task"
	"asyvy/web/pkg/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (ctrl *Controller) ImageScan(c *gin.Context) {
	var req model.ImageScanReq
	var resp model.ImageScanResp
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.Error(errors.WithStack(err))
		resp.Message = fmt.Sprintf("BindJSON error %s", err.Error())
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	imgstr := fmt.Sprintf("%s:%s", req.Image, req.Tag)
	// TODO: MYSQL INSERT DB
	task, err := task.ImageScanTask(imgstr, req.Username, req.Password)
	if err != nil {
		logrus.Error(errors.WithStack(err))
		resp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	info, err := ctrl.TaskClient.EnqueueContext(c.Copy().Request.Context(), task)
	if err != nil {
		logrus.Error(errors.WithStack(err))
		resp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	resp.Success = true
	resp.TaskID = info.ID
	c.JSON(http.StatusOK, resp)
}
