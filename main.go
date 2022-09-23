package main

import "fmt"

func rankingTwoDimensional(mainMatrix [][]int) [][]float64 {
	n := len(mainMatrix)    //4
	m := len(mainMatrix[1]) //6
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

func pairComparsion(mainMatrix [][]int) [][][]int {
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
					compMatrix[p][i][j] = 1
					//fmt.Println(1)
				}
				if mainMatrix[p][i] < mainMatrix[p][j] {
					compMatrix[p][i][j] = -1
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

func main() {
	/*
		{2, 1, 1, 1},
		{1, 2, 1, 1},
		{3, 2, 1, 2},
		{4, 3, 3, 3},
		{4, 4, 2, 3},
		{5, 5, 4, 4},
	*/

	mainMatrix := [][]int{
		{2, 1, 3, 4, 4, 5},
		{1, 2, 2, 3, 4, 5},
		{1, 1, 1, 3, 2, 4},
		{1, 1, 2, 3, 3, 4},
	}

	//fmt.Println(rankingTwoDimensional(mainMatrix))
	resultRanking := rankingTwoDimensional(mainMatrix)
	result := pairComparsion(mainMatrix)
	fmt.Println("Ранжировка:", resultRanking, "\n")
	for i := range result {
		fmt.Println("Эксперт №", i+1, ":", result[i], "\n")
	}
}
