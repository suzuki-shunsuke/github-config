package lambda

import (
	"github.com/sirupsen/logrus"
)

func (handler *Handler) StartRepo() error {
	if err := handler.start("repo"); err != nil {
		logrus.WithError(err).Error("start")
		return err
	}
	return nil
}
