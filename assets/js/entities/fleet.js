import { Ship } from './ship.js';
import { ShipGroup } from './ship_group.js';

export class Fleet {
    /**
     * Array of ships in this fleet
     * @type {Ship[]}
     */
    ships = [];

    /**
     * Fleet owner identifier
     * @type {string}
     */
    owner = '';

    /**
     *
     * @type {Map<string, Ship>}
     */
    shipsMap = new Map();

    /**
     * Ships grouped by tech key (not received from backend)
     * @type {Map<string, ShipGroup>}
     */
    shipGroupMap = new Map();

    /**
     *
     * @param ships {Ship[]}
     */
    fillShipsMap(ships) {
        for ( let ship of ships ) {
            this.shipsMap.set(ship.id, ship)
        }
    }

    findShip(shipId) {
        return this.shipsMap.get(shipId)
    }

    /**
     * Fills shipGroupMap from the ships array
     */
    fillShipGroupMap() {
        this.shipGroupMap.clear();
        const groups = ShipGroup.groupShips(this.ships);
        for (const group of groups) {
            const key = group.ship.buildTechKey();
            this.shipGroupMap.set(key, group);
        }
    }

    /**
     * Finds a ship group that contains the given ship
     * @param {Ship} ship
     * @return {ShipGroup|undefined}
     */
    findShipGroup(ship) {
        const key = ship.buildTechKey();
        return this.shipGroupMap.get(key);
    }

    /**
     * Notifies that a ship was destroyed, decreases group amount
     * @param {Ship} ship
     */
    notifyDestroyed(ship) {
        if (!ship.destroyed) {
            return;
        }
        const group = this.findShipGroup(ship);
        if (group) {
            group.amount--;
        }
    }

    /**
     * Updates fleet properties from DTO data
     * @param {Object} data - Fleet data from API
     */
    updateFromDTO(data) {
        this.ships = data.ships.map(shipData => {
            const ship = new Ship();
            ship.updateFromDTO(shipData);
            return ship;
        });
        this.fillShipsMap(this.ships);
        this.owner = data.owner;

        return this;
    }
}
