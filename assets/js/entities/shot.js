class Shot {
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

        // Get coordinates from ships
        const x1 = this.sourceShip.battleX;
        const y1 = this.sourceShip.battleY;
        const x2 = this.destinationShip.battleX;
        const y2 = this.destinationShip.battleY;

        // Create line from source to destination
        const line = document.createElementNS(SVG_NS, 'line');
        line.setAttribute('x1', x1);
        line.setAttribute('y1', y1);
        line.setAttribute('x2', x2);
        line.setAttribute('y2', y2);
        line.setAttribute('stroke', 'yellow');
        line.setAttribute('stroke-width', 2);
        group.appendChild(line);

        // Draw hit or miss indicator
        if (this.result) {
            // Hit - draw cross
            const line1 = document.createElementNS(SVG_NS, 'line');
            line1.setAttribute('x1', x2 - 15);
            line1.setAttribute('y1', y2 - 15);
            line1.setAttribute('x2', x2 + 15);
            line1.setAttribute('y2', y2 + 15);
            line1.setAttribute('stroke', 'red');
            line1.setAttribute('stroke-width', 3);

            const line2 = document.createElementNS(SVG_NS, 'line');
            line2.setAttribute('x1', x2 + 15);
            line2.setAttribute('y1', y2 - 15);
            line2.setAttribute('x2', x2 - 15);
            line2.setAttribute('y2', y2 + 15);
            line2.setAttribute('stroke', 'red');
            line2.setAttribute('stroke-width', 3);

            group.appendChild(line1);
            group.appendChild(line2);
        } else {
            // Miss - draw circle
            const circle = document.createElementNS(SVG_NS, 'circle');
            circle.setAttribute('cx', x2);
            circle.setAttribute('cy', y2);
            circle.setAttribute('r', 20);
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
