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
	err := database.Connect()
	if err != nil {
		logger.Fatal(err)
		return
	}
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
			logger.Warn("Received a termination signal, server shutdown...")
			if err := app.Shutdown(); err != nil {
				logger.Errorf("Server forced to shutdown error: %v", err)
			}
		}
	}()

	handler.RedirectHandler(app)
	handler.CreateShortUrlHandler(app)
	handler.ClickHandler(app)
	err = app.Listen(":8080")
	if err != nil {
		logger.Errorf("Server startup error: %v", err)
		os.Kill.Signal()
	}
}
