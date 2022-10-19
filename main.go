package main

import (
	"ffmraz/decisions"
	decisionsAp "ffmraz/decisionsAp"
	//"ffmraz/decisionsAp"
	expert "ffmraz/expert_opinions"
	plot "ffmraz/plots"
	predict "ffmraz/predictions"

	//decisions "ffmraz/decisions"
	"fmt"
)

func main() {
	/*
		mainMatrix1 := [][]float64{
			{2, 1, 3, 4, 4, 5},
			{1, 2, 2, 3, 4, 5},
			{1, 1, 1, 3, 2, 4},
			{1, 1, 2, 3, 3, 4},
			//{3, 1, 3, 4, 3, 2},
		}
		r := expert.RankingTwoDimensional(mainMatrix1, false)
		fmt.Println(expert.EvlanovKutuzov(r))

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

	for p := range fullMatrix {
		mainMatrix := fullMatrix[p]
		fmt.Println("Критерий №", p+1)
		resultRanking := expert.RankingTwoDimensional(mainMatrix, true)
		//если 10 - максимальная оценка, то в функцию commonMatrix нужно отправить true
		result := expert.CommonMatrix(mainMatrix, true)
		fmt.Print("Ранжировка (1 критерий):\n")
		for i := range result {
			fmt.Println("Эксперт №", i+1, ":", resultRanking[i])
		}
		fmt.Println()
		sumAlt := expert.SumAlts(resultRanking)
		for i, v := range sumAlt {
			fmt.Println("Сумма ранга по альтернативе №", i+1, " = ", v)
		}
		fmt.Print("\nОбобщённые ранжировки:\n")
		for i := range result {
			fmt.Println("Эксперт №", i+1, ":", result[i])
		}
		fmt.Println("\nЗначения для расчёта медианы: ", expert.PairComparison(result, 4), "\n")
		fmt.Println("Итоговая компетентность экспертов: ")
		koefComp := expert.EvlanovKutuzov(resultRanking)
		_, maxCompExp := expert.MinMax(koefComp)
		for i, v := range koefComp {
			fmt.Print("Эксперт №", i+1, " : ", v)
			if maxCompExp == v {
				fmt.Print(" (наиболее компетентный эксперт)")
			}
			fmt.Println()
		}
		fmt.Print("\nГипотезы наличия корреляционной связи: \n")
		expert.MatrixSpearman(mainMatrix)
		kendalW, s, kendallCrit := expert.KendallW(mainMatrix)
		fmt.Println("\nОценка согласованности экспертов:\nW: ", kendalW, "\nS: ", s, "\nКритическое значение: ", kendallCrit)
		if s > kendallCrit {
			fmt.Println("Мнения экспертов согласованны (S > критического значения)\n")
		} else {
			fmt.Println("Мнения экспертов несогласованны (S < критического значения)\n")
		}
	}

	x := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	y := []float64{490, 477, 430, 399, 383, 357, 278, 205, 182, 214, 206, 192, 186, 188, 190}

	yNew, predict1 := predict.LessSquare(x, y)
	fmt.Println("Прогноз на следующие 3 периода (линейная модель МНК): ")
	for i, v := range predict1 {
		fmt.Print("Прогнозный период ", i+1, " : ")
		fmt.Printf("%.3f", v)
		fmt.Println()
	}
	fmt.Println()
	yNewMean := predict.SmoothMean3(x, y)
	yNewSqr, predictSqr := predict.LessSquareSqr(x, y)
	fmt.Println("Прогноз на следующие 3 периода (квадратичная модель МНК): ")
	for i, v := range predictSqr {
		fmt.Print("Прогнозный период ", i+1, " : ")
		fmt.Printf("%.3f", v)
		fmt.Println()
	}
	fmt.Println()
	fmt.Println("Прогноз на следующие 3 периода (экспоненциальная линейная модель): ")
	yNewExp, predictExp := predict.ExponentialSmooth(x, y, 0.7)
	//predictExp := ExpPredict(x, y, 0.7)
	for i, v := range predictExp {
		fmt.Print("Прогнозный период ", i+1, " : ")
		fmt.Printf("%.2f", v)
		fmt.Println()
	}
	fmt.Println()
	fmt.Println("Прогноз на следующие 3 периода (экспоненциальная квадратичная модель): ")
	yNewExpSqr, predictExpSqr := predict.ExponentialSmoothSqr(x, y, 0.7)
	for i, v := range predictExpSqr {
		fmt.Print("Прогнозный период ", i+1, " : ")
		fmt.Printf("%.2f", v)
		fmt.Println()
	}
	plot.DrawLines(x, y, yNew, "bar.html", "lsq")
	plot.DrawLines(x, y, yNewMean, "bar1.html", "mean3")
	plot.DrawLines(x, y, yNewSqr, "bar2.html", "lsq")
	plot.DrawLines(x, y, yNewExp, "bar3.html", "exp")
	plot.DrawLines(x, y, yNewExpSqr, "bar4.html", "exp")

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

	fmt.Println()
	fmt.Println("Ошибки для линейного мнк: ")
	errsLSQ := predict.Errors(y, yNew, x, 2)
	for k, v := range errsLSQ {
		fmt.Println(k, ": ", v)
	}
	fmt.Println()
	fmt.Println("Ошибки для квадратичного мнк: ")
	errsLSQSqr := predict.Errors(y, yNewSqr, x, 3)
	for k, v := range errsLSQSqr {
		fmt.Println(k, ": ", v)
	}
	fmt.Println()
	fmt.Println("Ошибки для линейной экспоненты: ")
	errsExp := predict.Errors(y, yNewExp, x, 2)
	for k, v := range errsExp {
		fmt.Println(k, ": ", v)
	}
	fmt.Println()
	fmt.Println("Ошибки для квадратичной экспоненты: ")
	errsExpSqr := predict.Errors(y, yNewExpSqr, x, 3)
	for k, v := range errsExpSqr {
		fmt.Println(k, ": ", v)
	}

	z1 := [][]float64{
		{1, 0.541566065, 0.061842549, 0.870458937, 0},
		{0.728254671, 0.582059749, 0.024211536, 0, 1},
		{1, 0.511876922, 0.051906008, 0.760278934, 0},
	}

	w := []float64{1, 1, 1}

	normMatrix := decisions.FullNormalized(z1)
	fmt.Print("Нормализированные значения:\n")
	for i := 0; i < len(normMatrix); i++ {
		fmt.Printf("%.2f", normMatrix[i])
		fmt.Print("\n")
	}
	ipArr, ip, i := decisions.IdealPoint(normMatrix, w)
	fmt.Println("\nРасстояние до идеальной точки: ", ip, ", Альтернатива: ", i)
	fmt.Print("Матрица всех расстояний ид.т: ")
	fmt.Printf("%.2f", ipArr)
	fmt.Print("\n")
	aipArr, aip, ai := decisions.AntiIdealPoint(normMatrix, w)
	fmt.Println("Расстояние до антиидеальной точки: ", aip, ", Альтернатива: ", ai)
	fmt.Print("Матрица всех расстояний а.ид.т: ")
	fmt.Printf("%.2f", aipArr)
	fmt.Print("\n")
	abMatr, abp, ab := decisions.AbsoluteСoncession(normMatrix, w)
	fmt.Println("Принцип абсолютной уступки: ", abp, ", Альтернатива: ", ab)
	fmt.Print("Матрица: ")
	fmt.Printf("%.2f", abMatr)
	fmt.Print("\n")
	reMatr, rebp, rb := decisions.RelateСoncession(normMatrix, w)
	fmt.Println("Принцип относительной уступки: ", rebp, ", Альтернатива: ", rb)
	fmt.Print("Матрица: ")
	fmt.Printf("%.2f", reMatr)
	fmt.Print("\n")

	p := []float64{0.3, 0.15, 0.55}
	zz := [][]float64{
		{0.8, 0.7, 0.8},
		{0.9, 0.5, 0.8},
		{0.7, 0.6, 0.7},
		{0.9, 0.8, 0.8},
	}
	zz2 := [][]float64{
		{0.7, 0.7, 0.6},
		{0.8, 0.7, 0.8},
		{0.7, 0.7, 0.7},
		{0.9, 0.8, 0.9},
	}
	zz3 := [][]float64{
		{0.8, 0.7, 0.8},
		{0.6, 0.5, 0.6},
		{0.7, 0.6, 0.7},
		{0.5, 0.4, 0.5},
	}
	zz4 := [][]float64{
		{0.6, 0.4, 0.5},
		{0.9, 0.7, 0.9},
		{0.6, 0.5, 0.6},
		{0.8, 0.7, 0.8},
	}

	b1 := decisionsAp.Bl(zz, p)
	s1 := decisionsAp.Std(zz, p)
	c := decisionsAp.ConvolutionStdBl(b1, s1)

	b2 := decisionsAp.Bl(zz2, p)
	s2 := decisionsAp.Std(zz2, p)
	c2 := decisionsAp.ConvolutionStdBl(b2, s2)

	b3 := decisionsAp.Bl(zz3, p)
	s3 := decisionsAp.Std(zz3, p)
	c3 := decisionsAp.ConvolutionStdBl(b3, s3)

	b4 := decisionsAp.Bl(zz4, p)
	s4 := decisionsAp.Std(zz4, p)
	c4 := decisionsAp.ConvolutionStdBl(b4, s4)

	cArr := [][][]float64{c, c2, c3, c4}

	decisionsAp.IdealPointAp(cArr)

	//cRow := decisionsAp.RowToCol(c)

	//maxArr := decisionsAp.MaxInConvolution(cRow)

	//fmt.Println(maxArr)

}
