package main

import (
	"github.com/zuijinbuzai/fundtop/api"
	"os"
	"os/exec"
	"fmt"
)

func main() {
	if len(os.Args) == 2 && len(os.Args[1]) == 6 {
		url := fmt.Sprintf("http://fund.eastmoney.com/%06s.html", os.Args[1])
		exec.Command("cmd", "/c", "start", url).Start()
		return
	}
	api.Work()
}