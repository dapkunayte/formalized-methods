package main

import (
	"fmt"
	"github.com/aclements/go-moremath/mathx"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"math"
	"os"
	"sort"
)

var KenwallWMatrix = map[int]map[int]float64{
	3: {
		5: 64.4,
		6: 103.9,
		7: 157.3,
	},
	4: {
		4: 49.5,
		5: 88.4,
		6: 143.3,
		7: 217.0,
	},
	5: {
		4: 62.6,
		5: 112.3,
		6: 182.4,
		7: 276.2,
	},
	6: {
		4: 75.7,
		5: 136.1,
		6: 221.4,
		7: 335.2,
	},
	8: {
		3: 48.1,
		4: 101.7,
		5: 183.7,
		6: 229.0,
		7: 571.0,
	},
	10: {
		3: 60.0,
		4: 127.8,
		5: 231.2,
		6: 376.7,
		7: 571.0,
	},
}

type sort2 struct {
	x []float64
	y []float64
}

func (s sort2) Len() int           { return len(s.x) }
func (s sort2) Less(i, j int) bool { return s.x[i] < s.x[j] }
func (s sort2) Swap(i, j int) {
	s.x[i], s.x[j] = s.x[j], s.x[i]
	s.y[i], s.y[j] = s.y[j], s.y[i]
}

// crank overwrites the entries in with their ranks
func crank(w []float64) float64 {
	j, ji, jt, n := 1, 0, 0, len(w)
	var rank float64

	var s float64

	for j < n {
		if w[j] != w[j-1] {
			w[j-1] = float64(j)
			j++
		} else {
			for jt = j + 1; jt <= n && w[jt-1] == w[j-1]; jt++ {
				// empty
			}
			rank = 0.5 * (float64(j) + float64(jt) - 1)
			for ji = j; ji <= (jt - 1); ji++ {
				w[ji-1] = rank
			}
			t := float64(jt - j)
			s += (t*t*t - t)
			j = jt
		}
	}
	if j == n {
		w[n-1] = float64(n)
	}
	return s
}

// Spearman returns the rank correlation coefficient between data1 and data2, and the associated p-value
func Spearman(data1, data2 []float64) (rs float64, p float64) {
	n := len(data1)
	wksp1, wksp2 := make([]float64, n), make([]float64, n)
	copy(wksp1, data1)
	copy(wksp2, data2)

	sort.Sort(sort2{wksp1, wksp2})
	sf := crank(wksp1)
	sort.Sort(sort2{wksp2, wksp1})
	sg := crank(wksp2)
	d := 0.0
	for j := 0; j < n; j++ {
		sq := wksp1[j] - wksp2[j]
		d += (sq * sq)
	}

	en := float64(n)
	en3n := en*en*en - en

	fac := (1.0 - sf/en3n) * (1.0 - sg/en3n)
	// без math.Sqrt(fac) работает аналогично ходашинской формуле
	rs = (1.0 - (6.0/en3n)*(d+(sf+sg)/12.0))

	if fac = (rs + 1.0) * (1.0 - rs); fac > 0 {
		t := rs * math.Sqrt((en-2.0)/fac)
		df := en - 2.0
		p = mathx.BetaInc(df/(df+t*t), 0.5*df, 0.5)
	}
	return rs, p
}

func sqr(x float64) float64 { return x * x }

