package predictions

import (
	compute "github.com/go-compute/mean/pkg"
	"math"
)

type PredictModel struct {
}

func LessSquare(x, y []float64) ([]float64, []float64) {
	f := func(a, x, b float64) float64 { return a + x*b }
	yNew := make([]float64, len(y))
	predict := make([]float64, 3)
	n := float64(len(y))
	//var a, b float64
	var sumXY, sumX, sumY, sumXSqr float64
	for i := 0; i < len(y); i++ {
		sumXY += x[i] * y[i]
		sumY += y[i]
		sumX += x[i]
		sumXSqr += x[i] * x[i]
	}
	matrix := [][]float64{
		{n, sumX, sumY},
		{sumX, sumXSqr, sumXY},
	}
	ans := GaussMehtodforSqr(matrix)
	for i := 0; i < len(y); i++ {
		yNew[i] = f(ans[0], x[i], ans[1])
	}
	k := 1
	for i := 0; i < 3; i++ {
		predict[i] = f(ans[0], float64(len(x)+k), ans[1])
		k++
	}

	return yNew, predict
}

func SmoothMean3(x, y []float64) []float64 {
	yNew := make([]float64, len(y)-2)
	xNew := make([]float64, len(x)-2)
	m := make([]float64, len(x)-2)
	for i := 0; i < len(xNew); i++ {
		xNew[i] = x[i] + 1
		m[i] = (y[i] + y[i+1] + y[i+2]) / 3
	}
	for i := 0; i < len(yNew); i++ {
		yNew[i] = m[i] + 1/3*(y[i+1]-y[i])
	}
	return yNew
}

func ExponentialSmooth(x, y []float64, a float64) ([]float64, []float64) {

	n := float64(len(y))
	var sumXY, sumX, sumY, sumXSqr float64

	for i := 0; i < len(y); i++ {
		sumXY += x[i] * y[i]
		sumY += y[i]
		sumX += x[i]
		sumXSqr += x[i] * x[i]
	}
	matrix := [][]float64{
		{n, sumX, sumY},
		{sumX, sumXSqr, sumXY},
	}
	ans := GaussMehtodforSqr(matrix)

	s1 := make([]float64, len(y))
	s2 := make([]float64, len(y))
	s10 := ans[0] - ((1-a)/a)*ans[1]
	s20 := ans[0] - (2*(1-a)/a)*ans[1]
	s1[0] = a*y[0] + (1-a)*s10
	s2[0] = a*s1[0] + (1-a)*s20
	for i := 1; i < len(s1); i++ {
		s1[i] = a*y[i] + (1-a)*s1[i-1]
		s2[i] = a*s1[i] + (1-a)*s2[i-1]
	}

	f := func(a, x, b float64) float64 { return a + x*b }
	predict := make([]float64, 3)
	yNew := make([]float64, len(y))
	for i := 0; i < len(y); i++ {
		a0New := 2*s1[i] - s2[i]
		a1New := (a / (1 - a)) * (s1[i] - s2[i])
		yNew[i] = f(a0New, 0, a1New)
		//fmt.Println(a0New, a1New)
	}
	k := 1.0
	for i := 0; i < len(predict); i++ {
		a0New := 2*s1[len(y)-1] - s2[len(y)-1]
		a1New := (a / (1 - a)) * (s1[len(y)-1] - s2[len(y)-1])
		predict[i] = f(a0New, float64(k), a1New)
		k++
	}

	return yNew, predict
}

func LessSquareSqr(x, y []float64) ([]float64, []float64) {
	n := float64(len(y))
	yNew := make([]float64, len(y))
	predict := make([]float64, 3)
	var (
		sumX   float64
		sumX2  float64
		sumX3  float64
		sumX4  float64
		sumY   float64
		sumXY  float64
		sumX2Y float64
	)
	for i := 0; i < len(y); i++ {
		sumX += x[i]
		sumX2 += x[i] * x[i]
		sumX3 += x[i] * x[i] * x[i]
		sumX4 += x[i] * x[i] * x[i] * x[i]
		sumY += y[i]
		sumXY += y[i] * x[i]
		sumX2Y += y[i] * x[i] * x[i]
	}
	matrix := [][]float64{
		{sumX4, sumX3, sumX2, sumX2Y},
		{sumX3, sumX2, sumX, sumXY},
		{sumX2, sumX, n, sumY},
	}
	a := GaussMehtodforSqr(matrix)
	f := func(a, x, b, c float64) float64 { return a*x*x + b*x + c }
	for i := 0; i < len(y); i++ {
		yNew[i] = f(a[0], x[i], a[1], a[2])
	}
	k := 1
	for i := 0; i < 3; i++ {
		predict[i] = f(a[0], float64(len(x)+k), a[1], a[2])
		k++
	}

	return yNew, predict
}

