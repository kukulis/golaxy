import { Ship } from './ship.js';

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