func KendallW(mainMatrix [][]float64) (kendallw float64, s float64, kendallCrit float64) {
	numExp := len(mainMatrix)
	numAlt := len(mainMatrix[0])
	var sumCRank float64
	sumRank := make([]float64, numAlt)
	wksp1, wksp2 := make([]float64, numExp), make([]float64, numExp)
	for i := 0; i < numExp; i++ {
		copy(wksp1, mainMatrix[i])
		copy(wksp2, mainMatrix[i])
		sort.Sort(sort2{wksp1, wksp2})
		sf := crank(wksp1)
		//fmt.Println(mainMatrix[i])
		sumCRank += sf
	}

	rankMatrix := rankingTwoDimensional(mainMatrix, true)

	for i := 0; i < len(rankMatrix[0]); i++ {
		for j := 0; j < len(rankMatrix); j++ {
			sumRank[i] += rankMatrix[j][i]
		}
	}
	for i := 0; i < len(sumRank); i++ {
		s += math.Pow(sumRank[i]-0.5*float64(numExp)*(float64(numAlt)+1), 2)
	}
	//math.Pow(rankMatrix[i][j]-0.5*float64(numExp)*(float64(numAlt)+1), 2)
	kendallCrit = KenwallWMatrix[numExp][numAlt]
	down := math.Pow(float64(numExp), 2)*(math.Pow(float64(numAlt), 3)-float64(numAlt)) - float64(numExp)*sumCRank
	return 12 * s / down, s, kendallCrit
}

func rankingTwoDimensional(mainMatrix [][]float64, max10 bool) [][]float64 {
	n := len(mainMatrix)    //4
	m := len(mainMatrix[1]) //6
	if max10 {
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				mainMatrix[i][j] = math.Abs(mainMatrix[i][j] - 10)
			}
		}
	}
	rankMatrix := make([][]float64, n)
	for i := 0; i < n; i++ {
		rankMatrix[i] = make([]float64, m)
	}
	for i := range mainMatrix {
		for p := 0; p < len(mainMatrix[1]); p++ {
			r := 1
			s := 1
			for j := 0; j < len(mainMatrix[1]); j++ {
				if j != p && mainMatrix[i][j] < mainMatrix[i][p] {
					r += 1
				}
				if j != p && mainMatrix[i][p] == mainMatrix[i][j] {
					s += 1
				}
			}
			rankMatrix[i][p] = float64(r) + float64(s-1)/2.0
			//rankMatrix[i][p] = float64(r) + float64(s-1)/2.0 запасной вариант
		}
	}

	/* для одной строки
	for i := 0; i < len(mainMat); i++ {
		r := 1
		s := 1
		for j := 0; j < len(mainMat); j++ {
			if j != i && mainMat[j] < mainMat[i] {
				r += 1
			}
			if j != i && mainMat[j] == mainMat[i] {
				s += 1
			}
		}
		rankMatrix = append(rankMatrix, float64(r)+float64(s-1)/2.0)
	}

	*/
	return rankMatrix
}

func sumAlts(matrix [][]float64) []float64 {
	numCrit := len(matrix)
	numAlt := len(matrix[0])
	sumAltsArr := make([]float64, numAlt)
	for j := 0; j < numAlt; j++ {
		for i := 0; i < numCrit; i++ {
			sumAltsArr[j] += matrix[i][j]
		}
	}
	return sumAltsArr
}

func commonMatrix(mainMatrix [][]float64, max10 bool) [][][]int {
	n := len(mainMatrix)    //4
	m := len(mainMatrix[1]) //6

	compMatrix := make([][][]int, n)
	for i := 0; i < n; i++ {
		compMatrix[i] = make([][]int, m)
		for j := 0; j < m; j++ {
			compMatrix[i][j] = make([]int, m)
		}
	}

	for p := 0; p < n; p++ {
		//fmt.Println("----------")
		for i := 0; i < m; i++ {
			//fmt.Println("/////////")
			for j := 0; j < m; j++ {
				if mainMatrix[p][i] > mainMatrix[p][j] {
					if max10 {
						compMatrix[p][i][j] = -1
					} else {
						compMatrix[p][i][j] = 1
					}
					//fmt.Println(1)
				}
				if mainMatrix[p][i] < mainMatrix[p][j] {
					if max10 {
						compMatrix[p][i][j] = 1
					} else {
						compMatrix[p][i][j] = -1
					}
					//fmt.Println(-1)
				}
				if mainMatrix[p][i] == mainMatrix[p][j] {
					compMatrix[p][i][j] = 0
					//fmt.Println(0)
				}
			}
		}
	}
	return compMatrix
}

