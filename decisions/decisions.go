package decisions

import (
	"math"
)

type Point struct {
	X float64
	Y float64
	Z float64
}

func ToPoint(z [][]float64) []Point {
	points := make([]Point, len(z[0]))

	for i := 0; i < len(z[0]); i++ {
		points[i].X = z[0][i]
		points[i].Y = z[1][i]
		if len(z) > 2 {
			points[i].Z = z[2][i]
		}
	}
	return points
}

func FullNormalized(optMatrix [][]float64) [][]float64 {
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

func IdealPoint(normMatrix [][]float64, w []float64) ([]float64, float64, int) {
	numCrit := len(normMatrix)
	numAlt := len(normMatrix[0])
	idealPointMatrix := make([][]float64, numCrit)
	for i := 0; i < numCrit; i++ {
		idealPointMatrix[i] = make([]float64, numAlt)
	}

	for i := 0; i < numCrit; i++ {
		_, max := MinMax(normMatrix[i])
		for j := 0; j < numAlt; j++ {
			idealPointMatrix[i][j] = w[i] * w[i] * sqr(normMatrix[i][j]-max)
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

func AntiIdealPoint(normMatrix [][]float64, w []float64) ([]float64, float64, int) {
	numCrit := len(normMatrix)
	numAlt := len(normMatrix[0])
	idealPointMatrix := make([][]float64, numCrit)
	for i := 0; i < numCrit; i++ {
		idealPointMatrix[i] = make([]float64, numAlt)
	}

	for i := 0; i < numCrit; i++ {
		min, _ := MinMax(normMatrix[i])
		for j := 0; j < numAlt; j++ {
			idealPointMatrix[i][j] = w[i] * w[i] * sqr(normMatrix[i][j]-min)
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

func sqr(x float64) float64 { return x * x }

func AbsoluteXYZ(normMatrix [][]float64, w []float64) ([]float64, float64, int) {
	numCrit := len(normMatrix)
	numAlt := len(normMatrix[0])
	absoluteMatrix := make([]float64, numAlt)
	p := ToPoint(normMatrix)
	for i := 0; i < numAlt; i++ {
		if numCrit > 2 {
			absoluteMatrix[i] = p[i].X*w[0] + p[i].Y*w[1] + p[i].Z*w[2]
		} else {
			absoluteMatrix[i] = p[i].X*w[0] + p[i].Y*w[1]
		}

	}
	_, abMax := MinMax(absoluteMatrix)
	alt := 0
	for i, v := range absoluteMatrix {
		if v == abMax {
			alt = i + 1
		}
	}
	return absoluteMatrix, abMax, alt
}

func RelateXYZ(normMatrix [][]float64, w []float64) ([]float64, float64, int) {
	numCrit := len(normMatrix)
	numAlt := len(normMatrix[0])
	relateMatrix := make([]float64, numAlt)
	p := ToPoint(normMatrix)
	for i := 0; i < numAlt; i++ {
		if numCrit > 2 {
			relateMatrix[i] = math.Pow(p[i].X, w[0]) * math.Pow(p[i].Y, w[1]) * math.Pow(p[i].Z, w[2])
		} else {
			relateMatrix[i] = math.Pow(p[i].X, w[0]) * math.Pow(p[i].Y, w[1])
		}

	}
	_, reMax := MinMax(relateMatrix)
	alt := 0
	for i, v := range relateMatrix {
		if v == reMax {
			alt = i + 1
		}
	}
	return relateMatrix, reMax, alt
}

func IdealPointXYZ(normMatrix [][]float64, w []float64) ([]float64, float64, int) {
	numCrit := len(normMatrix)
	numAlt := len(normMatrix[0])
	idealPointMatrix := make([]float64, numAlt)
	_, maxZ1 := MinMax(normMatrix[0])
	_, maxZ2 := MinMax(normMatrix[1])
	_, maxZ3 := MinMax(normMatrix[2])
	p := ToPoint(normMatrix)
	for i := 0; i < numAlt; i++ {
		if numCrit > 2 {
			idealPointMatrix[i] = math.Pow((p[i].X-maxZ1), 2)*w[0]*w[0] + math.Pow((p[i].Y-maxZ2), 2)*w[1]*w[1] + math.Pow((p[i].Z-maxZ3), 2)*w[2]*w[2]
		} else {
			idealPointMatrix[i] = math.Pow((p[i].X-maxZ1), 2)*w[0]*w[0] + math.Pow((p[i].Y-maxZ2), 2)*w[1]*w[1]
		}

	}

	ip, _ := MinMax(idealPointMatrix)
	alt := 0
	for i, v := range idealPointMatrix {
		if v == ip {
			alt = i + 1
		}
	}
	return idealPointMatrix, ip, alt
}

func AntiIdealPointXYZ(normMatrix [][]float64, w []float64) ([]float64, float64, int) {
	numCrit := len(normMatrix)
	numAlt := len(normMatrix[0])
	antiIdealPointMatrix := make([]float64, numAlt)
	minZ1, _ := MinMax(normMatrix[0])
	minZ2, _ := MinMax(normMatrix[1])
	minZ3, _ := MinMax(normMatrix[2])
	p := ToPoint(normMatrix)
	for i := 0; i < numAlt; i++ {
		if numCrit > 2 {
			antiIdealPointMatrix[i] = math.Pow((p[i].X-minZ1), 2)*w[0]*w[0] + math.Pow((p[i].Y-minZ2), 2)*w[1]*w[1] + math.Pow((p[i].Z-minZ3), 2)*w[2]*w[2]
		} else {
			antiIdealPointMatrix[i] = math.Pow((p[i].X-minZ1), 2)*w[0]*w[0] + math.Pow((p[i].Y-minZ2), 2)*w[1]*w[1]
		}

	}

	_, aip := MinMax(antiIdealPointMatrix)
	alt := 0
	for i, v := range antiIdealPointMatrix {
		if v == aip {
			alt = i + 1
		}
	}
	return antiIdealPointMatrix, aip, alt
}

func AbsoluteСoncession(normMatrix [][]float64, w []float64) ([]float64, float64, int) {

	numCrit := len(normMatrix)
	numAlt := len(normMatrix[0])
	absoluteMatrix := make([][]float64, numCrit)

	for i := 0; i < numCrit; i++ {
		absoluteMatrix[i] = make([]float64, numAlt)
	}

	for i := 0; i < numCrit; i++ {
		for j := 0; j < numAlt; j++ {
			absoluteMatrix[i][j] = w[i] * normMatrix[i][j]
		}

	}
	//abRow := decisionsAp.RowToCol(absoluteMatrix)

	abArr := make([]float64, numAlt)
	for j := 0; j < numAlt; j++ {
		for i := 0; i < numCrit; i++ {
			//fmt.Println(idealPointMatrix[i])
			abArr[j] += absoluteMatrix[i][j]
		}
	}

	_, abMax := MinMax(abArr)
	alt := 0
	for i, v := range abArr {
		if v == abMax {
			alt = i + 1
		}
	}
	return abArr, abMax, alt
}

func RelateСoncession(normMatrix [][]float64, w []float64) ([]float64, float64, int) {
	numCrit := len(normMatrix)
	numAlt := len(normMatrix[0])
	relateMatrix := make([][]float64, numCrit)

	for i := 0; i < numCrit; i++ {
		relateMatrix[i] = make([]float64, numAlt)
	}

	for i := 0; i < numCrit; i++ {
		for j := 0; j < numAlt; j++ {
			relateMatrix[i][j] = math.Pow(normMatrix[i][j], w[i])
		}
	}
	//fmt.Println(relateMatrix)
	//abRow := decisionsAp.RowToCol(absoluteMatrix)
	reArr := make([]float64, numAlt)
	for j := 0; j < numAlt; j++ {
		reArr[j] = 1
		for i := 0; i < numCrit; i++ {
			//fmt.Println(relateMatrix[i][j])
			reArr[j] *= relateMatrix[i][j]
			//fmt.Println(reArr[j])
		}
	}
	_, reMax := MinMax(reArr)
	alt := 0
	for i, v := range reArr {
		if v == reMax {
			alt = i + 1
		}
	}
	return reArr, reMax, alt
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
