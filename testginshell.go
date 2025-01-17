package main

import (
	"github.com/bigartists/Modi/client"
	"github.com/bigartists/Modi/src/helpers"
	"github.com/bigartists/Modi/src/wscore"
	"github.com/gin-gonic/gin"
	"k8s.io/client-go/tools/remotecommand"
	"log"
)

func main() {
	r := gin.New()
	r.GET("/", func(c *gin.Context) {
		wsClient, err := wscore.Upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		shellClient := wscore.NewWsShellClient(wsClient)

		err = helpers.HandleCommand("infra", "taichu-web-66f995596f-6wc8s", "taichu-web", client.K8sClient, client.K8sClientRestConfig, []string{"sh"}).Stream(
			remotecommand.StreamOptions{
				Stdin:  shellClient,
				Stdout: shellClient,
				Stderr: shellClient,
				Tty:    true,
			})
		if err != nil {
			log.Println(err)
		}
	})
	r.Run(":7777")
}
