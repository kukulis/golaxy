import { Ship } from './ship.js';

export class ShipGroup {
    /**
     * The ship template for this group
     * @type {Ship}
     */
    ship = null;

    /**
     * Number of ships in this group
     * @type {number}
     */
    amount = 0;

    /**
     * List of ships in this group
     * @type {Ship[]}
     */
    shipList = [];

    /**
     * Groups ships by their tech key
     * @param {Ship[]} ships - Array of ships to group
     * @return {ShipGroup[]}
     */
    static groupShips(ships) {
        // TODO: implement
        return [];
    }
}
