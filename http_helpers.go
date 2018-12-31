package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	astilectron "github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	astilog "github.com/asticode/go-astilog"
	"github.com/pkg/errors"
	//tb "github.com/nsf/termbox-go"
)

type DelayedQuote struct {
	Symbol           string
	DelayedPrice     float64
	DelayedSize      float64
	DelayedPriceTime float64
	ProcessedTime    float64
	Text             string
}

func (c *DelayedQuote) GetStr() string {
	return fmt.Sprintf("%v is at %v, %v", c.Symbol, c.DelayedPrice, time.Now())
}

func pollForSymbol(symbol string, quitChannel chan struct{}, primaryWindow *astilectron.Window, updateInterval int) {
	dataChannel := make(chan DelayedQuote)
	getFunc := makeGetIntervalFunc("stock/"+symbol+"/delayed-quote", dataChannel)
	setInterval(getFunc, updateInterval, quitChannel)
	sendMsgFunc := makeSendMessageFunc(dataChannel, quitChannel, primaryWindow)
	go sendMsgFunc()
	//printingFunc := makePrintFunc(dataChannel, quitChannel)
	//go printingFunc()
}

func makeSendMessageFunc(dataChannel chan DelayedQuote, quitChannel chan struct{}, primaryWindow *astilectron.Window) func() {
	return func() {
		for {
			select {
			case data := <-dataChannel:
				bootstrap.SendMessage(primaryWindow, "ticker.track", data)
			case <-quitChannel:
				fmt.Println("Stopped the ticker!")
				close(dataChannel)
				return
			}
		}
	}
}

func makeGetIntervalFunc(apiStrPartial string, dataChannel chan DelayedQuote) func() {

	return func() {
		d := DelayedQuote{}
		apiStr := "https://api.iextrading.com/1.0/" + apiStrPartial
		resp, err := http.Get(apiStr)
		if err == nil {
			json.NewDecoder(resp.Body).Decode(&d)
			d.Text = d.GetStr()
			dataChannel <- d
			//fmt.Println(d.DelayedPrice)
		} else {
			fmt.Println(err)
		}
	}
}

//need to send message as a struct here
func makeSendMsgFunc(dataChannel chan DelayedQuote, quitChannel chan struct{}, primaryWindow *astilectron.Window) func() {
	return func() {
		defer close(quitChannel)
		for {
			select {
			case data := <-dataChannel:
				fmt.Println(data.GetStr())
				if err := bootstrap.SendMessage(primaryWindow, "check.out.menu", "Don't forget to check out the menu!"); err != nil {
					astilog.Error(errors.Wrap(err, "sending check.out.menu event failed"))
				}
			case <-quitChannel:
				fmt.Println("Stopped the ticker!")
				close(dataChannel)
				break
			}
		}
	}
}
