class Ship {
    /**
     * Ship identifier
     * @type {string}
     */
    id = '';

    /**
     * Ship technology specifications
     * @type {{attack: number, guns: number, defense: number, speed: number, cargo_capacity: number, mass: number}}
     */
    tech = {
        attack: 0,
        guns: 0,
        defense: 0,
        speed: 0,
        cargo_capacity: 0,
        mass: 0
    };

    /**
     * Whether the ship has been destroyed in battle
     * @type {boolean}
     */
    destroyed = false;

    /**
     * Ship name
     * @type {string}
     */
    name = '';

    /**
     * Ship owner identifier
     * @type {string}
     */
    owner = '';

    /**
     * X coordinate for battle rendering (not received from API)
     * @type {number}
     */
    battleX = 0;

    /**
     * Y coordinate for battle rendering (not received from API)
     * @type {number}
     */
    battleY = 0;

    /**
     *
     * @type {null|SVGElement}
     */
    svgElement = null;


    creteShipSvg(color, side) {

        let x = this.battleX;
        let y = this.battleY;

        // Create ship group
        const group = document.createElementNS(SVG_NS, 'g');
        group.setAttribute('class', 'ship');

        // // Draw ship as a circle
        const circle = document.createElementNS(SVG_NS, "circle");

        const circleRadius = 20;

        circle.setAttribute("cx", x);  // center x
        circle.setAttribute("cy", y);  // center y
        circle.setAttribute("r", circleRadius);    // radius
        circle.setAttribute("fill", color);
        circle.setAttribute('fill', this.destroyed ? 'gray' : color);
        circle.setAttribute('stroke', 'white');
        circle.setAttribute('stroke-width', 2);




        // Add ship label
        // TODO remake using css class
        const text = document.createElementNS(SVG_NS, 'text');
        text.setAttribute('x', x);
        text.setAttribute('y', y - 20);
        text.setAttribute('text-anchor', 'middle');
        text.setAttribute('fill', 'white');
        text.setAttribute('font-size', '12');
        text.textContent = this.name;

        if (this.tech.speed > 0) {
            group.appendChild(this.creteMotorSvg(color, side));
        }

        group.appendChild(circle);



        group.appendChild(text);



        // Add click handler
        group.addEventListener('click', () => this.handleShipClick());

        // Add hover effect
        group.style.cursor = 'pointer';

        this.svgElement = group;

        return group;
    }

    handleShipClick() {
        console.log('Ship clicked:', this);
        const info = `Ship: ${this.name}\nID: ${this.id}\nSpeed: ${this.tech.speed}\nAttack: ${this.tech.attack}\nDefense: ${this.tech.defense}\nDestroyed: ${this.destroyed}`;
        alert(info);
    }

    /**
     * Updates ship properties from DTO data
     * @param {Object} data - Ship data from API
     */
    updateFromDTO(data) {
        this.id = data.id;
        this.tech = {
            attack: data.tech.attack,
            guns: data.tech.guns,
            defense: data.tech.defense,
            speed: data.tech.speed,
            cargo_capacity: data.tech.cargo_capacity,
            mass: data.tech.mass
        };
        this.destroyed = data.destroyed;
        this.name = data.name;
        this.owner = data.owner;

        return this;
    }

    creteMotorSvg(color, side) {

        let x = this.battleX;
        let y = this.battleY;


        // Draw motor as a triangle
        const triangle = document.createElementNS(SVG_NS, 'polygon');

        // Create triangle points based on side
        let points;
        if (side === 'a') {
            // Right-pointing triangle for side A
            x = x - 20
            points = `${x - 10},${y - 10} ${x + 15},${y} ${x - 10},${y + 10}`;
        } else {
            // Left-pointing triangle for side B
            x = x + 20
            points = `${x + 10},${y - 10} ${x - 15},${y} ${x + 10},${y + 10}`;
        }

        triangle.setAttribute('points', points);
        triangle.setAttribute('fill', this.destroyed ? 'gray' : color);
        triangle.setAttribute('stroke', 'white');
        triangle.setAttribute('stroke-width', 2);

        return triangle;
    }
}
