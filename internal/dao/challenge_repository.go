package dao

import (
	"maps"
	"slices"

	"glaktika.eu/galaktika/pkg/galaxy"
)

type ChallengeRepository struct {
	challengesMap map[string]*galaxy.Challenge
}

func NewChallengeRepository() *ChallengeRepository {
	return &ChallengeRepository{
		challengesMap: make(map[string]*galaxy.Challenge),
	}
}

func (r *ChallengeRepository) Get(id string) *galaxy.Challenge {
	return r.challengesMap[id]
}

func (r *ChallengeRepository) GetAll() []*galaxy.Challenge {
	return slices.Collect(maps.Values(r.challengesMap))
}

func (r *ChallengeRepository) Upsert(challenge *galaxy.Challenge) {
	r.challengesMap[challenge.ID] = challenge
}

func (r *ChallengeRepository) Delete(id string) {
	delete(r.challengesMap, id)
}

func (r *ChallengeRepository) ResetData() {
	r.challengesMap = make(map[string]*galaxy.Challenge)
}
