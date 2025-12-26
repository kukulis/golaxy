package galaxy

import (
	"errors"
	"glaktika.eu/galaktika/pkg/util"
	"math"
)

type ShipModel struct {
	ID            string
	Name          string
	Guns          int
	OneGunMass    float64
	DefenseMass   float64
	EngineMass    float64
	validateError error
}

func (shipModel *ShipModel) CalculateTotalMass() float64 {
	return float64(shipModel.Guns)*shipModel.OneGunMass +
		shipModel.DefenseMass +
		shipModel.EngineMass
}

func (shipModel *ShipModel) GetValidateError() error {
	return shipModel.validateError
}

func (shipModel *ShipModel) ValidateModel() bool {
	if shipModel.OneGunMass > 0 && shipModel.OneGunMass < 1 {
		shipModel.validateError = errors.New("OneGunMass must be equal to 0 or not lower than 1")
		return false
	}

	if shipModel.DefenseMass > 0 && shipModel.DefenseMass < 1 {
		shipModel.validateError = errors.New("DefenseMass must be equal to 0 or not lower than 1")
		return false
	}

	if shipModel.EngineMass > 0 && shipModel.EngineMass < 1 {
		shipModel.validateError = errors.New("EngineMass must be equal to 0 or not lower than 1")
		return false
	}

	return true
}

func (shipModel *ShipModel) CalculateShipTech(t *Technologies) ShipTech {
	mass := shipModel.CalculateTotalMass()

	speed := shipModel.EngineMass * t.Engine / mass
	defense := shipModel.DefenseMass * t.Defense / math.Sqrt(mass)
	attack := shipModel.OneGunMass * t.Attack

	return ShipTech{
		Guns:    shipModel.Guns,
		Speed:   speed,
		Defense: defense,
		Attack:  attack,
		Mass:    mass,
	}
}

func (shipModel *ShipModel) GenerateShips(t *Technologies, amount int, generator util.IdGenerator, ownerId string) []*Ship {
	shipTech := shipModel.CalculateShipTech(t)

	ships := make([]*Ship, amount)

	for i := 0; i < amount; i++ {
		id := generator.NextId()
		ships[i] = &Ship{
			VersionID: id,
			ID:        id,
			Name:      shipModel.Name,
			Tech:      shipTech,
			Owner:     ownerId,
			Destroyed: false,
		}
	}

	return ships
}
