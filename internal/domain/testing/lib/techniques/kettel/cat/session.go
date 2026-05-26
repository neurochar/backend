package cat

import (
	"fmt"
)

type ResponsePair struct {
	ItemID   string `json:"item_id"`
	Response int    `json:"response"`
}

type Result struct {
	Factor             string         `json:"factor"`
	FactorName         string         `json:"factor_name"`
	Theta              float64        `json:"theta"`
	Sten               int            `json:"sten"`
	SE                 float64        `json:"se"`
	NItemsAdministered int            `json:"n_items_administered"`
	Converged          bool           `json:"converged"`
	StopReason         string         `json:"stop_reason"`
	Responses          []ResponsePair `json:"responses"`
}

type CATSession struct {
	factor       string
	factorName   string
	maxItems     int
	seThreshold  float64
	priorMean    float64
	priorSD      float64
	thetaGrid    []float64
	items        []Item
	itemLookup   map[string]Item
	remainingIDs []string
	answeredIDs  []string
	responses    []int
	thetaHat     float64
	se           float64
	done         bool
}

const (
	thetaLo = -4.0
	thetaHi = 4.0
	thetaN  = 801
)

func NewCATSession(factorParams Factor, maxItems int, seThreshold float64) *CATSession {
	s := &CATSession{
		factor:       factorParams.KettelLabel,
		factorName:   factorParams.Name,
		maxItems:     maxItems,
		seThreshold:  seThreshold,
		priorMean:    0.0,
		priorSD:      1.0,
		thetaGrid:    linspace(thetaLo, thetaHi, thetaN),
		items:        factorParams.Items,
		itemLookup:   make(map[string]Item),
		remainingIDs: make([]string, len(factorParams.Items)),
		answeredIDs:  make([]string, 0),
		responses:    make([]int, 0),
		thetaHat:     0.0,
		se:           1.0,
		done:         false,
	}

	for i, item := range factorParams.Items {
		s.itemLookup[item.ID] = item
		s.remainingIDs[i] = item.ID
	}

	s.updateEstimate()
	return s
}

func (s *CATSession) IsDone() bool {
	return s.done
}

func (s *CATSession) NextItem() *Item {
	if s.done {
		return nil
	}

	var bestID string
	var bestInfo float64 = -1

	for _, id := range s.remainingIDs {
		item := s.itemLookup[id]
		info := ItemInformation(s.thetaHat, item.A, item.B)
		if info > bestInfo {
			bestInfo = info
			bestID = id
		}
	}

	item := s.itemLookup[bestID]
	return &item
}

func (s *CATSession) RecordResponse(itemID string, response int) error {
	found := false
	for _, id := range s.remainingIDs {
		if id == itemID {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("item %s not available (already answered or unknown)", itemID)
	}
	if response < 0 || response > 2 {
		return fmt.Errorf("response must be 0, 1, or 2; got %d", response)
	}

	newRemaining := make([]string, 0, len(s.remainingIDs)-1)
	for _, id := range s.remainingIDs {
		if id != itemID {
			newRemaining = append(newRemaining, id)
		}
	}
	s.remainingIDs = newRemaining

	s.answeredIDs = append(s.answeredIDs, itemID)
	s.responses = append(s.responses, response)

	s.updateEstimate()
	s.checkDone()
	return nil
}

func (s *CATSession) Result() Result {
	return Result{
		Factor:             s.factor,
		FactorName:         s.factorName,
		Theta:              roundTo(s.thetaHat, 4),
		Sten:               ThetaToSten(s.thetaHat),
		SE:                 roundTo(s.se, 4),
		NItemsAdministered: len(s.answeredIDs),
		Converged:          s.se < s.seThreshold,
		StopReason:         s.stopReason(),
		Responses:          s.responsePairs(),
	}
}

func (s *CATSession) Theta() float64 {
	return s.thetaHat
}

func (s *CATSession) SE() float64 {
	return s.se
}

func (s *CATSession) updateEstimate() {
	post := posterior(s.thetaGrid, s.answeredIDs, s.responses, s.itemLookup, s.priorMean, s.priorSD)
	s.thetaHat, s.se = eapEstimate(post, s.thetaGrid)
}

func (s *CATSession) checkDone() {
	switch {
	case s.se < s.seThreshold:
		s.done = true
	case len(s.answeredIDs) >= s.maxItems:
		s.done = true
	case len(s.remainingIDs) == 0:
		s.done = true
	}
}

func (s *CATSession) stopReason() string {
	switch {
	case s.se < s.seThreshold:
		return "se_threshold"
	case len(s.answeredIDs) >= s.maxItems:
		return "max_items"
	default:
		return "all_items_exhausted"
	}
}

func (s *CATSession) responsePairs() []ResponsePair {
	pairs := make([]ResponsePair, len(s.answeredIDs))
	for i, id := range s.answeredIDs {
		pairs[i] = ResponsePair{
			ItemID:   id,
			Response: s.responses[i],
		}
	}
	return pairs
}
