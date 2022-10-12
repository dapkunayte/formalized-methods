package plots

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"os"
)

func DrawLines(x, y, yNew []float64, name string, method string) {
	line := charts.NewLine()
	// set some global options like Title/Legend/ToolTip or anything else
	items1 := make([]opts.LineData, len(y))
	items2 := make([]opts.LineData, len(y))
	for i := 0; i < len(y); i++ {
		items1[i] = opts.LineData{Value: y[i]}
	}
	for i := 0; i < len(yNew); i++ {
		items2[i] = opts.LineData{Value: yNew[i]}
	}
	x2 := x[1 : len(x)-1]
	// Put data into instance
	switch method {
	case "mean3":
		line.SetXAxis(x2).
			AddSeries("Category A", items1[1:len(y)-1]).
			AddSeries("Category B", items2).
			SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: false}))
	case "lsq":
		line.SetXAxis(x).
			AddSeries("Category A", items1).
			AddSeries("Category B", items2).
			SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: false}))
	case "exp":
		line.SetXAxis(x).
			AddSeries("Category A", items1).
			AddSeries("Category B", items2).
			SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: false}))
	}

	f, _ := os.Create(name)
	line.Render(f)
}
