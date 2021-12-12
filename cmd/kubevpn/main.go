package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/wencaiwulue/kubevpn/cmd/kubevpn/cmds"
	"github.com/wencaiwulue/kubevpn/util"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	if !util.IsAdmin() {
		util.RunWithElevated()
	} else {
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
		_ = cmds.RootCmd.Execute()
	}
}