func GaussMehtodforSqr(matrix [][]float64) []float64 {
	n := len(matrix)
	answer := make([]float64, n)
	var (
		//k   int
		tmp float64
	)
	for i := 0; i < n; i++ {
		tmp = matrix[i][i]
		for j := n; j >= i; j-- {
			matrix[i][j] /= tmp
		}
		for j := i + 1; j < n; j++ {
			tmp = matrix[j][i]
			for k := n; k >= i; k-- {
				matrix[j][k] -= tmp * matrix[i][k]
			}
		}
	}
	answer[n-1] = matrix[n-1][n]
	for i := n - 2; i >= 0; i-- {
		answer[i] = matrix[i][n]
		for j := i + 1; j < n; j++ {
			answer[i] -= matrix[i][j] * answer[j]
		}
	}
	return answer
}

func ExponentialSmoothSqr(x, y []float64, a float64) ([]float64, []float64) {
	n := float64(len(y))
	yNew := make([]float64, len(y))
	predict := make([]float64, 3)
	var (
		sumX   float64
		sumX2  float64
		sumX3  float64
		sumX4  float64
		sumY   float64
		sumXY  float64
		sumX2Y float64
	)
	for i := 0; i < len(y); i++ {
		sumX += x[i]
		sumX2 += x[i] * x[i]
		sumX3 += x[i] * x[i] * x[i]
		sumX4 += x[i] * x[i] * x[i] * x[i]
		sumY += y[i]
		sumXY += y[i] * x[i]
		sumX2Y += y[i] * x[i] * x[i]
	}
	matrix := [][]float64{
		{sumX4, sumX3, sumX2, sumX2Y},
		{sumX3, sumX2, sumX, sumXY},
		{sumX2, sumX, n, sumY},
	}
	ans := GaussMehtodforSqr(matrix)

	s1 := make([]float64, len(y))
	s2 := make([]float64, len(y))
	s3 := make([]float64, len(y))
	s10 := ans[0] - (((1 - a) / a) * ans[1]) + (((1-a)*(2-a))/(2*a*a))*ans[2]
	s20 := ans[0] - (((2 * (1 - a)) / a) * ans[1]) + (((1-a)*(3-2*a))/(2*a*a))*ans[2]
	s30 := ans[0] - ((3 * (1 - a) / a) * ans[1]) + (((1-a)*(4-3*a))/(2*a*a))*ans[2]

	s1[0] = a*y[0] + (1-a)*s10
	s2[0] = a*s1[0] + (1-a)*s20
	s3[0] = a*s2[0] + (1-a)*s30
	for i := 1; i < len(s1); i++ {
		s1[i] = a*y[i] + (1-a)*s1[i-1]
		s2[i] = a*s1[i] + (1-a)*s2[i-1]
		s3[i] = a*s2[i] + (1-a)*s3[i-1]
	}

	f := func(a, b, c, x float64) float64 { return a + b*x + (c/2)*x*x }
	for i := 0; i < len(y); i++ {
		a0New := 3*(s1[i]-s2[i]) + s3[i]
		a1New := (a / (2 * (1 - a))) * (((6 - 5*a) * s1[i]) - (2 * (5 - 4*a) * s2[i]) + ((4 - 3*a) * s3[i]))
		a2New := ((a / (1 - a)) * (a / (1 - a))) * (s1[i] - 2*s2[i] + s3[i])
		yNew[i] = f(a0New, a1New, a2New, 0)
	}

	k := 1.0
	for i := 0; i < len(predict); i++ {
		a0New := 3*(s1[len(y)-1]-s2[len(y)-1]) + s3[len(y)-1]
		a1New := (a / (2 * (1 - a))) * (((6 - 5*a) * s1[len(y)-1]) - (2 * (5 - 4*a) * s2[len(y)-1]) + ((4 - 3*a) * s3[len(y)-1]))
		a2New := ((a / (1 - a)) * (a / (1 - a))) * (s1[len(y)-1] - 2*s2[len(y)-1] + s3[len(y)-1])
		predict[i] = f(a0New, a1New, a2New, k)
		k++
	}
	return yNew, predict
}

func Errors(y, yPred, x []float64, p float64) map[string]float64 {
	errs := make(map[string]float64)
	var s1, s1f, b, ma, xm, nu float64
	n := float64(len(y))

	for i := 0; i < int(n); i++ {
		s1 += math.Pow(y[i]-compute.Mean(y), 2) / (n - 1) //по формуле в квадрате
		s1f += math.Pow(y[i]-yPred[i], 2)                 //числитель
		b += math.Abs(y[i]-yPred[i]) / math.Sqrt(n*(n-1))
		ma += math.Abs((y[i] - yPred[i]) / (y[i]))
		xm += math.Pow(x[i]-compute.Mean(x), 2)
	}

	s1f = math.Sqrt(s1f / (n - p))
	s1 = math.Sqrt(s1)
	nu = math.Sqrt((s1f * s1f) / (s1 * s1))
	errs["n"] = nu
	errs["b"] = b
	errs["ma"] = (1 / n) * ma * 100
	errs["s1f"] = s1f
	errs["f"] = (s1 * s1) / (s1f * s1f)
	errs["s"] = s1f * math.Sqrt(1+(1/n)+(xm/(s1*s1)))

	return errs
}
