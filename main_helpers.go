package main

import (
	"fmt"
	"time"
)

func setInterval(funcToLoop func(), timeInMs int, quitChannel chan struct{}) {
	interval := time.Duration(timeInMs) * time.Millisecond
	ticker := time.NewTicker(interval)

	go func() {
		for {
			select {
			case <-ticker.C:
				funcToLoop()
				fmt.Println(time.Now())
			case <-quitChannel:
				ticker.Stop()
				fmt.Println("Stopped the ticker!")
				return
			}
		}
	}()

}
