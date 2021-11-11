package main

import (
	dither "github.com/esimov/dithergo"
	"github.com/flipperdevices/go-flipper"
	"github.com/go-vgo/robotgo"
	"github.com/wmarbut/goxbm"
	"image"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var scale = robotgo.SysScale()

var ditherer = dither.Dither{
	Type: "Sierra-Lite",
	Settings: dither.Settings{
		Filter: [][]float32{
			{0.0, 0.0, 2.0 / 4.0},
			{1.0 / 4.0, 1.0 / 4.0, 0.0},
			{0.0, 0.0, 0.0},
		},
	},
}

func main() {
	rand.Seed(time.Now().Unix())

	ser, err := initCli(os.Args[1])
	if err != nil {
		log.Println("Can't init RPC", err)
		os.Exit(1)
	}

	f, err := flipper.Connect(ser)
	if err != nil {
		log.Println("Can't connect to RPC", err)
		os.Exit(1)
	}

	log.Println("Connected!")

	err = f.Gui.StartVirtualDisplay()
	if err != nil {
		log.Fatalln("Can't start Virtual Display", err)
	}

	t := time.NewTicker(time.Millisecond * 50)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-t.C:
			bimg := getScreenPart()
			err = f.Gui.UpdateVirtualDisplay(goxbm.ToRawXBMBytes(bimg))
			if err != nil {
				log.Fatalln("Can't update Virtual Display", err)
			}
			break
		case <-c:
			log.Println("Stopping...")
			err = f.Gui.StopVirtualDisplay()
			if err != nil {
				log.Fatalln("Can't stop Virtual Display", err)
			}
			return
		}
	}
}

func getScreenPart() image.Image {
	x, y := robotgo.GetMousePos()

	img := robotgo.CaptureScreen(
		x-int(float64(64)/scale),
		y-int(float64(32)/scale),
		int(float64(128)/scale),
		int(float64(64)/scale))
	defer robotgo.FreeBitmap(img)

	bimg := ditherer.Monochrome(robotgo.ToImage(img), 1.18)

	return bimg
}
