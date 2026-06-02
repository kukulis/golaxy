package dao

import (
	"maps"
	"slices"

	"glaktika.eu/galaktika/pkg/galaxy"
	"glaktika.eu/galaktika/pkg/util"
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

func (r *ChallengeRepository) GetAll(filter galaxy.ChallengesFilter) []*galaxy.Challenge {
	return util.ArrayFilter(slices.Collect(maps.Values(r.challengesMap)), filter.Match)
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
