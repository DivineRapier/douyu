package main

import (
	"fmt"

	"github.com/DivineRapier/douyu"
)

func main() {
	dy, b := douyu.OpenDanmu(534740)
	fmt.Println(dy)
	fmt.Println(b)
	fmt.Println()
	fmt.Println()

	dy.JoinGroupRequest(0)
	dy.ShowChatmessage()
	wait()
}

func wait() {
	c := make(chan struct{})
	<-c
}
