package draw

import (
	"bytes"
	"cloudProject/pkg/utils"
	"github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
)

func SellsCompare2Game(data map[string][]chart.Value) ([]byte, error) {
	graph := chart.BarChart{
		BarSpacing: 20,
		Background: chart.Style{
			Padding: chart.Box{
				Top: 50,
			},

			StrokeWidth: 5,
			StrokeColor: drawing.ColorBlack,
		},
		Height: 512,
		//BarWidth: 50,
		//Width: 700,
	}
	var series []chart.Series //help for graph
	counter := 0
	colCounter := 0
	var colors = make([]string, 0)
	for k, v := range data {
		for i, vv := range v {
			if len(colors) == 0 || colCounter >= len(colors) {
				colors = append(colors, utils.RandInt(1)+utils.RandInt(1)+utils.RandInt(1))
			}

			if i == 1 {
				//set one label only
				graph.Bars = append(graph.Bars, chart.Value{Value: vv.Value, Label: k, Style: chart.Style{FillColor: drawing.ColorFromHex(colors[colCounter]), StrokeColor: drawing.ColorFromHex(colors[colCounter])}})
			} else {
				graph.Bars = append(graph.Bars, chart.Value{Value: vv.Value, Style: chart.Style{FillColor: drawing.ColorFromHex(colors[colCounter]), StrokeColor: drawing.ColorFromHex(colors[colCounter])}})
			}
			if counter == 0 {
				series = append(series, chart.ContinuousSeries{Name: vv.Label, Style: chart.Style{FillColor: drawing.ColorFromHex(colors[colCounter]), StrokeColor: drawing.ColorFromHex(colors[colCounter])}})
			}
			colCounter++
		}
		//between of them
		graph.Bars = append(graph.Bars, []chart.Value{{Value: 0, Label: "", Style: chart.Style{FillColor: drawing.ColorWhite, StrokeColor: drawing.ColorWhite}},
			{Value: 0, Label: "", Style: chart.Style{FillColor: drawing.ColorWhite, StrokeColor: drawing.ColorWhite}},
			{Value: 0, Label: "", Style: chart.Style{FillColor: drawing.ColorWhite, StrokeColor: drawing.ColorWhite}},
			{Value: 0, Label: "", Style: chart.Style{FillColor: drawing.ColorWhite, StrokeColor: drawing.ColorWhite}},
		}...)
		//for new graph on old graph
		counter++
		colCounter = 0
	}
	//adding graph details to graph
	graph.Elements = append(graph.Elements, []chart.Renderable{
		chart.Legend(&chart.Chart{Series: series}),
	}...)

	buff := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buff)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}
