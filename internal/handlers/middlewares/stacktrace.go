package middlewares

import (
	"fmt"
	"os"
	"runtime"

	"beer-api/internal/core/config"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func stackTrace(c *fiber.Ctx, r interface{}) {
	var ok bool
	err, ok := r.(error)
	if !ok {
		err = fmt.Errorf("%v", r)
	}

	buf := make([]byte, 4<<10)
	buf = buf[:runtime.Stack(buf, false)]
	_, _ = os.Stderr.WriteString(fmt.Sprintf("[PANIC RECOVER]: %v\n%s\n", r, buf))
	logrus.Error(err)
	result := config.RR.CustomMessage(err.Error(), err.Error())
	_ = c.Status(result.HTTPStatusCode())
	_ = c.JSON(result)
}
