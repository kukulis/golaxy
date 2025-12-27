package dao

import (
	"glaktika.eu/galaktika/pkg/galaxy"
	"maps"
	"slices"
	"strings"
)

type DivisionRepository struct {
	divisionMap map[string]*galaxy.Division
}

func NewDivisionRepository() *DivisionRepository {
	return &DivisionRepository{
		divisionMap: make(map[string]*galaxy.Division),
	}
}

func (r *DivisionRepository) Get(id string) *galaxy.Division {
	return r.divisionMap[id]
}

func (r *DivisionRepository) GetAll() []*galaxy.Division {
	divisions := slices.Collect(maps.Values(r.divisionMap))

	slices.SortFunc(divisions, func(a, b *galaxy.Division) int {
		return strings.Compare(a.ID, b.ID)
	})

	return divisions
}

func (r *DivisionRepository) Upsert(division *galaxy.Division) {
	r.divisionMap[division.ID] = division
}

func (r *DivisionRepository) Delete(id string) {
	delete(r.divisionMap, id)
}

func (r *DivisionRepository) ResetData() {
	r.divisionMap = make(map[string]*galaxy.Division)
}
