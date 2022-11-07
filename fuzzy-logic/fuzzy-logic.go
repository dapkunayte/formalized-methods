package fuzzy_logic

var ruleBase = [][]string{
	{"bad", "bad", "bad"},          //r1
	{"bad", "bad", "medium"},       //r2
	{"bad", "medium", "medium"},    //r3
	{"bad", "medium", "bad"},       //r4
	{"medium", "bad", "bad"},       //r5
	{"medium", "medium", "bad"},    //r6
	{"medium", "medium", "medium"}, //r7
	{"medium", "bad", "medium"},    //r8
	{"medium", "medium", "good"},   //r9
	{"medium", "good", "medium"},   //r10
	{"medium", "good", "good"},     //r11
	{"good", "good", "good"},       //r12
	{"good", "good", "medium"},     //r13
	{"good", "medium", "medium"},   //r14
	{"good", "medium", "good"},     //r15
}

var ruleAgg = []string{
	"bad", "bad", "medium", "bad", "bad", "medium", "medium", "medium", "medium", "medium", "high", "high", "high", "medium", "high",
}

func Triangle(x float64, fuzzyPow []float64) float64 {
	switch {
	case x <= fuzzyPow[0]:
		return 0.0
	case fuzzyPow[0] < x && x <= fuzzyPow[1]:
		return (x - fuzzyPow[0]) / (fuzzyPow[1] - fuzzyPow[0])
	case fuzzyPow[1] < x && x < fuzzyPow[2]:
		return (fuzzyPow[2] - x) / (fuzzyPow[2] - fuzzyPow[1])
	case x >= fuzzyPow[2]:
		return 0
	}
	return 0
}

func Fuzzification(alts [][]float64, fuzzyPow []float64) [][]float64 {
	fuzzyMatrix := make([][]float64, len(alts))
	for i := 0; i < len(alts); i++ {
		fuzzyMatrix[i] = make([]float64, len(alts[0]))
	}
	for i := 0; i < len(alts); i++ {
		for j := 0; j < len(alts[0]); j++ {
			fuzzyMatrix[i][j] = Triangle(alts[i][j], fuzzyPow)
		}
	}
	return fuzzyMatrix
}

func RulesConv(bad, medium, high [][]float64, alt int) [][]float64 {
	rulesMatrix := make([][]float64, len(ruleBase))
	for i := 0; i < len(ruleBase); i++ {
		rulesMatrix[i] = make([]float64, len(ruleBase[0]))
	}
	for i := 0; i < len(ruleBase); i++ {
		for j := 0; j < len(ruleBase[0]); j++ {
			switch ruleBase[i][j] {
			case "bad":
				rulesMatrix[i][j] = bad[j][alt]
			case "medium":
				rulesMatrix[i][j] = medium[j][alt]
			case "high":
				rulesMatrix[i][j] = high[j][alt]
			}
		}
	}
	tNorm := make([]float64, len(rulesMatrix))
	for i := 0; i < len(rulesMatrix); i++ {
		min, _ := MinMax(rulesMatrix[i])
		tNorm[i] = min
	}

	fuzzyPowMaxMatrix := make([]float64, 3)

	for j := 0; j < len(ruleAgg); j++ {
		switch ruleAgg[j] {
		case "bad":
			if fuzzyPowMaxMatrix[0] <= tNorm[j] {
				fuzzyPowMaxMatrix[0] = tNorm[j]
			}
		case "medium":
			if fuzzyPowMaxMatrix[1] <= tNorm[j] {
				fuzzyPowMaxMatrix[1] = tNorm[j]
			}
		case "high":
			if fuzzyPowMaxMatrix[2] <= tNorm[j] {
				fuzzyPowMaxMatrix[2] = tNorm[j]
			}
		}
	}
	return rulesMatrix
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
