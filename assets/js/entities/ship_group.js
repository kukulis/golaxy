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
     * Scale used for drawing
     * @type {number}
     */
    static SCALE = 3;

    /**
     * Font size for label
     * @type {number}
     */
    static LABEL_FONT_SIZE = 12;

    /**
     * Padding between label and ship
     * @type {number}
     */
    static LABEL_PADDING = 5;

    /**
     * Calculates the total height of the ship group drawing
     * @return {number}
     */
    calculateHeight() {
        const dp = this.ship.buildDrawParams();
        const circleRadius = dp.drawMass * ShipGroup.SCALE;
        const labelHeight = ShipGroup.LABEL_FONT_SIZE + ShipGroup.LABEL_PADDING;
        // Total height: label + padding + ship diameter
        return labelHeight + (circleRadius * 2);
    }

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
        const scale = ShipGroup.SCALE;
        const dp = this.ship.buildDrawParams();

        // Calculate ship center offset (circle is at x + trapezeHeight + circleRadius)
        const circleRadius = dp.drawMass * scale;
        const trapezeHeight = dp.drawSpeed * scale;
        const shipCenterOffset = trapezeHeight + circleRadius;

        // Adjust label x based on side (ship flips 180 for side 'b')
        const labelX = side === 'b' ? x - shipCenterOffset : x + shipCenterOffset;
        const labelOffsetY = circleRadius + ShipGroup.LABEL_FONT_SIZE + ShipGroup.LABEL_PADDING;

        // Add label above ship
        const label = document.createElementNS(SVG_NS, 'text');
        label.setAttribute('x', labelX);
        label.setAttribute('y', y - labelOffsetY);
        label.setAttribute('fill', color);
        label.setAttribute('font-size', ShipGroup.LABEL_FONT_SIZE);
        label.setAttribute('text-anchor', 'middle');
        label.textContent = `${this.ship.name} [${this.amount}]`;
        group.appendChild(label);

        // Draw ship using ShipDraw
        const shipDraw = new ShipDraw(scale, color, 'yellow');
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
