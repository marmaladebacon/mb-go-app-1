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
	AppName string
	BuiltAt string
	debug   = flag.Bool("d", false, "enables the debug mode")
	w       *astilectron.Window
)

var onWait = func(_ *astilectron.Astilectron, ws []*astilectron.Window, _ *astilectron.Menu, _ *astilectron.Tray, _ *astilectron.Menu) error {
	w = ws[0]
	go func() {
		time.Sleep(5 * time.Second)
		if err := bootstrap.SendMessage(w, "check.out.menu", "Don't forget to check out the menu!"); err != nil {
			astilog.Error(errors.Wrap(err, "sending check.out.menu event failed"))
		}
	}()
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
					if err := bootstrap.SendMessage(w,
						"about", htmlAbout, sendMsgErrorFunc); err != nil {
						astilog.Error(errors.Wrap(err, "sending about event failed"))
					}
					return
				},
			},
			{Role: astilectron.MenuItemRoleClose},
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
