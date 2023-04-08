package beer

import (
	"beer-api/internal/models"
	"fmt"
	"io"
	"os"
	"time"

	"beer-api/internal/core/context"

	"github.com/sirupsen/logrus"
)

func (s *service) uploadFile(c *context.Context) (*models.File, error) {
	var fileCreate *models.File
	file, _ := c.FormFile("file")
	if file != nil {
		src, err := file.Open()
		if err != nil {
			logrus.Errorf("open file error: %s", err)
			return nil, err
		}

		defer func() { _ = src.Close() }()

		newFileName := fmt.Sprintf("beer_%s_%s", time.Now().Format("2006010215040506"), file.Filename)

		f, err := os.Create("./public/" + newFileName)
		if err != nil {
			logrus.Errorf("create file err: %s", err)
			return nil, err
		}

		defer func() {
			_ = f.Close()
		}()

		_, err = io.Copy(f, src)
		if err != nil {
			logrus.Errorf("copy file err: %s", err)
			return nil, err
		}

		fileCreate = &models.File{
			Name: newFileName,
		}
		err = s.fileRepository.Create(c.GetDatabase(), fileCreate)
		if err != nil {
			logrus.Errorf("create file error: %s", err)
			return nil, err
		}
	}

	return fileCreate, nil
}
