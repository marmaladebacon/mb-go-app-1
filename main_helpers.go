package main

import (
	"fmt"
	"time"
)

func setInterval(funcToLoop func(), timeInMs int) chan struct{} {
	interval := time.Duration(timeInMs) * time.Millisecond
	ticker := time.NewTicker(interval)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				funcToLoop()
				fmt.Println(time.Now())
			case <-quit:
				ticker.Stop()
				fmt.Println("Stopped the ticker!")
				return
			}
		}
	}()
	return quit
}
