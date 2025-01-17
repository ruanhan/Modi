package controllers

import (
	"github.com/bigartists/Modi/client"
	"github.com/bigartists/Modi/src/handler"
	"github.com/bigartists/Modi/src/helpers"
	"github.com/bigartists/Modi/src/wscore"
	"github.com/gin-gonic/gin"
	"k8s.io/client-go/tools/remotecommand"
)

type TerminalController struct {
}

func ProviderTerminalController() *TerminalController {
	return &TerminalController{}
}

func (this *TerminalController) PodConnect(c *gin.Context) {
	ns := c.Query("ns")
	pod := c.Query("pod")
	container := c.Query("c")
	wsClient, err := wscore.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		handler.NewRespBodyFromError(handler.NewCustomError().BadRequest(err.Error()))
	}
	shellClient := wscore.NewWsShellClient(wsClient)

	err = helpers.HandleCommand(ns, pod, container, client.K8sClient, client.K8sClientRestConfig, []string{"sh"}).
		Stream(remotecommand.StreamOptions{
			Stdin:  shellClient,
			Stdout: shellClient,
			Stderr: shellClient,
			Tty:    true,
		})
}

func (this *TerminalController) Build(r *gin.RouterGroup) {
	r.GET("/podws", this.PodConnect)
}
