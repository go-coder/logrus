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

	log.Info("infommmmmmsg", "infokey", "infoval")
	log.V(2).Info("v2infomsg", "v2infokey", "v2infoval")
	log.Error(nil, "errmsg", "errkey", "errval")
	log.Error(&Err{err_msg:"myerr"}, "errmsg", "errkey", "errval")
}