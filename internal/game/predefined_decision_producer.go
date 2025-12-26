package game

type PredefinedDecisionProducer struct {
	shots []ShotDecision
	index int
}

// NewPredefinedDecisionProducer creates a new predefined decision producer
func NewPredefinedDecisionProducer(shots []ShotDecision) *PredefinedDecisionProducer {
	return &PredefinedDecisionProducer{
		shots: shots,
		index: 0,
	}
}

// ProduceNextShot returns the next predefined shot
func (p *PredefinedDecisionProducer) ProduceNextShot() *ShotDecision {
	if p.index >= len(p.shots) {
		return nil
	}
	shot := p.shots[p.index]
	p.index++
	return &shot
}
