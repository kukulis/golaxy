export class Shot {
    /**
     * Source ship ID that fired the shot
     * @type {string|number}
     */
    source = '';

    /**
     * Destination ship ID that is the target of the shot
     * @type {string|number}
     */
    destination = '';

    /**
     * Whether the shot hit the target (true) or missed (false)
     * @type {boolean}
     */
    result = false;

    /**
     * Not received from api, but assigned later.
     * @type {Ship}
     */
    sourceShip = null;

    /**
     * Not received from api, but assigned later.
     * @type {Ship}
     */
    destinationShip = null;

    /**
     *
     * @type {null|SVGElement}
     */
    svgElement = null;

    /**
     * @return {SVGElement}
     */
    buildSvg() {
        const SVG_NS = 'http://www.w3.org/2000/svg';

        // Create group to hold all shot elements
        const group = document.createElementNS(SVG_NS, 'g');
        group.setAttribute('class', 'shot');

        if (!this.sourceShip || !this.destinationShip) {
            console.error('Could not find ships for shot:', this);

            return group;
        }

        // Get coordinates from ships (nose for source, center for destination)
        const x1 = this.sourceShip.noseX;
        const y1 = this.sourceShip.noseY;
        const x2 = this.destinationShip.centerX;
        const y2 = this.destinationShip.centerY;

        // Create line from source to destination
        const line = document.createElementNS(SVG_NS, 'line');
        line.setAttribute('x1', x1);
        line.setAttribute('y1', y1);
        line.setAttribute('x2', x2);
        line.setAttribute('y2', y2);
        line.setAttribute('stroke', 'yellow');
        line.setAttribute('stroke-width', 2);
        group.appendChild(line);

        // Calculate ship size for indicators
        const scale = 3;
        const dp = this.destinationShip.buildDrawParams();
        const circleRadius = dp.drawMass * scale;

        // Draw hit or miss indicator
        if (this.result) {
            // Hit - draw cross (orange for visibility on both blue and red ships)
            const crossSize = circleRadius;
            const line1 = document.createElementNS(SVG_NS, 'line');
            line1.setAttribute('x1', x2 - crossSize);
            line1.setAttribute('y1', y2 - crossSize);
            line1.setAttribute('x2', x2 + crossSize);
            line1.setAttribute('y2', y2 + crossSize);
            line1.setAttribute('stroke', 'orange');
            line1.setAttribute('stroke-width', 3);

            const line2 = document.createElementNS(SVG_NS, 'line');
            line2.setAttribute('x1', x2 + crossSize);
            line2.setAttribute('y1', y2 - crossSize);
            line2.setAttribute('x2', x2 - crossSize);
            line2.setAttribute('y2', y2 + crossSize);
            line2.setAttribute('stroke', 'orange');
            line2.setAttribute('stroke-width', 3);

            group.appendChild(line1);
            group.appendChild(line2);
        } else {
            // Miss - draw circle matching ship size
            const circle = document.createElementNS(SVG_NS, 'circle');
            circle.setAttribute('cx', x2);
            circle.setAttribute('cy', y2);
            circle.setAttribute('r', circleRadius+10);
            circle.setAttribute('fill', 'none');
            circle.setAttribute('stroke', 'white');
            circle.setAttribute('stroke-width', 2);
            group.appendChild(circle);
        }

        this.svgElement = group;

        return group;
    }

    /**
     * Updates shot properties from DTO data
     * @param {Object} data - Shot data from API
     */
    updateFromDTO(data) {
        this.source = data.source;
        this.destination = data.destination;
        this.result = data.result;

        return this;
    }
}
