package controllers

import (
	"net/http"

	"github.com/bigartists/Modi/src/handler"

	"github.com/bigartists/Modi/src/repo"
	"github.com/gin-gonic/gin"
)

type ServiceController struct {
	serviceRepo *repo.ServiceRepo
}

func ProviderServiceController(serviceRepo *repo.ServiceRepo) *ServiceController {
	return &ServiceController{
		serviceRepo: serviceRepo,
	}
}

// ListServices godoc
// @Summary 获取Service列表
// @Description 获取指定namespace下的Service列表，如果namespace为空则获取所有namespace的Service
// @Tags Service
// @Accept json
// @Produce json
// @Param namespace query string false "命名空间"
// @Success 200 {array} corev1.Service
// @Router /api/v1/services [get]
func (s *ServiceController) ListServices(c *gin.Context) {
	namespace := c.Query("ns")
	services, err := s.serviceRepo.List(namespace)
	if err != nil {
		handler.NewRespBodyFromError(handler.NewCustomError().BadRequest(err.Error()))
		return
	}
	handler.HandleResponse(c, nil, services)
}

// GetService godoc
// @Summary 获取Service详情
// @Description 获取指定namespace下的Service详情
// @Tags Service
// @Accept json
// @Produce json
// @Param namespace path string true "命名空间"
// @Param name path string true "Service名称"
// @Success 200 {object} corev1.Service
// @Router /api/v1/namespaces/{namespace}/services/{name} [get]
func (s *ServiceController) GetService(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	service := s.serviceRepo.Get(namespace, name)

	if service == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
		return
	}
	c.JSON(http.StatusOK, service)
}

func (s *ServiceController) Build(api *gin.RouterGroup) {
	api.GET("/services", s.ListServices)
	api.GET("/services/:namespace/:name", s.GetService)
}
