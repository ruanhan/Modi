package app

import (
	"fmt"
	dbs "github.com/bigartists/Modi/src/dao"
	"github.com/bigartists/Modi/src/listener"
	"github.com/bigartists/Modi/src/middlewares"
	"github.com/bigartists/Modi/src/routes"
	"github.com/bigartists/Modi/src/validators"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

// 本函数用于初始化gin
// 有一个测试路由 GET  /test
// 函数接收一个参数，用于指定监听的端口

func Run(port int) error {
	// 执行命令行
	dbs.InitDB()
	// 初始化k8s client
	listener.InitInformerListener()
	r := gin.New()

	r.Use(middlewares.JwtAuthMiddleware())
	r.Use(middlewares.ErrorHandler())

	// 加载路由
	routes.Build(r)
	// 加载 validator
	validators.Build()

	err := r.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	return nil
}

// NewApiServerCommand 初始化命令行参数
func NewApiServerCommand() (cmd *cobra.Command) {
	// 集成 cobra命令
	cmd = &cobra.Command{
		Use: "appServer",
		RunE: func(cmd *cobra.Command, args []string) error {
			port, err := cmd.Flags().GetInt("port")
			if err != nil {
				return err
			}
			return Run(port)
		},
	}
	// 添加 flag, name=port, 默认值是 9090
	cmd.Flags().Int("port", 9090, "appserver port")
	return
}
