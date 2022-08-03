package miyabi

import (
	"fmt"
	"log"
	"os"
)

func requestLog(path, method string, status int64) {
	out := os.Stdout
	log.SetPrefix("[miyabi] ")
	logStr := fmt.Sprintf("| %d | %s | %s \n", status, method, path)
	fmt.Fprint(out, logStr)
}

/*
func routerLog(myb *Miyabi) {
	out := os.Stdout
	logStr := "\n"
	router := myb.Router
	for i := 0; i < len(router.RouterInfo); i++ {
		info := router.RouterInfo[i]
		logStr = fmt.Sprintf("%s[miyabi] %s \t %s \n", logStr, info.method, info.path)
	}
	for i := 0; i < len(router.Groups); i++ {
		info := router.Groups[i]
		logStr = fmt.Sprintf("%s[miyabi] group %s \n", logStr, info.basePath)
		for j := 0; j < len(info.GroupInfo); j++ {
			logStr = fmt.Sprintf("%s[miyabi]\t%s \t %s \n", logStr, info.GroupInfo[j].method, info.GroupInfo[j].path)
		}
	}
	fmt.Fprint(out, logStr)
}
*/
