package fuzzySets

import (
	"fmt"
	"math"
	"strconv"
)

func Rm(f [][]float64) ([][]float64, [][]float64) {
	rmMatrix := make([][]float64, len(f[0]))
	for i := 0; i < len(f[0]); i++ {
		rmMatrix[i] = make([]float64, len(f[0]))
	}
	for i := 0; i < len(rmMatrix); i++ {
		//fmt.Println(i)
		for j := 0; j < len(rmMatrix[0]); j++ {
			s := fmt.Sprintf("%.2f", (1 - f[0][i]))
			aM1, _ := strconv.ParseFloat(s, 64)
			rmMatrix[i][j] = math.Max(math.Min(f[0][i], f[1][j]), math.Min(aM1, f[2][j]))
			//fmt.Println(f[0][i], f[1][j], 1-f[0][i], f[2][j])
		}
	}
	convolutionMatrix := make([][]float64, len(f))
	for i := 0; i < len(f); i++ {
		convolutionMatrix[i] = make([]float64, len(f[0]))
	}
	for j := 0; j < len(f); j++ {
		for p := 0; p < len(f[0]); p++ {
			for i := 0; i < len(f[0]); i++ {
				if convolutionMatrix[j][p] < math.Min(f[j][i], rmMatrix[i][p]) {
					convolutionMatrix[j][p] = math.Min(f[j][i], rmMatrix[i][p])
				}

				//fmt.Println(f[j][i], rmMatrix[i][p])
			}
		}
	}

	return rmMatrix, convolutionMatrix
}

func Ra(f [][]float64) ([][]float64, [][]float64) {
	raMatrix := make([][]float64, len(f[0]))
	for i := 0; i < len(f[0]); i++ {
		raMatrix[i] = make([]float64, len(f[0]))
	}
	for i := 0; i < len(raMatrix); i++ {
		//fmt.Println(i)
		for j := 0; j < len(raMatrix[0]); j++ {
			s := fmt.Sprintf("%.2f", (1 - f[0][i]))
			aM1, _ := strconv.ParseFloat(s, 64)
			raMatrix[i][j] = math.Min(math.Min(f[0][i]+f[2][j], aM1+f[1][j]), 1)
			//fmt.Println(f[0][i], f[1][j], 1-f[0][i], f[2][j])
		}
	}
	convolutionMatrix := make([][]float64, len(f))
	for i := 0; i < len(f); i++ {
		convolutionMatrix[i] = make([]float64, len(f[0]))
	}
	for j := 0; j < len(f); j++ {
		for p := 0; p < len(f[0]); p++ {
			for i := 0; i < len(f[0]); i++ {
				if convolutionMatrix[j][p] < math.Min(f[j][i], raMatrix[i][p]) {
					convolutionMatrix[j][p] = math.Min(f[j][i], raMatrix[i][p])
				}

				//fmt.Println(f[j][i], rmMatrix[i][p])
			}
		}
	}

	return raMatrix, convolutionMatrix
}

func Rb(f [][]float64) ([][]float64, [][]float64) {
	rbMatrix := make([][]float64, len(f[0]))
	for i := 0; i < len(f[0]); i++ {
		rbMatrix[i] = make([]float64, len(f[0]))
	}
	for i := 0; i < len(rbMatrix); i++ {
		//fmt.Println(i)
		for j := 0; j < len(rbMatrix[0]); j++ {
			s := fmt.Sprintf("%.2f", (1 - f[0][i]))
			aM1, _ := strconv.ParseFloat(s, 64)
			rbMatrix[i][j] = math.Min(math.Max(f[0][i], f[2][j]), math.Max(aM1, f[1][j]))
			//fmt.Println(f[0][i], f[1][j], 1-f[0][i], f[2][j])
		}
	}
	convolutionMatrix := make([][]float64, len(f))
	for i := 0; i < len(f); i++ {
		convolutionMatrix[i] = make([]float64, len(f[0]))
	}
	for j := 0; j < len(f); j++ {
		for p := 0; p < len(f[0]); p++ {
			for i := 0; i < len(f[0]); i++ {
				if convolutionMatrix[j][p] < math.Min(f[j][i], rbMatrix[i][p]) {
					convolutionMatrix[j][p] = math.Min(f[j][i], rbMatrix[i][p])
				}

				//fmt.Println(f[j][i], rmMatrix[i][p])
			}
		}
	}

	return rbMatrix, convolutionMatrix
}

func Rgg(f [][]float64) ([][]float64, [][]float64) {
	rggMatrix := make([][]float64, len(f[0]))
	for i := 0; i < len(f[0]); i++ {
		rggMatrix[i] = make([]float64, len(f[0]))
	}
	a, b := 0.0, 0.0
	for i := 0; i < len(rggMatrix); i++ {
		//fmt.Println(i)
		for j := 0; j < len(rggMatrix[0]); j++ {
			if f[0][i] <= f[1][j] {
				a = 1
			} else {
				a = f[1][j]
			}
			s := fmt.Sprintf("%.2f", (1 - f[0][i]))
			aM1, _ := strconv.ParseFloat(s, 64)
			if aM1 <= f[2][j] {
				b = 1
			} else {
				b = f[2][j]
			}
			rggMatrix[i][j] = math.Min(a, b)
			//fmt.Println(f[0][i], f[1][j], 1-f[0][i], f[2][j])
		}
	}
	convolutionMatrix := make([][]float64, len(f))
	for i := 0; i < len(f); i++ {
		convolutionMatrix[i] = make([]float64, len(f[0]))
	}
	for j := 0; j < len(f); j++ {
		for p := 0; p < len(f[0]); p++ {
			for i := 0; i < len(f[0]); i++ {
				if convolutionMatrix[j][p] < math.Min(f[j][i], rggMatrix[i][p]) {
					convolutionMatrix[j][p] = math.Min(f[j][i], rggMatrix[i][p])
				}

				//fmt.Println(f[j][i], rmMatrix[i][p])
			}
		}
	}

	return rggMatrix, convolutionMatrix
}