func pairComparison(commonMatrix [][][]int, numberExpert int) []float64 {
	//comparisonNumber := ((len(commonMatrix) * len(commonMatrix)) - len(commonMatrix)) / 2
	matrixNumber := len(commonMatrix)
	matrixRows := len(commonMatrix[0])
	matrixNumberInRows := len(commonMatrix[0][0])

	var (
		elemsPairComparison      []float64 //для сохранение элементов матриц после выполнения сравнения
		sum                      float64   //для сложения элементов матриц после сравнения для добавления в массив
		sum1                     float64   //для сложения строк суммированных элементов
		sumOfElemsPairComparison []float64 //массив из сумм строк элементов после сравнения
		k0                       int       //счётчик для получения суммированных строк
		k1                       int       //счётчик для получения итоговых значений попарных сравнений
		result                   []float64 //итоговые значения попарных сравнений
		//mean                     []float64 //медиана
	)

	for p := 0; p < matrixNumber; p++ {
		for q := p + 1; q < matrixNumber; q++ {
			for i := 0; i < matrixRows; i++ {
				for j := 0; j < matrixNumberInRows; j++ {
					k := math.Abs(float64(commonMatrix[p][i][j] - commonMatrix[q][i][j]))
					elemsPairComparison = append(elemsPairComparison, k)
				}
			}
		}
	}

	for j := 0; j < len(elemsPairComparison); j++ {
		k0++
		sum += elemsPairComparison[j]
		if k0 == matrixRows {
			k0 = 0
			sumOfElemsPairComparison = append(sumOfElemsPairComparison, sum)
			sum = 0
		}
	}

	for i := 0; i < len(sumOfElemsPairComparison); i++ {
		k1++
		sum1 += sumOfElemsPairComparison[i]
		if k1 == matrixRows {
			k1 = 0
			result = append(result, 0.5*sum1)
			sum1 = 0
		}
	}

	for p := 0; p < numberExpert; p++ {
		for i := 0; i < len(sumOfElemsPairComparison)-numberExpert-1; i++ {

		}
	}
	return result
}

func matrixSpearman(mainMatrix [][]float64) {
	for i := 0; i < len(mainMatrix); i++ {
		for j := i + 1; j < len(mainMatrix); j++ {
			value, pValue := Spearman(mainMatrix[i], mainMatrix[j])
			if pValue < 0.05 {
				fmt.Print("э", i+1, "/э", j+1, " ")
				fmt.Printf("%.4f", value)
				fmt.Println(" | отвергается гипотеза об отсутствии корр. связи")
			} else {
				fmt.Print("э", i+1, "/э", j+1, " ")
				fmt.Printf("%.4f", value)
				fmt.Println(" | подтверждается гипотеза об отсутствии корр. связи")
			}
		}
	}
}

func fullNormalized(optMatrix [][]float64) [][]float64 {
	numCrit := len(optMatrix)
	numAlt := len(optMatrix[0])
	normMatrix := make([][]float64, numCrit)
	for i := 0; i < numCrit; i++ {
		normMatrix[i] = make([]float64, numAlt)
	}
	for i := 0; i < numCrit; i++ {
		min, max := MinMax(optMatrix[i])
		for j := 0; j < numAlt; j++ {
			normMatrix[i][j] = (optMatrix[i][j] - min) / (max - min)
		}

	}
	return normMatrix
}

func MinMax(array []float64) (float64, float64) {
	var max float64 = array[0]
	var min float64 = array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return min, max
}

func idealPoint(normMatrix [][]float64) ([]float64, float64, int) {
	numCrit := len(normMatrix)
	numAlt := len(normMatrix[0])
	idealPointMatrix := make([][]float64, numCrit)
	for i := 0; i < numCrit; i++ {
		idealPointMatrix[i] = make([]float64, numAlt)
	}

	for i := 0; i < numCrit; i++ {
		_, max := MinMax(normMatrix[i])
		for j := 0; j < numAlt; j++ {
			idealPointMatrix[i][j] = sqr(normMatrix[i][j] - max)
		}

	}
	ipArr := make([]float64, numAlt)
	for j := 0; j < numAlt; j++ {
		for i := 0; i < numCrit; i++ {
			//fmt.Println(idealPointMatrix[i])
			ipArr[j] += idealPointMatrix[i][j]
		}
	}
	ip, _ := MinMax(ipArr)
	alt := 0
	for i, v := range ipArr {
		if v == ip {
			alt = i + 1
		}
	}
	return ipArr, ip, alt
}

