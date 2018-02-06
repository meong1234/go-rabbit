package util

import "github.com/sirupsen/logrus"

type Logger logrus.FieldLogger

var Log logrus.FieldLogger

func BuildContext(context string) logrus.Fields {
	return logrus.Fields{
		"context": context,
	}
}