func Rss(f [][]float64) ([][]float64, [][]float64) {
	rssMatrix := make([][]float64, len(f[0]))
	for i := 0; i < len(f[0]); i++ {
		rssMatrix[i] = make([]float64, len(f[0]))
	}
	a, b := 0.0, 0.0
	for i := 0; i < len(rssMatrix); i++ {
		//fmt.Println(i)
		for j := 0; j < len(rssMatrix[0]); j++ {
			if f[0][i] <= f[1][j] {
				a = 1
			} else {
				a = 0
			}
			s := fmt.Sprintf("%.2f", (1 - f[0][i]))
			aM1, _ := strconv.ParseFloat(s, 64)
			if aM1 <= f[2][j] {
				b = 1
			} else {
				b = 0
			}
			rssMatrix[i][j] = math.Min(a, b)
			//fmt.Println(f[0][i], f[1][j], 1-f[0][i], f[2][j])
		}
	}
	convolutionMatrix := make([][]float64, len(f))
	for i := 0; i < len(f); i++ {
		convolutionMatrix[i] = make([]float64, len(f[0]))
	}
	for j := 0; j < len(f); j++ {
		for p := 0; p < len(f[0]); p++ {
			for i := 0; i < len(f[0]); i++ {
				if convolutionMatrix[j][p] < math.Min(f[j][i], rssMatrix[i][p]) {
					convolutionMatrix[j][p] = math.Min(f[j][i], rssMatrix[i][p])
				}

				//fmt.Println(f[j][i], rmMatrix[i][p])
			}
		}
	}

	return rssMatrix, convolutionMatrix
}

func Rsg(f [][]float64) ([][]float64, [][]float64) {
	rsgMatrix := make([][]float64, len(f[0]))
	for i := 0; i < len(f[0]); i++ {
		rsgMatrix[i] = make([]float64, len(f[0]))
	}
	a, b := 0.0, 0.0
	for i := 0; i < len(rsgMatrix); i++ {
		//fmt.Println(i)
		for j := 0; j < len(rsgMatrix[0]); j++ {
			if f[0][i] <= f[1][j] {
				a = 1
			} else {
				a = 0
			}
			s := fmt.Sprintf("%.2f", (1 - f[0][i]))
			aM1, _ := strconv.ParseFloat(s, 64)
			if aM1 <= f[2][j] {
				b = 1
			} else {
				b = f[2][j]
			}
			rsgMatrix[i][j] = math.Min(a, b)
			//fmt.Println(f[0][i], f[1][j], 1-f[0][i], f[2][j])
		}
	}
	convolutionMatrix := make([][]float64, len(f))
	for i := 0; i < len(f); i++ {
		convolutionMatrix[i] = make([]float64, len(f[0]))
	}
	for j := 0; j < len(f); j++ {
		for p := 0; p < len(f[0]); p++ {
			for i := 0; i < len(f[0]); i++ {
				if convolutionMatrix[j][p] < math.Min(f[j][i], rsgMatrix[i][p]) {
					convolutionMatrix[j][p] = math.Min(f[j][i], rsgMatrix[i][p])
				}

				//fmt.Println(f[j][i], rmMatrix[i][p])
			}
		}
	}

	return rsgMatrix, convolutionMatrix
}

func Rgs(f [][]float64) ([][]float64, [][]float64) {
	rgsMatrix := make([][]float64, len(f[0]))
	for i := 0; i < len(f[0]); i++ {
		rgsMatrix[i] = make([]float64, len(f[0]))
	}
	a, b := 0.0, 0.0
	for i := 0; i < len(rgsMatrix); i++ {
		//fmt.Println(i)
		for j := 0; j < len(rgsMatrix[0]); j++ {
			if f[0][i] <= f[1][j] {
				a = 1
			} else {
				a = f[1][j]
			}
			s := fmt.Sprintf("%.2f", (1 - f[0][i]))
			aM1, _ := strconv.ParseFloat(s, 64)
			if aM1 <= f[2][j] {
				b = 1
			} else {
				b = 0
			}
			rgsMatrix[i][j] = math.Min(a, b)
			//fmt.Println(f[0][i], f[1][j], 1-f[0][i], f[2][j])
		}
	}
	convolutionMatrix := make([][]float64, len(f))
	for i := 0; i < len(f); i++ {
		convolutionMatrix[i] = make([]float64, len(f[0]))
	}
	for j := 0; j < len(f); j++ {
		for p := 0; p < len(f[0]); p++ {
			for i := 0; i < len(f[0]); i++ {
				if convolutionMatrix[j][p] < math.Min(f[j][i], rgsMatrix[i][p]) {
					convolutionMatrix[j][p] = math.Min(f[j][i], rgsMatrix[i][p])
				}

				//fmt.Println(f[j][i], rmMatrix[i][p])
			}
		}
	}

	return rgsMatrix, convolutionMatrix
}