func antiIdealPoint(normMatrix [][]float64) ([]float64, float64, int) {
	numCrit := len(normMatrix)
	numAlt := len(normMatrix[0])
	idealPointMatrix := make([][]float64, numCrit)
	for i := 0; i < numCrit; i++ {
		idealPointMatrix[i] = make([]float64, numAlt)
	}

	for i := 0; i < numCrit; i++ {
		min, _ := MinMax(normMatrix[i])
		for j := 0; j < numAlt; j++ {
			idealPointMatrix[i][j] = sqr(normMatrix[i][j] - min)
		}

	}
	ipArr := make([]float64, numAlt)
	for j := 0; j < numAlt; j++ {
		for i := 0; i < numCrit; i++ {
			//fmt.Println(idealPointMatrix[i])
			ipArr[j] += idealPointMatrix[i][j]
		}
	}
	_, ip := MinMax(ipArr)
	alt := 0
	for i, v := range ipArr {
		if v == ip {
			alt = i + 1
		}
	}
	return ipArr, ip, alt
}

func EvlanovKutuzov(rankMatrix [][]float64) []float64 {
	const e = 0.001 //условие останова
	numCrit := len(rankMatrix)
	numAlt := len(rankMatrix[0])
	sumRanks := sumAlts(rankMatrix)

	avgRanks := make([]float64, len(sumRanks)) //массив для хранения средних оценок

	normCoefArr := make([]float64, len(sumRanks)) //массив для хранения нормированных коэффициентов

	normCoef := 0.0 //нормировочный коэффициент

	koefCompMatrix := make([][]float64, numCrit) //матрица для получения критериев
	for i := 0; i < numCrit; i++ {
		koefCompMatrix[i] = make([]float64, numAlt)
	}

	kArr := make([]float64, numCrit) //массив критериев

	avgCompMatrix := make([][]float64, numCrit) //ранги умноженные на кртиерии
	for i := 0; i < numCrit; i++ {
		avgCompMatrix[i] = make([]float64, numAlt)
	}

	var end = true //условие останова

	avgRanks1 := make([]float64, len(sumRanks)) //второй массив средних оценок для проверки останова

	for end == true {
		for i, v := range sumRanks {
			avgRanks[i] = v / float64(len(rankMatrix))
			normCoefArr[i] = avgRanks[i] * v
			normCoef += normCoefArr[i]
		}
		//fmt.Println(normCoef)
		for i := 0; i < len(rankMatrix); i++ {
			for j := 0; j < len(rankMatrix[0]); j++ {
				koefCompMatrix[i][j] = rankMatrix[i][j] * avgRanks[j]
				kArr[i] += koefCompMatrix[i][j]
			}
			kArr[i] = kArr[i] * (1 / normCoef)
		}

		for i := 0; i < numCrit; i++ {
			for j := 0; j < numAlt; j++ {
				avgCompMatrix[i][j] = rankMatrix[i][j] * kArr[i]
				avgRanks1[j] += avgCompMatrix[i][j]
			}
		}

		for i := range avgRanks {
			if e < math.Abs(avgRanks[i]-avgRanks1[i]) {
				end = false
			}
		}
	}
	return kArr
}

