package main

import (
	"ffmraz/decisions"
	decisionsAp "ffmraz/decisionsAp"
	"ffmraz/fuzzySets"

	//"ffmraz/decisionsAp"
	expert "ffmraz/expert_opinions"
	"ffmraz/fuzzy-logic"
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
		{1.90, 2.25, 1.95, 2.5, 2.14, 2.5},
		{17, 16, 20, 9, 12, 16},
		{4.2, 4.4, 4.8, 4.6, 4.2, 4.8},
	}

	w := []float64{0.5, 0.3, 0.2}
	options := []float64{1}
	normMatrix := decisions.FullNormalized(z1, options)
	fmt.Print("Нормализированные значения:\n")
	for i := 0; i < len(normMatrix); i++ {
		fmt.Printf("%.4f", normMatrix[i])
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

	p := []float64{0.4, 0.45, 0.15}
	zz1 := [][]float64{
		{17, 20, 21},
		{16, 19, 23},
		{20, 23, 23},
		{9, 12, 15},
		{12, 15, 16},
		{16, 19, 25},
	}
	zz2 := [][]float64{
		{1.9, 1.85, 1.8},
		{2.25, 2.23, 2.2},
		{1.95, 1.90, 1.89},
		{2.5, 2.47, 2.4},
		{2.14, 2.1, 2},
		{2.5, 2.3, 2.1},
	}
	/*
		zTest := [][]float64{
			{2, 4, 7},
			{1, 3, 6},
			{2, 3, 5},
			{2, 4, 6},
			{1, 5, 7},
		}

		fmt.Println("-----", decisionsAp.NormalizedMatrix(zTest, "min"))

	*/
	zz1Norm := decisionsAp.NormalizedMatrix(zz1, "min")
	zz2Norm := decisionsAp.NormalizedMatrix(zz2, "max")

	b1 := decisionsAp.Bl(zz1Norm, p)
	s1 := decisionsAp.Std(zz1Norm, p)
	c1 := decisionsAp.ConvolutionStdBl(b1, s1)

	b2 := decisionsAp.Bl(zz2Norm, p)
	s2 := decisionsAp.Std(zz2Norm, p)
	c2 := decisionsAp.ConvolutionStdBl(b2, s2)

	fmt.Println("\n===ПеРвАя СиТуАцИя====")
	fmt.Println("\nСнятие неопределенности для первого критерия (свертка B и СКО)")
	for r := 0; r < len(c1); r++ {
		fmt.Printf("%.4f", c1[r])
		fmt.Println()
	}

	fmt.Println()

	fmt.Println("Снятие неопределенности для второго критерия (свертка B и СКО)")
	for r := 0; r < len(c1); r++ {
		fmt.Printf("%.4f", c2[r])
		fmt.Println()
	}

	cArr := [][][]float64{c1, c2} //для однокритериального оставьте только одно значение в зависимости от того, какой критерий смотрите. то же самое с другими переменными *Arr

	ipAp := decisionsAp.IdealPointAp(cArr)
	reAp := decisionsAp.RealteСoncessionAp(cArr)
	abAp := decisionsAp.AbsoluteСoncessionAp(cArr)
	aipAp := decisionsAp.AntiIdealPointAp(cArr)
	/*
		fmt.Println(ipAp[0])
		fmt.Println(reAp[0])
		fmt.Println(abAp[0])
		fmt.Println(aipAp[0])

	*/

	fmt.Println("\nМногокритериальная оценка ид.т")
	for r := 0; r < len(ipAp); r++ {
		fmt.Printf("%.4f", ipAp[r])
		fmt.Println()
	}

	fmt.Println()

	fmt.Println("\nМногокритериальная оценка а.ид.т")
	for r := 0; r < len(aipAp); r++ {
		fmt.Printf("%.4f", aipAp[r])
		fmt.Println()
	}

	fmt.Println()

	fmt.Println("\nМногокритериальная оценка относ. уст.")
	for r := 0; r < len(reAp); r++ {
		fmt.Printf("%.4f", reAp[r])
		fmt.Println()
	}

	fmt.Println()

	fmt.Println("\nМногокритериальная оценка абс. устр")
	for r := 0; r < len(abAp); r++ {
		fmt.Printf("%.4f", abAp[r])
		fmt.Println()
	}

	fmt.Println()

	//вторая ситуация
	z10_1 := decisionsAp.Z10(zz1Norm) //гурвиц
	z10_2 := decisionsAp.Z10(zz2Norm)

	fmt.Println("\n===ВтОрАя СиТуАцИя====")
	fmt.Println("\nСнятие неопределенности для первого критерия (Гурвиц)")
	for r := 0; r < len(z10_1); r++ {
		fmt.Printf("%.4f", z10_1[r])
		fmt.Println()
	}

	fmt.Println()

	fmt.Println("\nСнятие неопределенности для второго критерия (Гурвиц)")
	for r := 0; r < len(z10_2); r++ {
		fmt.Printf("%.4f", z10_2[r])
		fmt.Println()
	}

	z10Arr := [][][]float64{z10_1, z10_2}

	z10pAp := decisionsAp.IdealPointAp(z10Arr)
	z10aiAp := decisionsAp.AntiIdealPointAp(z10Arr)
	z10abAp := decisionsAp.AbsoluteСoncessionAp(z10Arr)
	z10ReAp := decisionsAp.RealteСoncessionAp(z10Arr)

	fmt.Println("\nМногокритериальная оценка ид.т")
	for r := 0; r < len(z10pAp); r++ {
		fmt.Printf("%.4f", z10pAp[r])
		fmt.Println()
	}

	fmt.Println()

	fmt.Println("\nМногокритериальная оценка а.ид.т")
	for r := 0; r < len(z10aiAp); r++ {
		fmt.Printf("%.4f", z10aiAp[r])
		fmt.Println()
	}

	fmt.Println()

	fmt.Println("\nМногокритериальная оценка абс. уст.")
	for r := 0; r < len(z10abAp); r++ {
		fmt.Printf("%.4f", z10abAp[r])
		fmt.Println()
	}

	fmt.Println()

	fmt.Println("\nМногокритериальная оценка относ . устр")
	for r := 0; r < len(z10ReAp); r++ {
		fmt.Printf("%.4f", z10ReAp[r])
		fmt.Println()
	}

	fmt.Println()
	fmt.Println("\n===ТрЕтЬя СиТуАцИя====")
	//третья ситуация
	betta := []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1}
	for b := 0; b < len(betta); b++ {
		zI_1 := decisionsAp.ZI(c1, z10_1, betta[b]) //свертка свертки байеса+ско и гурвица
		zI_2 := decisionsAp.ZI(c2, z10_2, betta[b])
		fmt.Println("---------Бетта = ", betta[b], "---------------")
		fmt.Println("\nСнятие неопределенности для первого критерия (свертка свертки B и СКО и Гурвица)")
		for r := 0; r < len(zI_1); r++ {
			fmt.Printf("%.4f", zI_1[r])
			fmt.Println()
		}

		fmt.Println()

		fmt.Println("\nСнятие неопределенности для второго критерия (свертка свертки B и СКО и Гурвица)")
		for r := 0; r < len(zI_2); r++ {
			fmt.Printf("%.4f", zI_2[r])
			fmt.Println()
		}

		zIArr := [][][]float64{zI_1, zI_2}

		zIpAp := decisionsAp.IdealPointAp(zIArr)
		zIaiAp := decisionsAp.AntiIdealPointAp(zIArr)
		zIabAp := decisionsAp.AbsoluteСoncessionAp(zIArr)
		zIReAp := decisionsAp.RealteСoncessionAp(zIArr)

		fmt.Println("\nМногокритериальная оценка ид.т")
		for r := 0; r < len(zIpAp); r++ {
			fmt.Printf("%.4f", zIpAp[r])
			fmt.Println()
		}

		fmt.Println()

		fmt.Println("\nМногокритериальная оценка а.ид.т")
		for r := 0; r < len(zIaiAp); r++ {
			fmt.Printf("%.4f", zIaiAp[r])
			fmt.Println()
		}

		fmt.Println()

		fmt.Println("\nМногокритериальная оценка абс. уст.")
		for r := 0; r < len(zIabAp); r++ {
			fmt.Printf("%.4f", zIabAp[r])
			fmt.Println()
		}

		fmt.Println()

		fmt.Println("\nМногокритериальная оценка относ. уст.")
		for r := 0; r < len(zIReAp); r++ {
			fmt.Printf("%.4f", zIReAp[r])
			fmt.Println()
		}

		fmt.Println()
	}

	//нечеткие множества

	abc := [][]float64{
		{0, 0.3, 0.7, 0.9, 1},
		{1, 0.7, 0.3, 0.1, 0},
		{0, 0.1, 0.3, 0.7, 1},
		{0, 0.01, 0.09, 0.49, 1},
	}

	_, cnvRm := fuzzySets.Rm(abc)
	_, cnvRa := fuzzySets.Ra(abc)
	_, cnvRb := fuzzySets.Rb(abc)
	_, cnvRgg := fuzzySets.Rgg(abc)
	_, cnvRss := fuzzySets.Rss(abc)
	_, cnvRsg := fuzzySets.Rsg(abc)
	_, cnvRgs := fuzzySets.Rgs(abc)
	fmt.Println(cnvRm)
	fmt.Println(cnvRa)
	fmt.Println(cnvRb)
	fmt.Println(cnvRgg)
	fmt.Println(cnvRss)
	fmt.Println(cnvRsg)
	fmt.Println(cnvRgs)

	rng := [][]float64{
		{2.53, 1.02, 1.84, 0.11, 1.43, 0.40},
		{0.07, 0.58, 1.01, 2.44, 0.69, 2.13},
		{0.00, 0.07, 0.00, 0.67, 0.00, 0.36},
		{0.27, 1.28, 1.08, 2.67, 1.13, 2.36},
	}
	fmt.Println(expert.RankingTwoDimensional(rng, true))

	fuzzyAlts := [][]float64{
		{0.40, 0.20, 0.70, 0.70, 0.30, 0.90, 0.10, 0.30, 0.60, 0.40},
		{0.20, 0.10, 0.90, 0.30, 0.10, 0.90, 0.20, 0.40, 0.70, 0.80},
		{0.10, 0.30, 0.60, 0.30, 0.60, 0.20, 0.30, 0.60, 0.80, 0.80},
	}

	fuzzyWeights := [][]float64{
		{0, 0, 0.33},    //bad
		{0, 0.33, 0.66}, //medium
		{0.66, 1, 1},    //high
	}

	fuzzBad := fuzzy_logic.Fuzzification(fuzzyAlts, fuzzyWeights[0])
	fuzzMed := fuzzy_logic.Fuzzification(fuzzyAlts, fuzzyWeights[1])
	fuzzHigh := fuzzy_logic.Fuzzification(fuzzyAlts, fuzzyWeights[2])
	rules := fuzzy_logic.RulesConv(fuzzBad, fuzzMed, fuzzHigh, 0)
	fmt.Println(rules)

}
