package expert_opinions

import (
	"fmt"
	"github.com/aclements/go-moremath/mathx"
	"math"
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

	rankMatrix := RankingTwoDimensional(mainMatrix, true)

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

func RankingTwoDimensional(mainMatrix [][]float64, max10 bool) [][]float64 {
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

func SumAlts(matrix [][]float64) []float64 {
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

func CommonMatrix(mainMatrix [][]float64, max10 bool) [][][]int {
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

func PairComparison(commonMatrix [][][]int, numberExpert int) []float64 {
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

func MatrixSpearman(mainMatrix [][]float64) {
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

func EvlanovKutuzov(rankMatrix [][]float64) []float64 {
	const e = 0.001 //условие останова
	numCrit := len(rankMatrix)
	numAlt := len(rankMatrix[0])
	sumRanks := SumAlts(rankMatrix)

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
		k := 0
		for i := range avgRanks {
			if e < math.Abs(avgRanks[i]-avgRanks1[i]) {
				k++
				if k == len(avgRanks) {
					end = false
				}
			}
		}
	}

	return kArr
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