func LessSquare(x, y []float64) ([]float64, []float64, float64) {
	f := func(a, x, b float64) float64 { return a + x*b }
	yNew := make([]float64, len(y))
	predict := make([]float64, 3)
	n := float64(len(y))
	//var a, b float64
	var sumXY, sumX, sumY, sumXSqr float64
	for i := 0; i < len(y); i++ {
		sumXY += x[i] * y[i]
		sumY += y[i]
		sumX += x[i]
		sumXSqr += x[i] * x[i]
	}
	matrix := [][]float64{
		{n, sumX, sumY},
		{sumX, sumXSqr, sumXY},
	}
	ans := GaussMehtodforSqr(matrix)
	for i := 0; i < len(y); i++ {
		yNew[i] = f(ans[0], x[i], ans[1])
	}
	k := 1
	for i := 0; i < 3; i++ {
		predict[i] = f(ans[0], float64(len(x)+k), ans[1])
		k++
	}
	s := 0.0
	for i := 0; i < len(y); i++ {
		s += math.Pow(y[i]-yNew[i], 2)
	}
	return yNew, predict, math.Sqrt(s/n - 2)
}

func SmoothMean3(x, y []float64) []float64 {
	yNew := make([]float64, len(y)-2)
	xNew := make([]float64, len(x)-2)
	m := make([]float64, len(x)-2)
	for i := 0; i < len(xNew); i++ {
		xNew[i] = x[i] + 1
		m[i] = (y[i] + y[i+1] + y[i+2]) / 3
	}
	for i := 0; i < len(yNew); i++ {
		yNew[i] = m[i] + 1/3*(y[i+1]-y[i])
	}
	return yNew
}

func ExponentialSmooth(x, y []float64, a float64) ([]float64, []float64, float64) {

	n := float64(len(y))
	var sumXY, sumX, sumY, sumXSqr float64

	for i := 0; i < len(y); i++ {
		sumXY += x[i] * y[i]
		sumY += y[i]
		sumX += x[i]
		sumXSqr += x[i] * x[i]
	}
	matrix := [][]float64{
		{n, sumX, sumY},
		{sumX, sumXSqr, sumXY},
	}
	ans := GaussMehtodforSqr(matrix)

	s1 := make([]float64, len(y))
	s2 := make([]float64, len(y))
	s10 := ans[0] - ((1-a)/a)*ans[1]
	s20 := ans[0] - (2*(1-a)/a)*ans[1]
	s1[0] = a*y[0] + (1-a)*s10
	s2[0] = a*s1[0] + (1-a)*s20
	for i := 1; i < len(s1); i++ {
		s1[i] = a*y[i] + (1-a)*s1[i-1]
		s2[i] = a*s1[i] + (1-a)*s2[i-1]
	}

	f := func(a, x, b float64) float64 { return a + x*b }
	predict := make([]float64, 3)
	yNew := make([]float64, len(y))
	for i := 0; i < len(y); i++ {
		a0New := 2*s1[i] - s2[i]
		a1New := (a / (1 - a)) * (s1[i] - s2[i])
		yNew[i] = a0New + a1New
		//fmt.Println(a0New, a1New)
	}
	k := 2
	for i := 0; i < len(predict); i++ {
		a0New := 2*s1[len(y)-1] - s2[len(y)-1]
		a1New := (a / (1 - a)) * (s1[len(y)-1] - s2[len(y)-1])
		predict[i] = f(a0New, float64(k), a1New)
		k++
	}
	q0 := 0.0
	for i := 0; i < len(y); i++ {
		q0 += math.Pow(y[i]-yNew[i], 2)
	}
	q := math.Sqrt(q0/(n-2)) * math.Sqrt((a/(math.Pow(2-a, 3)))*(1+4*(1-a)+5*((1-a)*(1-a))+2*a*(4-3*a)*3+2*a*a*3*3))
	fmt.Println(q)
	return yNew, predict, q
}

