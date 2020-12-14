package controller

import (
	"github.com/gin-gonic/gin"
	"gitlab.eoitek.net/EOI/ckman/config"
	"gitlab.eoitek.net/EOI/ckman/model"
	"os"
	"syscall"
)

type ConfigController struct {
	signal chan os.Signal
}

func NewConfigController(ch chan os.Signal) *ConfigController {
	cf := &ConfigController{}
	cf.signal = ch
	return cf
}

// @Summary 更新配置
// @Description 更新配置
// @version 1.0
// @Security ApiKeyAuth
// @Param req body model.UpdateConfigReq true "request body"
// @Failure 200 {string} json "{"code":400,"msg":"请求参数错误","data":""}"
// @Failure 200 {string} json "{"code":5070,"msg":"更新配置失败","data":""}"
// @Success 200 {string} json "{"code":200,"msg":"success","data":nil}"
// @Router /api/v1/config [put]
func (cf *ConfigController) UpdateConfig(c *gin.Context) {
	var req model.UpdateConfigReq

	if err := model.DecodeRequestBody(c.Request, &req); err != nil {
		model.WrapMsg(c, model.INVALID_PARAMS, model.GetMsg(model.INVALID_PARAMS), err.Error())
		return
	}

	if len(req.Peers) > 0 {
		config.GlobalConfig.Server.Peers = req.Peers
	}

	if len(req.Prometheus) > 0 {
		config.GlobalConfig.Prometheus.Hosts = req.Prometheus
	}

	if err := config.MarshConfigFile(); err != nil {
		model.WrapMsg(c, model.UPDATE_CONFIG_FAIL, model.GetMsg(model.UPDATE_CONFIG_FAIL), err.Error())
		return
	}

	model.WrapMsg(c, model.SUCCESS, model.GetMsg(model.SUCCESS), nil)
	cf.signal <- syscall.SIGHUP
}
