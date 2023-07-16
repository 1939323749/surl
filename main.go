package main

import (
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/mritd/logger"
	"os"
	"os/signal"
	"surl/database"
	"surl/handler"
	"syscall"
	"time"
)

type CommonResp struct {
	Code      int
	Message   string
	Timestamp int64
}

func main() {
	database.Connect()
	app := fiber.New(fiber.Config{
		JSONEncoder: jsoniter.Marshal,
		Network:     "tcp",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(CommonResp{
				Code:      code,
				Message:   err.Error(),
				Timestamp: time.Now().Unix(),
			})
		}})

	go func() {
		sigs := make(chan os.Signal)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		for range sigs {
			logger.Warn("Received a termination signal, bark server shutdown...")
			if err := app.Shutdown(); err != nil {
				logger.Errorf("Server forced to shutdown error: %v", err)
			}
		}
	}()

	handler.RedirectHandler(app)
	handler.CreateShortUrlHandler(app)
	err := app.Listen(":8080")
	if err != nil {
		os.Kill.Signal()
	}
}
