import { test } from 'node:test'
import assert from 'node:assert/strict'
import { ShipTech } from '../../../assets/js/entities/ship_tech.js'

test('updateFromDTO maps all fields', () => {
    const tech = new ShipTech()
    tech.updateFromDTO({
        attack:         4.0,
        guns:           2,
        defense:        3.0,
        speed:          0.5,
        cargo_capacity: 10.0,
        mass:           8.0,
    })
    assert.equal(tech.attack,         4.0)
    assert.equal(tech.guns,           2)
    assert.equal(tech.defense,        3.0)
    assert.equal(tech.speed,          0.5)
    assert.equal(tech.cargo_capacity, 10.0)
    assert.equal(tech.mass,           8.0)
})

test('updateFromDTO defaults missing fields to 0', () => {
    const tech = new ShipTech()
    tech.updateFromDTO({})
    assert.equal(tech.attack,         0)
    assert.equal(tech.guns,           0)
    assert.equal(tech.defense,        0)
    assert.equal(tech.speed,          0)
    assert.equal(tech.cargo_capacity, 0)
    assert.equal(tech.mass,           0)
})

test('updateFromDTO returns this for chaining', () => {
    const tech = new ShipTech()
    const result = tech.updateFromDTO({})
    assert.equal(result, tech)
})
