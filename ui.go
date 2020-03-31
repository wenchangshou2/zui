package main

import (
	"fmt"
	"github.com/lxn/walk"
	"log"
)
type Zui struct {
	mw        *walk.MainWindow
	exit_chan chan bool
}
func (zui *Zui) Init()error{
	var (
		err error
	)
	zui.mw,err=walk.NewMainWindow()
	if err!=nil{
		return err
	}
	return nil
}

func (zui *Zui) Start() {
	fmt.Println("zui",zui.mw)
	// We load our icon from a file.
	//iconPath,_:=zutil.GetFullPath("img/stop.ico")
	icon, err := walk.Resources.Icon("img/stop.png")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	// Create the notify icon and make sure we clean it up on exit.
	ni, err := walk.NewNotifyIcon(zui.mw)
	if err != nil {
		log.Fatal(err)
	}
	defer ni.Dispose()

	// Set the icon and a tool tip text.
	//if err := ni.SetIcon(icon); err != nil {
	//	log.Fatal(err)
	//}
	if err := ni.SetToolTip("Click for info or use the context menu to exit."); err != nil {
		log.Fatal(err)
	}

	// When the left mouse button is pressed, bring up our balloon.
	ni.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
		if button != walk.LeftButton {
			return
		}

		if err := ni.ShowCustom(
			"Walk NotifyIcon Example",
			"There are multiple ShowX methods sporting different icons.",
			icon); err != nil {

			log.Fatal(err)
		}
	})

	// We put an exit action into the context menu.
	exitAction := walk.NewAction()
	if err := exitAction.SetText("E&xit"); err != nil {
		log.Fatal(err)
	}
	exitAction.Triggered().Attach(func() { zui.exit_chan<-true })
	if err := ni.ContextMenu().Actions().Add(exitAction); err != nil {
		log.Fatal(err)
	}

	// The notify icon is hidden initially, so we have to make it visible.
	if err := ni.SetVisible(true); err != nil {
		log.Fatal(err)
	}

	// Run the message loop.
	zui.mw.Run()
}
var (
	G_UI *Zui
)
func InitUI(exit chan bool) error {
	G_UI=&Zui{
		exit_chan:exit,
	}
	if err:=G_UI.Init();err!=nil{
		fmt.Println("init ui err"+err.Error())
		return err
	}
	G_UI.Start()
	return nil
}
