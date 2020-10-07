package miyabi

import (
	"fmt"
	"log"
)

func requestLog(path, method string, status int64) {
	log.SetPrefix("[miyabi] ")
	logStr := fmt.Sprintf("| %d | %s | %s \n", status, method, path)
	log.Print(logStr)
}

func routerLog(myb *Miyabi) {
	logStr := "\n"
	router := myb.Router
	for i := 0; i < len(router.RouterInfo); i++ {
		info := router.RouterInfo[i]
		logStr += fmt.Sprintf("[miyabi] %s \t %s \n", info.method, info.path)
	}
	for i := 0; i < len(router.Groups); i++ {
		info := router.Groups[i]
		logStr += fmt.Sprintf("[miyabi] group %s \n", info.basePath)
		for j := 0; j < len(info.GroupInfo); j++ {
			logStr += fmt.Sprintf("[miyabi]\t%s \t %s \n", info.GroupInfo[j].method, info.GroupInfo[j].path)
		}
	}
	log.Print(logStr)
}
