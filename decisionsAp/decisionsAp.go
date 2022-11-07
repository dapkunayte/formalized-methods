package decisionsAp

import (
	"math"
)

func Bl(z [][]float64, p []float64) []float64 {
	numP := len(p)
	numAlt := len(z)
	b := make([]float64, numAlt)
	for i := 0; i < numAlt; i++ {
		for j := 0; j < numP; j++ {
			b[i] += z[i][j] * p[j]
		}
	}
	return b
}

func Std(z [][]float64, p []float64) []float64 {
	numP := len(p)
	numAlt := len(z)
	s := make([]float64, numAlt)
	b := Bl(z, p)
	for i := 0; i < numAlt; i++ {
		for j := 0; j < numP; j++ {
			s[i] += math.Pow(z[i][j]-b[i], 2) * p[j]
		}
		s[i] = math.Sqrt(s[i])
	}
	return s
}

func ConvolutionStdBl(b []float64, s []float64) [][]float64 {
	numBStd := len(b)
	l := []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1}
	convolutionMatrix := make([][]float64, numBStd)
	for i := 0; i < numBStd; i++ {
		convolutionMatrix[i] = make([]float64, len(l))
	}
	for i := 0; i < numBStd; i++ {
		for j := 0; j < len(l); j++ {
			convolutionMatrix[i][j] = ((1 - l[j]) * b[i]) - (s[i] * l[j])
		}
	}
	return convolutionMatrix
}

func IdealPointAp(c [][][]float64) [][]float64 {
	numCrit := len(c)
	//numAlt := len(c[0])
	numAlt := len(c[0])
	numL := len(c[0][0])
	idealPoint := make([][]float64, numL)
	for i := 0; i < numL; i++ {
		idealPoint[i] = make([]float64, numAlt)
	}
	//k := 0
	//z := 0
	for p := 0; p < numCrit; p++ {
		rc := RowToCol(c[p])
		_, max := MinMaxInConvolution(rc)
		for i := 0; i < numL; i++ {
			for j := 0; j < numAlt; j++ {
				idealPoint[i][j] += (rc[i][j] - max[i]) * (rc[i][j] - max[i])
				//fmt.Println(rc[i][j], max[i])
			}
			//fmt.Println("----------")
		}
	}
	return idealPoint
}

func RealteСoncessionAp(c [][][]float64) [][]float64 {
	numCrit := len(c)
	//numAlt := len(c[0])
	numAlt := len(c[0])
	numL := len(c[0][0])
	relateMatrix := make([][]float64, numL)
	for i := 0; i < numL; i++ {
		relateMatrix[i] = make([]float64, numAlt)
		for j := 0; j < numAlt; j++ {
			relateMatrix[i][j] = 1
		}
	}
	for p := 0; p < numCrit; p++ {
		rc := RowToCol(c[p])
		for i := 0; i < numL; i++ {
			for j := 0; j < numAlt; j++ {
				relateMatrix[i][j] *= rc[i][j]
			}
		}
	}
	return relateMatrix
}

func AbsoluteСoncessionAp(c [][][]float64) [][]float64 {
	numCrit := len(c)
	//numAlt := len(c[0])
	numAlt := len(c[0])
	numL := len(c[0][0])
	absoluteMatrix := make([][]float64, numL)
	for i := 0; i < numL; i++ {
		absoluteMatrix[i] = make([]float64, numAlt)
	}
	for p := 0; p < numCrit; p++ {
		rc := RowToCol(c[p])
		for i := 0; i < numL; i++ {
			for j := 0; j < numAlt; j++ {
				absoluteMatrix[i][j] += rc[i][j]
			}
		}
	}
	return absoluteMatrix
}

func AntiIdealPointAp(c [][][]float64) [][]float64 {
	numCrit := len(c)
	//numAlt := len(c[0])
	numAlt := len(c[0])
	numL := len(c[0][0])
	antiIdealPoint := make([][]float64, numL)
	for i := 0; i < numL; i++ {
		antiIdealPoint[i] = make([]float64, numAlt)
	}
	//k := 0
	//z := 0
	for p := 0; p < numCrit; p++ {
		rc := RowToCol(c[p])
		min, _ := MinMaxInConvolution(rc)
		for i := 0; i < numL; i++ {
			for j := 0; j < numAlt; j++ {
				antiIdealPoint[i][j] += (rc[i][j] - min[i]) * (rc[i][j] - min[i])
				//fmt.Println(rc[i][j], max[i])
			}
			//fmt.Println("----------")
		}
	}
	return antiIdealPoint
}

func RowToCol(c [][]float64) [][]float64 {
	numBStd := len(c)
	lenL := len(c[0])
	reverseMatrix := make([][]float64, lenL)
	for i := 0; i < lenL; i++ {
		reverseMatrix[i] = make([]float64, numBStd)
	}
	for i := 0; i < lenL; i++ {
		for j := 0; j < numBStd; j++ {
			reverseMatrix[i][j] = c[j][i]
		}
	}
	return reverseMatrix
}

func MinMaxInConvolution(r [][]float64) ([]float64, []float64) {
	numBStd := len(r)
	//lenL := len(r[0])
	max := make([]float64, numBStd)
	min := make([]float64, numBStd)
	for i := 0; i < numBStd; i++ {
		minArr, maxArr := MinMax(r[i])
		max[i] = maxArr
		min[i] = minArr
	}
	return min, max
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

func Z10(z [][]float64) [][]float64 {
	numAlt := len(z)
	//numCond := len(z[0])
	l := []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1}
	z10Matrix := make([][]float64, numAlt)
	for i := 0; i < numAlt; i++ {
		z10Matrix[i] = make([]float64, len(l))
	}
	for i := 0; i < numAlt; i++ {
		min, max := MinMax(z[i])
		for j := 0; j < len(l); j++ {
			z10Matrix[i][j] = l[j]*min + (1-l[j])*max
		}
	}
	return z10Matrix
}

func ZI(z7 [][]float64, z10 [][]float64) [][]float64 {
	numAlt := len(z7)
	//numCond := len(z[0])
	l := []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1}
	ziMatrix := make([][]float64, numAlt)
	for i := 0; i < numAlt; i++ {
		ziMatrix[i] = make([]float64, len(l))
	}
	for i := 0; i < numAlt; i++ {
		for j := 0; j < len(l); j++ {
			ziMatrix[i][j] = (1-l[j])*z7[i][j] + l[j]*z10[i][j]
		}
	}
	return ziMatrix
}

func NormalizedMatrix(matrix [][]float64, minOrMax string) [][]float64 {
	min, max := 0.0, 0.0
	switch minOrMax {
	case "max":
		min, max = MinMax(matrix[0])
		for i := 1; i < len(matrix); i++ {
			min1, max1 := MinMax(matrix[i])
			if min > min1 {
				min = min1
			}
			if max < max1 {
				max = max1
			}
		}
	case "min":
		max, min = MinMax(matrix[0])
		for i := 1; i < len(matrix); i++ {
			max1, min1 := MinMax(matrix[i])
			if max > max1 {
				max = max1
			}
			if min < min1 {
				min = min1
			}
		}
	}

	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[0]); j++ {
			matrix[i][j] = (matrix[i][j] - min) / (max - min)
		}
	}

	return matrix
}
