package main

import (
	// "encoding/json"
	"fmt"

	"github.com/rakyll/portmidi"
	"github.com/scgolang/sc"
	// "os"
)

func main() {
	var synthID int32
	const synthName = "sineTone"

	// Set up SuperCollider client.
	client, err := sc.NewClient("udp", "127.0.0.1:57121", "127.0.0.1:57120")
	if err != nil {
		panic(err)
	}
	_, err = client.AddDefaultGroup()
	if err != nil {
		panic(err)
	}
	err = client.DumpOSC(int32(1))
	if err != nil {
		panic(err)
	}
	def := sc.NewSynthdef(synthName, func(p sc.Params) sc.Ugen {
		freq := p.Add("freq", 440)
		gain := p.Add("gain", 0.5)
		bus := sc.C(0)
		env := sc.EnvGen{
			Env:        sc.EnvPerc{},
			Done:       sc.FreeEnclosing,
			LevelScale: gain,
		}.Rate(sc.KR)
		sig := sc.SinOsc{Freq: freq}.Rate(sc.AR).Mul(env)
		return sc.Out{bus, sig}.Rate(sc.AR)
	})
	err = client.SendDef(def)
	if err != nil {
		panic(err)
	}

	// initialize midi
	portmidi.Initialize()
	// this code can be uncommented to discover the
	// device ID's portmidi comes up with
	// deviceCount := portmidi.CountDevices()
	// enc := json.NewEncoder(os.Stdout)
	// for i := 0; i < deviceCount; i++ {
	// 	info := portmidi.GetDeviceInfo(portmidi.DeviceId(i))
	// 	fmt.Printf("device %d - ", i)
	// 	err = enc.Encode(info)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	// setup midi input stream and listen for midi events
	in, err := portmidi.NewInputStream(3, 1024)
	if err != nil {
		panic(err)
	}
	ch := in.Listen()
	for event := range ch {
		if event.Status == 144 {
			// MIDI note
			fmt.Printf("Note %-3d Velocity %-3d\n", event.Data1, event.Data2)
			if event.Data2 > 0 {
				// Note On
				synthID = client.NextSynthID()
				ctls := map[string]float32{
					"freq": sc.Midicps(float32(event.Data1)),
					"gain": float32(event.Data2 / 127.0),
				}
				_, err = client.Synth(synthName, synthID, sc.AddToTail, sc.DefaultGroupID, ctls)
			}
		}
	}
	portmidi.Terminate()
}
