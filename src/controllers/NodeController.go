package controllers

import (
	"github.com/bigartists/Modi/src/handler"
	"github.com/bigartists/Modi/src/service"
	"github.com/gin-gonic/gin"
)

type NodeController struct {
	nodeService service.INode
}

func ProviderNodeController(nodeService service.INode) *NodeController {
	return &NodeController{
		nodeService: nodeService,
	}
}

func (this *NodeController) Build(r *gin.RouterGroup) {
	r.GET("/nodes", this.listNodes)
}

func (this *NodeController) listNodes(c *gin.Context) {
	namespace := c.Query("ns")
	nodes := this.nodeService.ListAllNodes(namespace)
	if nodes == nil {
		handler.NewRespBodyFromError(handler.NewCustomError().BadRequest("no such node"))
		return
	}
	handler.HandleResponse(c, nil, nodes)
}
