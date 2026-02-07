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
        const groupMap = new Map();

        for (const ship of ships) {
            const key = ship.buildTechKey();

            if (groupMap.has(key)) {
                const group = groupMap.get(key);
                group.shipList.push(ship);
                group.amount++;
            } else {
                const group = new ShipGroup();
                group.ship = ship;
                group.shipList = [ship];
                group.amount = 1;
                groupMap.set(key, group);
            }
        }

        return Array.from(groupMap.values());
    }
}
