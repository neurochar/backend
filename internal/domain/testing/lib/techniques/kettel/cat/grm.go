package cat

import (
	"math"
)

const (
	defaultEps = 1e-12
)

type Item struct {
	ID string    `json:"id"`
	A  float64   `json:"a"`
	B  []float64 `json:"b"`
}

type Factor struct {
	Name        string `json:"name"`
	KettelLabel string `json:"kettel_label"`
	Items       []Item `json:"items"`
}

type GRMParams struct {
	Meta    map[string]interface{} `json:"meta"`
	Factors map[string]Factor      `json:"factors"`
}

func pge(theta, a float64, b []float64) []float64 {
	r := make([]float64, len(b))
	for i := range b {
		z := a * (theta - b[i])
		r[i] = 1.0 / (1.0 + math.Exp(-z))
	}
	return r
}

func pgeVec(theta []float64, a float64, b []float64) [][]float64 {
	m := len(theta)
	k := len(b)
	r := make([][]float64, m)
	for i := 0; i < m; i++ {
		r[i] = make([]float64, k)
		for j := 0; j < k; j++ {
			z := a * (theta[i] - b[j])
			r[i][j] = 1.0 / (1.0 + math.Exp(-z))
		}
	}
	return r
}

func pcatVec(theta []float64, a float64, b []float64) [][]float64 {
	pge := pgeVec(theta, a, b)
	m := len(theta)
	k := len(b) + 1

	res := make([][]float64, m)
	for i := 0; i < m; i++ {
		padded := make([]float64, k+1)
		padded[0] = 1.0
		for j := 0; j < len(b); j++ {
			padded[j+1] = pge[i][j]
		}
		padded[k] = 0.0

		pk := make([]float64, k)
		var sum float64
		for j := 0; j < k; j++ {
			pk[j] = padded[j] - padded[j+1]
			if pk[j] < defaultEps {
				pk[j] = defaultEps
			}
			sum += pk[j]
		}
		for j := 0; j < k; j++ {
			pk[j] /= sum
		}
		res[i] = pk
	}
	return res
}

func ItemInformation(theta, a float64, b []float64) float64 {
	pge := pge(theta, a, b)

	dPge := make([]float64, len(b))
	for i := range b {
		dPge[i] = a * pge[i] * (1.0 - pge[i])
	}

	k := len(b) + 1

	PgeP := make([]float64, k+1)
	PgeP[0] = 1.0
	for i := 0; i < len(b); i++ {
		PgeP[i+1] = pge[i]
	}
	PgeP[k] = 0.0

	dPgeP := make([]float64, k+1)
	dPgeP[0] = 0.0
	for i := 0; i < len(b); i++ {
		dPgeP[i+1] = dPge[i]
	}
	dPgeP[k] = 0.0

	var info float64
	for i := 0; i < k; i++ {
		pk := PgeP[i] - PgeP[i+1]
		if pk < defaultEps {
			pk = defaultEps
		}
		dPk := dPgeP[i] - dPgeP[i+1]
		info += dPk * dPk / pk
	}
	return info
}

func posterior(
	thetaGrid []float64,
	answeredIDs []string,
	responses []int,
	itemLookup map[string]Item,
	priorMean,
	priorSD float64,
) []float64 {
	m := len(thetaGrid)
	logp := make([]float64, m)

	for i := 0; i < m; i++ {
		diff := (thetaGrid[i] - priorMean) / priorSD
		logp[i] = -0.5 * diff * diff
	}

	for idx := range answeredIDs {
		item := itemLookup[answeredIDs[idx]]
		resp := responses[idx]
		pk := pcatVec(thetaGrid, item.A, item.B)

		for i := 0; i < m; i++ {
			prob := pk[i][resp]
			if prob < defaultEps {
				prob = defaultEps
			}
			logp[i] += math.Log(prob)
		}
	}

	maxLog := logp[0]
	for i := 1; i < m; i++ {
		if logp[i] > maxLog {
			maxLog = logp[i]
		}
	}

	post := make([]float64, m)
	var sum float64
	for i := 0; i < m; i++ {
		post[i] = math.Exp(logp[i] - maxLog)
		sum += post[i]
	}
	for i := 0; i < m; i++ {
		post[i] /= sum
	}
	return post
}

func eapEstimate(post, thetaGrid []float64) (float64, float64) {
	var thetaHat float64
	for i := range post {
		thetaHat += post[i] * thetaGrid[i]
	}

	var variance float64
	for i := range post {
		d := thetaGrid[i] - thetaHat
		variance += post[i] * d * d
	}
	if variance < 0 {
		variance = 0
	}
	se := math.Sqrt(variance)
	return thetaHat, se
}

func linspace(lo, hi float64, n int) []float64 {
	r := make([]float64, n)
	if n == 1 {
		r[0] = lo
		return r
	}
	step := (hi - lo) / float64(n-1)
	for i := 0; i < n; i++ {
		r[i] = lo + float64(i)*step
	}
	return r
}

func roundTo(x float64, decimals int) float64 {
	pow := math.Pow(10, float64(decimals))
	return math.Round(x*pow) / pow
}

func ThetaToSten(theta float64) int {
	s := 5.5 + 2.0*theta
	if s < 1 {
		return 1
	}
	if s > 10 {
		return 10
	}
	return int(math.Round(s))
}