func LessSquareSqr(x, y []float64) ([]float64, []float64, float64) {
	n := float64(len(y))
	yNew := make([]float64, len(y))
	predict := make([]float64, 3)
	var (
		sumX   float64
		sumX2  float64
		sumX3  float64
		sumX4  float64
		sumY   float64
		sumXY  float64
		sumX2Y float64
	)
	for i := 0; i < len(y); i++ {
		sumX += x[i]
		sumX2 += x[i] * x[i]
		sumX3 += x[i] * x[i] * x[i]
		sumX4 += x[i] * x[i] * x[i] * x[i]
		sumY += y[i]
		sumXY += y[i] * x[i]
		sumX2Y += y[i] * x[i] * x[i]
	}
	matrix := [][]float64{
		{sumX4, sumX3, sumX2, sumX2Y},
		{sumX3, sumX2, sumX, sumXY},
		{sumX2, sumX, n, sumY},
	}
	a := GaussMehtodforSqr(matrix)
	f := func(a, x, b, c float64) float64 { return a*x*x + b*x + c }
	for i := 0; i < len(y); i++ {
		yNew[i] = f(a[0], x[i], a[1], a[2])
	}
	k := 1
	for i := 0; i < 3; i++ {
		predict[i] = f(a[0], float64(len(x)+k), a[1], a[2])
		k++
	}
	s := 0.0
	for i := 0; i < len(y); i++ {
		s += math.Pow(y[i]-yNew[i], 2)
	}
	return yNew, predict, math.Sqrt(s/n - 3)
}

func GaussMehtodforSqr(matrix [][]float64) []float64 {
	n := len(matrix)
	answer := make([]float64, n)
	var (
		//k   int
		tmp float64
	)
	for i := 0; i < n; i++ {
		tmp = matrix[i][i]
		for j := n; j >= i; j-- {
			matrix[i][j] /= tmp
		}
		for j := i + 1; j < n; j++ {
			tmp = matrix[j][i]
			for k := n; k >= i; k-- {
				matrix[j][k] -= tmp * matrix[i][k]
			}
		}
	}
	answer[n-1] = matrix[n-1][n]
	for i := n - 2; i >= 0; i-- {
		answer[i] = matrix[i][n]
		for j := i + 1; j < n; j++ {
			answer[i] -= matrix[i][j] * answer[j]
		}
	}
	return answer
}

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

