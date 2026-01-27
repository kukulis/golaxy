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


    createShipSvg(color, side) {

        let shipDraw = new ShipDraw(3, color, 'yellow')
        let dp = this.buildDrawParams();

        const group = document.createElementNS(SVG_NS, 'g');
        group.setAttribute('class', 'ship');
        let x = this.battleX;
        let y = this.battleY;
        const rotation = side === 'b' ? 180 : 0;
        shipDraw.drawShipRaw(group, x, y, dp.drawMass, dp.drawSpeed, dp.drawGuns, dp.drawAttack, dp.drawDefence, rotation);

        return group;
    }

    handleShipClick() {
        console.log('Ship clicked:', this);
        const info = `Ship: ${this.name}\nID: ${this.id}\nSpeed: ${this.tech.speed}\nAttack: ${this.tech.attack}\nGuns:${this.tech.guns} \nDefense: ${this.tech.defense}\nDestroyed: ${this.destroyed}`;
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

    buildDrawParams() {
        return new ShipDrawParams(
            Math.log( this.tech.mass ),
            this.tech.speed,
            this.tech.guns,
            this.tech.attack,
            this.tech.defense
        );
    }
}
