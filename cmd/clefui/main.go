package main

import (
	"os"
	"context"
	"os/signal"
	"fmt"
	"log"

	"github.com/kyokan/clef-ui/pkg/rpc"
	"github.com/kyokan/clef-ui/pkg/clefclient"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/quick"
	"github.com/therecipe/qt/quickcontrols2"
	"github.com/therecipe/qt/widgets"
)


func startQT() {

	// enable high dpi scaling
	// useful for devices with high pixel density displays
	// such as smartphones, retina displays, ...
	core.QCoreApplication_SetAttribute(core.Qt__AA_EnableHighDpiScaling, true)

	// needs to be called once before you can start using QML/Quick
	widgets.NewQApplication(len(os.Args), os.Args)

	// use the material style
	// the other inbuild styles are:
	// Default, Fusion, Imagine, Universal
	quickcontrols2.QQuickStyle_SetStyle("Material")

	// create the quick view
	// with a minimum size of 250*200
	// set the window title to "Hello QML/Quick Example"
	// and let the root item of the view resize itself to the size of the view automatically
	view := quick.NewQQuickView(nil)
	view.SetMinimumSize(core.NewQSize2(250, 200))
	view.SetResizeMode(quick.QQuickView__SizeRootObjectToView)
	view.SetTitle("Hello QML/Quick Example")

	// load the embeeded qml file
	// created by either qtrcc or qtdeploy
	// view.SetSource(core.NewQUrl3("qrc:/qml/main.qml", 0))
	// you can also load a local file like this instead:
	//view.SetSource(core.QUrl_FromLocalFile("./qml/main.qml"))

	// make the view visible
	view.Show()

	// start the main Qt event loop
	// and block until app.Exit() is called
	// or the window is closed by the user
	widgets.QApplication_Exec()
}


func main() {
	// trap Ctrl+C and call cancel on the context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	cancelChannel := make(chan os.Signal, 1)
	done := make(chan bool)
	signal.Notify(cancelChannel, os.Interrupt)

	stdin, stdout, stderr, err := clefclient.StartClef()
	if err != nil {
		log.Panicf("Cannot stard clef: %s", err)
		return
	}

	server := rpc.NewServer()
	server.ListenStdIO(ctx, stdin, stdout, stderr)

	// Watch for os interrupt
	go func() {
		<-cancelChannel
		cancel()
		fmt.Print("\n")
		log.Print("Stopped Clef UI.")
		signal.Stop(cancelChannel)
		done <- true
	}()

	startQT()
	// Exit when done
	<-done
}
