package draw

import (
	"bytes"
	"github.com/wcharczuk/go-chart/v2"
	"strconv"
)

func ProducerSells(data map[string][]chart.Value) ([]byte, error) {

	graph := chart.Chart{
		Background: chart.Style{
			Padding: chart.Box{
				Top:  20,
				Left: 20,
			},
		},
	}
	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}
	for key, value := range data {
		ch := chart.ContinuousSeries{
			Name: key,
		}
		for _, v := range value {
			p, _ := strconv.ParseFloat(v.Label, 10)
			ch.YValues = append(ch.YValues, v.Value)
			ch.XValues = append(ch.XValues, p)
		}
		graph.Series = append(graph.Series, ch)
	}
	buff := bytes.NewBuffer([]byte{})
	_ = graph.Render(chart.PNG, buff)
	return buff.Bytes(), nil
}
