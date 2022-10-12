package decisions

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

func IdealPoint(normMatrix [][]float64) ([]float64, float64, int) {
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

func AntiIdealPoint(normMatrix [][]float64) ([]float64, float64, int) {
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

func sqr(x float64) float64 { return x * x }

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
