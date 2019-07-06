package main

import (
	"github.com/go-coder/logrus"
	rus "github.com/sirupsen/logrus"
)

type Err struct {
	err_msg string
}

func (e Err) Error() string {
	return e.err_msg
}

func main() {
	log := logrus.NewLogger("logname", rus.New())

	log.Info("infomsg", "infokey", "infoval")
	log.V(5).Info("v5infomsg", "v5infokey", "v5infoval")
	log.WithValues("withkey", "withval").Error(nil, "errmsg", "errkey", "errval")
	log.WithName("myprefix").Error(&Err{err_msg: "myerr"}, "errmsg", "errkey", "errval")
}
