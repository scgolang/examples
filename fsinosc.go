package main

import (
	"fmt"

	"github.com/scgolang/sc"
)

func main() {
	const synthName = "FSinOscExample"
	client := sc.NewClient("127.0.0.1:57112")
	err := client.Connect("127.0.0.1:57110")
	if err != nil {
		panic(err)
	}
	defaultGroup, err := client.AddDefaultGroup()
	if err != nil {
		panic(err)
	}
	def := sc.NewSynthdef(synthName, func(p sc.Params) sc.Ugen {
		bus := sc.C(0)
		line := sc.XLine{sc.C(4), sc.C(401), sc.C(8), 0}.Rate(sc.KR)
		sin1 := sc.SinOsc{line, sc.C(0)}.Rate(sc.AR).MulAdd(sc.C(200), sc.C(800))
		sin2 := sc.SinOsc{Freq: sin1}.Rate(sc.AR).Mul(sc.C(0.2))
		return sc.Out{bus, sin2}.Rate(sc.AR)
	})
	err = client.SendDef(def)
	if err != nil {
		panic(err)
	}
	synthID := client.NextSynthID()
	_, err = defaultGroup.Synth(synthName, synthID, sc.AddToTail, nil)
	fmt.Printf("created synth %d\n", synthID)
}
