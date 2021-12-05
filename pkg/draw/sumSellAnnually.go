package draw

import (
	"bytes"
	"cloudProject/pkg/utils"
	"github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
)

func SumSellAnnually(data []chart.Value) ([]byte, error) {
	graph := chart.BarChart{
		Background: chart.Style{
			Padding: chart.Box{
				Top:  50,
				Left: 10,
			},
			StrokeWidth: 5,
			StrokeColor: drawing.ColorBlack,
		},
		Height:   1024,
		Width:    2048,
		BarWidth: 30,
	}
	var series []chart.Series //help for graph

	for _, value := range data {
		r := utils.RandInt(1)
		g := utils.RandInt(1)
		b := utils.RandInt(1)

		graph.Bars = append(graph.Bars, chart.Value{
			Style: chart.Style{FillColor: drawing.ColorFromHex(r + g + b), StrokeColor: drawing.ColorFromHex(r + g + b)},
			Label: value.Label,
			Value: value.Value,
		})
		series = append(series, chart.ContinuousSeries{Name: value.Label, Style: chart.Style{FillColor: drawing.ColorFromHex(r + g + b), StrokeColor: drawing.ColorFromHex(r + g + b)}})
	}

	//adding graph details to graph
	graph.Elements = append(graph.Elements, []chart.Renderable{
		chart.Legend(&chart.Chart{Series: series}),
	}...)
	buff := bytes.NewBuffer([]byte{})
	_ = graph.Render(chart.PNG, buff)
	return buff.Bytes(), nil
}
