package draw

import (
	"bytes"
	"github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
)

func GenreSells(from, to string, data []chart.Value) ([]byte, error) {
	graph := chart.BarChart{
		Title: from + " to " + to,
		Background: chart.Style{
			Padding: chart.Box{
				Top:    50,
				Right:  20,
				Left:   20,
				Bottom: 30,
			},
			StrokeWidth: 5,
			StrokeColor: drawing.ColorBlack,
		},
		Height:   512,
		BarWidth: 30,
	}
	for _, value := range data {
		graph.Bars = append(graph.Bars, value)
	}

	buff := bytes.NewBuffer([]byte{})
	_ = graph.Render(chart.PNG, buff)
	return buff.Bytes(), nil
}
