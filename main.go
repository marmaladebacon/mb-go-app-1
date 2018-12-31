package main

import (
	"encoding/json"
	"flag"
	"time"

	astilectron "github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	astilog "github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

// Constants
const htmlAbout = `<b>Astilectron</b> demo base MarmaladeBacon test!<br>`

var (
	AppName       string
	BuiltAt       string
	debug         = flag.Bool("d", false, "enables the debug mode")
	primaryWindow *astilectron.Window
	app           *astilectron.Astilectron
	quitChannel   chan struct{}
)

var onWait = func(appPointer *astilectron.Astilectron, ws []*astilectron.Window, _ *astilectron.Menu, _ *astilectron.Tray, _ *astilectron.Menu) error {
	app = appPointer
	primaryWindow = ws[0]

	go func() {
		primaryWindow.OpenDevTools()
		time.Sleep(1 * time.Second)
		if err := bootstrap.SendMessage(primaryWindow, "check.out.menu", "Don't forget to check out the menu!"); err != nil {
			astilog.Error(errors.Wrap(err, "sending check.out.menu event failed"))
		}

	}()
	quitChannel := make(chan struct{})
	time.Sleep(1 * time.Second)
	pollForSymbol("aapl", quitChannel, primaryWindow, 3000)
	time.Sleep(1 * time.Second)
	pollForSymbol("fds", quitChannel, primaryWindow, 2000)
	setInterval(func() {
		if err := bootstrap.SendMessage(primaryWindow, "time.test", time.Now().Format(time.UnixDate)); err != nil {
			astilog.Error(errors.Wrap(err, "sending time.test event failed"))
		}
	}, 500, quitChannel)
	return nil
}

var sendMsgErrorFunc = func(m *bootstrap.MessageIn) {
	// Unmarshal payload
	var s string
	if err := json.Unmarshal(m.Payload, &s); err != nil {
		astilog.Error(errors.Wrap(err, "unmarshaling payload failed"))
		return
	}
	astilog.Infof("About modal has been displayed and payload is %s!", s)
}

func getMenuOptions() []*astilectron.MenuItemOptions {
	return []*astilectron.MenuItemOptions{{
		Label: astilectron.PtrStr("File"),
		SubMenu: []*astilectron.MenuItemOptions{
			{
				Label: astilectron.PtrStr("About"),
				OnClick: func(e astilectron.Event) (deleteListener bool) {
					if err := bootstrap.SendMessage(primaryWindow,
						"about", htmlAbout, sendMsgErrorFunc); err != nil {
						astilog.Error(errors.Wrap(err, "sending about event failed"))
					}
					return
				},
			},
			{
				Label: astilectron.PtrStr("Close App"),
				OnClick: func(e astilectron.Event) (deleteListener bool) {
					app.Quit()
					return
				},
			},
			//{Role: astilectron.MenuItemRoleClose},
		},
	}}
}

func main() {
	//Init
	flag.Parse()
	astilog.FlagInit()

	//Run bootstrap
	astilog.Debugf("Running app built at %s", BuiltAt)
	if err := bootstrap.Run(bootstrap.Options{

		AstilectronOptions: astilectron.Options{
			AppName:            AppName,
			AppIconDarwinPath:  "resources/icon.icns",
			AppIconDefaultPath: "resources/icon.png",
		},
		Debug:       *debug,
		MenuOptions: getMenuOptions(),
		OnWait:      onWait,
		//Unsure what this line below is for
		RestoreAssets: RestoreAssets,
		Windows: []*bootstrap.Window{{
			Homepage:       "index.html",
			MessageHandler: handleMessages,
			Options: &astilectron.WindowOptions{
				BackgroundColor: astilectron.PtrStr("#333"),
				Center:          astilectron.PtrBool(true),
				Height:          astilectron.PtrInt(480),
				Width:           astilectron.PtrInt(640),
			},
		}},
	}); err != nil {
		astilog.Fatal(errors.Wrap(err, "running bootstrap failed"))
	}
}
