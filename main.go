package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/love477/tools/demo"
	docs "github.com/love477/tools/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"k8s.io/client-go/util/workqueue"
)

// @BasePath /api/v1

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /example/helloworld [get]
func Helloworld(g *gin.Context) {
	g.JSON(http.StatusOK, "helloworld")
}

func main1() {
	go printGoroutine()
	go testWorkqueue()
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	{
		eg := v1.Group("/example")
		{
			eg.GET("/helloworld", Helloworld)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(":8080")
}

func main() {
	runGame()
}

func testWorkqueue() {
	stop := make(chan struct{})
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultItemBasedRateLimiter())
	for i := 0; i < 10000; i++ {
		v := i
		demo.HandleAndQueueEvent(demo.Create, v, queue)
	}
	c := demo.Controller{
		Queue:        queue,
		ResourceType: "pod",
		EventHandler: handler,
	}
	go c.Run(50, stop)
}

func handler(event *demo.Event) error {
	// fmt.Println("handler: ", event.Obj.(int))
	t := rand.Intn(3)
	time.Sleep(time.Duration(t) * time.Second)
	return nil
}

func printGoroutine() {
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		fmt.Println("goroutines: ", runtime.NumGoroutine())
	}
}

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello, World!")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func runGame() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
