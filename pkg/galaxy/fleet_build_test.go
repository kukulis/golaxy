package galaxy

import (
	"reflect"
	"testing"
)

func TestCalculateShipTech_WithAllResources(t *testing.T) {
	// Given: A fleet build with resources in all areas
	fleetBuild := &FleetBuild{
		AttackResources:  100, // +1 tech => 2.0
		DefenseResources: 50,  // +0.5 tech => 1.5
		EngineResources:  200, // +2 tech => 3.0
		CargoResources:   150, // +1.5 tech => 2.5 (not used in ShipTech currently)
	}

	// And: A ship model
	shipModel := &ShipModel{
		Guns:        3,
		OneGunMass:  5,
		DefenseMass: 15,
		EngineMass:  30,
	}
	// Total mass = 3*5 + 15 + 30 = 60

	// When: Calculating ship tech
	result := fleetBuild.CalculateShipTech(shipModel)

	// Then: All techs should be enhanced
	// speed = EngineMass * Engine / mass = 30 * 3.0 / 60 = 1.5
	// defense = DefenseMass * Defense / sqrt(mass) = 15 * 1.5 / sqrt(60) = 22.5 / 7.745... = 2.905...
	// attack = OneGunMass * Attack = 5 * 2.0 = 10
	expected := ShipTech{
		Guns:    3,
		Attack:  10,                       // 5 * 2.0
		Defense: 22.5 / 7.745966692414834, // 15 * 1.5 / sqrt(60)
		Speed:   1.5,                      // 30 * 3.0 / 60
		Mass:    60,
	}

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("CalculateShipTech with all resources:\nexpected %+v\ngot      %+v", expected, result)
	}
}

func TestCalculateShipTech_WithAllResources2(t *testing.T) {
	fleetBuild := &FleetBuild{
		AttackResources:  100, // +1 tech => 2.0
		DefenseResources: 50,  // +0.5 tech => 1.5
		EngineResources:  200, // +2 tech => 3.0
		CargoResources:   150, // +1.5 tech => 2.5
	}

	shipModel := &ShipModel{
		Guns:        4,
		OneGunMass:  4,
		DefenseMass: 16,
		EngineMass:  32,
	}

	result := fleetBuild.CalculateShipTech(shipModel)

	expected := ShipTech{
		Guns:    4,
		Attack:  8,
		Defense: 3,
		Speed:   1.5,
		Mass:    64,
	}

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("CalculateShipTech with all resources:\nexpected %+v\ngot      %+v", expected, result)
	}
}

func TestCalculateShipTech_Simple(t *testing.T) {
	fleetBuild := &FleetBuild{
		AttackResources:  0,
		DefenseResources: 0,
		EngineResources:  0,
		CargoResources:   0,
	}

	shipModel := &ShipModel{
		Guns:        0,
		OneGunMass:  0,
		DefenseMass: 0,
		EngineMass:  1,
	}

	result := fleetBuild.CalculateShipTech(shipModel)

	expected := ShipTech{
		Guns:    0,
		Attack:  0,
		Defense: 0,
		Speed:   1,
		Mass:    1,
	}

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("CalculateShipTech with all resources:\nexpected %+v\ngot      %+v", expected, result)
	}
}
func TestCalculateShipTech_Defensive(t *testing.T) {
	fleetBuild := &FleetBuild{
		AttackResources:  0,
		DefenseResources: 0,
		EngineResources:  0,
		CargoResources:   0,
	}

	shipModel := &ShipModel{
		Guns:        0,
		OneGunMass:  0,
		DefenseMass: 2,
		EngineMass:  2,
	}

	result := fleetBuild.CalculateShipTech(shipModel)

	expected := ShipTech{
		Guns:    0,
		Attack:  0,
		Defense: 1,
		Speed:   0.5,
		Mass:    4,
	}

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("CalculateShipTech with all resources:\nexpected %+v\ngot      %+v", expected, result)
	}
}

func TestCalculateShipTech_VeryDefensive(t *testing.T) {
	fleetBuild := &FleetBuild{
		AttackResources:  0,
		DefenseResources: 0,
		EngineResources:  0,
		CargoResources:   0,
	}

	shipModel := &ShipModel{
		Guns:        1,
		OneGunMass:  10,
		DefenseMass: 40,
		EngineMass:  50,
	}

	result := fleetBuild.CalculateShipTech(shipModel)

	expected := ShipTech{
		Guns:    1,
		Attack:  10,
		Defense: 4,
		Speed:   0.5,
		Mass:    100,
	}

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("CalculateShipTech with all resources:\nexpected %+v\ngot      %+v", expected, result)
	}
}

func TestCalculateShipTech_UltraDefensive(t *testing.T) {
	fleetBuild := &FleetBuild{
		AttackResources:  0,
		DefenseResources: 0,
		EngineResources:  0,
		CargoResources:   0,
	}

	shipModel := &ShipModel{
		Guns:        1,
		OneGunMass:  30,
		DefenseMass: 870,
		EngineMass:  0,
	}

	result := fleetBuild.CalculateShipTech(shipModel)

	expected := ShipTech{
		Guns:    1,
		Attack:  30,
		Defense: 29,
		Speed:   0,
		Mass:    900,
	}

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("CalculateShipTech with all resources:\nexpected %+v\ngot      %+v", expected, result)
	}
}

func TestCalculateShipTech_UltraDefensiveFlying(t *testing.T) {
	fleetBuild := &FleetBuild{
		AttackResources:  0,
		DefenseResources: 0,
		EngineResources:  0,
		CargoResources:   0,
	}

	shipModel := &ShipModel{
		Guns:        1,
		OneGunMass:  30,
		DefenseMass: 420,
		EngineMass:  450,
	}

	result := fleetBuild.CalculateShipTech(shipModel)

	expected := ShipTech{
		Guns:    1,
		Attack:  30,
		Defense: 14,
		Speed:   0.5,
		Mass:    900,
	}

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("CalculateShipTech with all resources:\nexpected %+v\ngot      %+v", expected, result)
	}
}

func TestCalculateShipTech_Perfo(t *testing.T) {
	fleetBuild := &FleetBuild{
		AttackResources:  0,
		DefenseResources: 0,
		EngineResources:  0,
		CargoResources:   0,
	}

	shipModel := &ShipModel{
		Guns:        10,
		OneGunMass:  1,
		DefenseMass: 40,
		EngineMass:  50,
	}

	result := fleetBuild.CalculateShipTech(shipModel)

	expected := ShipTech{
		Guns:    10,
		Attack:  1,
		Defense: 4,
		Speed:   0.5,
		Mass:    100,
	}

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("CalculateShipTech with all resources:\nexpected %+v\ngot      %+v", expected, result)
	}
}

func TestCalculateShipTech_GlassPerfo(t *testing.T) {
	fleetBuild := &FleetBuild{
		AttackResources:  0,
		DefenseResources: 0,
		EngineResources:  0,
		CargoResources:   0,
	}

	shipModel := &ShipModel{
		Guns:        30,
		OneGunMass:  1,
		DefenseMass: 20,
		EngineMass:  50,
	}

	result := fleetBuild.CalculateShipTech(shipModel)

	expected := ShipTech{
		Guns:    30,
		Attack:  1,
		Defense: 2,
		Speed:   0.5,
		Mass:    100,
	}

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("CalculateShipTech with all resources:\nexpected %+v\ngot      %+v", expected, result)
	}
}
