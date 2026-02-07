import { Ship } from '../entities/ship.js';
import { ShipGroup } from '../entities/ship_group.js';
import { TestRunner } from './test_runner.js';

function createShip(id, attack, guns, defense, speed, cargoCapacity, mass) {
    const ship = new Ship();
    ship.id = id;
    ship.tech = {
        attack: attack,
        guns: guns,
        defense: defense,
        speed: speed,
        cargo_capacity: cargoCapacity,
        mass: mass
    };
    return ship;
}

// Test: empty array returns empty array
const emptyResult = ShipGroup.groupShips([]);
TestRunner.assert(Array.isArray(emptyResult), 'groupShips returns an array');
TestRunner.assertEquals(emptyResult.length, 0, 'empty input returns empty array');

// Test: single ship returns single group
const ship1 = createShip('s1', 1, 2, 3, 4, 5, 6);
const singleResult = ShipGroup.groupShips([ship1]);
TestRunner.assertEquals(singleResult.length, 1, 'single ship returns one group');
if (singleResult.length > 0) {
    TestRunner.assertEquals(singleResult[0].amount, 1, 'single ship group has amount 1');
    TestRunner.assertEquals(singleResult[0].shipList.length, 1, 'single ship group has one ship in list');
}

// Test: two ships with same tech are grouped together
const ship2 = createShip('s2', 1, 2, 3, 4, 5, 6);
const sameTechResult = ShipGroup.groupShips([ship1, ship2]);
TestRunner.assertEquals(sameTechResult.length, 1, 'two ships with same tech return one group');
if (sameTechResult.length > 0) {
    TestRunner.assertEquals(sameTechResult[0].amount, 2, 'group has amount 2');
    TestRunner.assertEquals(sameTechResult[0].shipList.length, 2, 'group has two ships in list');
}

// Test: two ships with different tech are in separate groups
const ship3 = createShip('s3', 10, 2, 3, 4, 5, 6);
const diffTechResult = ShipGroup.groupShips([ship1, ship3]);
TestRunner.assertEquals(diffTechResult.length, 2, 'two ships with different tech return two groups');

// Print results
const success = TestRunner.printResults();
process.exit(success ? 0 : 1);
