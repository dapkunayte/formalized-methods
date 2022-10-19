package decisionsAp

import (
	"fmt"
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

func IdealPointAp(c [][][]float64) {
	numCrit := len(c)
	//numAlt := len(c[0])
	numL := len(c[0][0])
	idealPoint := make([][]float64, numL)
	for i := 0; i < numL; i++ {
		idealPoint[i] = make([]float64, numCrit)
	}
	p := 0
	k := 0
	for j := 0; j < numL; j++ {
		for i := 0; i < numCrit; i++ {
			rc := RowToCol(c[p])
			//fmt.Println(rc[0])
			max := MaxInConvolution(rc)
			//fmt.Println(rc)
			//fmt.Println(rc[j][k], max[j])
			idealPoint[j][i] += math.Pow(rc[j][k]-max[j], 2)
			p++
		}

		if k == numCrit-1 {
			k = 0
		} else {
			k++
		}
		p = 0
	}
	fmt.Println(idealPoint[0])
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

func MaxInConvolution(r [][]float64) []float64 {
	numBStd := len(r)
	//lenL := len(r[0])
	max := make([]float64, numBStd)
	for i := 0; i < numBStd; i++ {
		_, maxArr := MinMax(r[i])
		max[i] = maxArr
	}
	return max
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
