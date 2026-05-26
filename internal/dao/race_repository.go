package dao

import (
	"glaktika.eu/galaktika/pkg/galaxy"
	"maps"
	"slices"
	"strings"
)

type RaceRepository struct {
	raceMap map[string]*galaxy.Race
}

func NewRaceRepository() *RaceRepository {
	return &RaceRepository{
		raceMap: make(map[string]*galaxy.Race),
	}
}

func (r *RaceRepository) Get(id string) *galaxy.Race {
	return r.raceMap[id]
}

func (r *RaceRepository) GetByToken(token string) *galaxy.Race {
	for _, race := range r.raceMap {
		if race.Token == token {
			return race
		}
	}
	return nil
}

func (r *RaceRepository) GetAll() []*galaxy.Race {
	races := slices.Collect(maps.Values(r.raceMap))

	slices.SortFunc(races, func(a, b *galaxy.Race) int {
		return strings.Compare(a.ID, b.ID)
	})

	return races
}

func (r *RaceRepository) Upsert(race *galaxy.Race) {
	r.raceMap[race.ID] = race
}

func (r *RaceRepository) Delete(id string) {
	delete(r.raceMap, id)
}

func (r *RaceRepository) ResetData() {
	r.raceMap = make(map[string]*galaxy.Race)
}
