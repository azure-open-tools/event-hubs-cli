package common

import (
	"fmt"
	"github.com/vbauerster/mpb/v5"
	"github.com/vbauerster/mpb/v5/decor"
	"log"
	"os"
)

func CloseFile(file *os.File) {
	if file != nil {
		err := file.Close()
		if err != nil {
			log.Fatalf("open file error: %v", err)
		}
	}
}

func GetBar(size int64, taskIndex int, progress *mpb.Progress) *mpb.Bar {
	task := fmt.Sprintf("Task - %02d:", taskIndex)

	var bar = progress.AddBar(size,
		mpb.BarStyle("╢▌▌░╟"),
		mpb.PrependDecorators(
			decor.Name(task, decor.WC{W: len(task) + 1, C: decor.DidentRight}),
			decor.Name("sending", decor.WCSyncSpaceR),
			decor.CountersNoUnit("%d / %d - ", decor.WCSyncWidth),
			decor.Elapsed(decor.ET_STYLE_GO, decor.WC{W: 3}),
			//decor.OnComplete(
			//	decor.AverageETA(decor.ET_STYLE_GO, decor.WC{W: 4}), "done",
			//	),
		),
		mpb.AppendDecorators(
			decor.Percentage(decor.WC{W: 5}),
		),
	)

	return bar
}