func main() {

	/*
		mainMatrix := [][]float64{
			{2, 1, 3, 4, 4, 5},
			{1, 2, 2, 3, 4, 5},
			{1, 1, 1, 3, 2, 4},
			{1, 1, 2, 3, 3, 4},
			//{3, 1, 3, 4, 3, 2},
		}
	*/

	crtMatrix1 := [][]float64{
		{8, 10, 5, 2, 3, 8},
		{9, 10, 5, 3, 5, 9},
		{8, 4, 7, 3, 5, 1},
		{9, 10, 7, 5, 6, 9},
	}
	crtMatrix2 := [][]float64{
		{5, 7, 5, 6, 5, 6},
		{7, 8, 7, 2, 7, 9},
		{7, 3, 9, 6, 5, 3},
		{6, 8, 8, 9, 8, 7},
	}
	crtMatrix3 := [][]float64{
		{5, 10, 2, 2, 5, 8},
		{5, 10, 2, 2, 3, 2},
		{5, 2, 4, 7, 9, 2},
		{10, 9, 6, 6, 7, 5},
	}

	fullMatrix := [][][]float64{crtMatrix1, crtMatrix2, crtMatrix3}
	/*
		optMatrix := [][]float64{
			{175, 190, 150, 160, 120, 111},
			{20, 10, 5, 14, 12, 9},
			{1, 3, 3, 4, 5, 6},
		}

	*/

	for p := range fullMatrix {
		mainMatrix := fullMatrix[p]
		fmt.Println("Критерий №", p+1)
		resultRanking := rankingTwoDimensional(mainMatrix, true)
		//если 10 - максимальная оценка, то в функцию commonMatrix нужно отправить true
		result := commonMatrix(mainMatrix, true)
		fmt.Print("Ранжировка (1 критерий):\n")
		for i := range result {
			fmt.Println("Эксперт №", i+1, ":", resultRanking[i])
		}
		fmt.Println()
		sumAlt := sumAlts(resultRanking)
		for i, v := range sumAlt {
			fmt.Println("Сумма ранга по альтернативе №", i+1, " = ", v)
		}
		fmt.Print("\nОбобщённые ранжировки:\n")
		for i := range result {
			fmt.Println("Эксперт №", i+1, ":", result[i])
		}
		fmt.Println("\nЗначения для расчёта медианы: ", pairComparison(result, 4), "\n")
		fmt.Println("Итоговая компетентность экспертов: ")
		koefComp := EvlanovKutuzov(resultRanking)
		_, maxCompExp := MinMax(koefComp)
		for i, v := range koefComp {
			fmt.Print("Эксперт №", i+1, " : ", v)
			if maxCompExp == v {
				fmt.Print(" (наиболее компетентный эксперт)")
			}
			fmt.Println()
		}
		fmt.Print("\nГипотезы наличия корреляционной связи: \n")
		matrixSpearman(mainMatrix)
		kendalW, s, kendallCrit := KendallW(mainMatrix)
		fmt.Println("\nОценка согласованности экспертов:\nW: ", kendalW, "\nS: ", s, "\nКритическое значение: ", kendallCrit)
		if s > kendallCrit {
			fmt.Println("Мнения экспертов согласованны (S > критического значения)\n")
		} else {
			fmt.Println("Мнения экспертов несогласованны (S < критического значения)\n")
		}
	}

	x := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	y := []float64{490, 477, 430, 399, 383, 357, 278, 205, 182, 214, 206, 192, 186, 188, 190}

	yNew, predict, errLsq := LessSquare(x, y)
	fmt.Println("Прогноз на следующие 3 периода (линейная модель МНК): ")
	for i, v := range predict {
		fmt.Print("Прогнозный период ", i+1, " : ")
		fmt.Printf("%.3f", v)
		fmt.Println()
	}
	fmt.Println("Стандартная ошибка :", errLsq)
	fmt.Println()
	yNewMean := SmoothMean3(x, y)
	yNewSqr, predictSqr, errLsq2 := LessSquareSqr(x, y)
	fmt.Println("Прогноз на следующие 3 периода (квадратичная модель МНК): ")
	for i, v := range predictSqr {
		fmt.Print("Прогнозный период ", i+1, " : ")
		fmt.Printf("%.3f", v)
		fmt.Println()
	}
	fmt.Println("Стандартная ошибка :", errLsq2)
	fmt.Println()
	fmt.Println("Прогноз на следующие 3 периода (экспоненциальная линейная модель): ")
	yNewExp, predictExp, err := ExponentialSmooth(x, y, 0.7)
	//predictExp := ExpPredict(x, y, 0.7)
	for i, v := range predictExp {
		fmt.Print("Прогнозный период ", i+1, " : ")
		fmt.Printf("%.3f", v)
		fmt.Println()
	}
	fmt.Println("Ошибка прогноза: ", err)
	DrawLines(x, y, yNew, "bar.html", "lsq")
	DrawLines(x, y, yNewMean, "bar1.html", "mean3")
	DrawLines(x, y, yNewSqr, "bar2.html", "lsq")
	DrawLines(x, y, yNewExp, "bar3.html", "exp")
	//fmt.Println(yNewExp)
	/*
		normMatrix := fullNormalized(optMatrix)
		fmt.Print("Нормализированные значения:\n")
		for i := 0; i < len(normMatrix); i++ {
			fmt.Printf("%.2f", normMatrix[i])
			fmt.Print("\n")
		}
		ipArr, ip, i := idealPoint(normMatrix)
		fmt.Println("\nРасстояние до идеальной точки: ", ip, "Альтернатива: ", i)
		fmt.Print("Матрица всех расстояний ид.т: ")
		fmt.Printf("%.2f", ipArr)
		fmt.Print("\n")
		aipArr, aip, ai := antiIdealPoint(normMatrix)
		fmt.Println("Расстояние до антиидеальной точки: ", aip, "Альтернатива: ", ai)
		fmt.Print("Матрица всех расстояний а.ид.т: ")
		fmt.Printf("%.2f", aipArr)
		fmt.Print("\n")
		//fmt.Println(EvlanovKutuzov(resultRanking))

	*/

}
