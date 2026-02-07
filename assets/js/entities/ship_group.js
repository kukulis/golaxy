import { Ship } from './ship.js';
import { ShipDraw } from '../ship_draw.js';

const SVG_NS = 'http://www.w3.org/2000/svg';

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
     * X coordinate for rendering (not received from API)
     * @type {number}
     */
    battleX = 0;

    /**
     * Y coordinate for rendering (not received from API)
     * @type {number}
     */
    battleY = 0;

    /**
     * @type {null|SVGElement}
     */
    svgElement = null;

    /**
     * Creates SVG for the ship group with name and amount label
     * @param {string} color - Ship color
     * @param {string} side - Side ('a' or 'b')
     * @return {SVGElement}
     */
    createGroupSvg(color, side) {
        const group = document.createElementNS(SVG_NS, 'g');
        group.setAttribute('class', 'ship-group');

        const x = this.battleX;
        const y = this.battleY;

        // Add label above ship
        const label = document.createElementNS(SVG_NS, 'text');
        label.setAttribute('x', x);
        label.setAttribute('y', y - 10);
        label.setAttribute('fill', color);
        label.setAttribute('font-size', '12');
        label.textContent = `${this.ship.name} [${this.amount}]`;
        group.appendChild(label);

        // Draw ship using ShipDraw
        const shipDraw = new ShipDraw(3, color, 'yellow');
        const dp = this.ship.buildDrawParams();
        const rotation = side === 'b' ? 180 : 0;
        shipDraw.drawShipRaw(group, x, y, dp.drawMass, dp.drawSpeed, dp.drawGuns, dp.drawAttack, dp.drawDefence, rotation);

        this.svgElement = group;
        return group;
    }

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
