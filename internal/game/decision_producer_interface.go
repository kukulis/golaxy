package game

type DecisionProducerInterface interface {
	ProduceNextShot() *ShotDecision
}
