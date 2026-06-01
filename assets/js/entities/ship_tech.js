export class ShipTech {
    /** @type {number} */
    attack = 0;

    /** @type {number} */
    guns = 0;

    /** @type {number} */
    defense = 0;

    /** @type {number} */
    speed = 0;

    /** @type {number} */
    cargo_capacity = 0;

    /** @type {number} */
    mass = 0;

    /**
     * @param {Object} data
     * @returns {ShipTech}
     */
    updateFromDTO(data) {
        this.attack         = data.attack         ?? 0;
        this.guns           = data.guns           ?? 0;
        this.defense        = data.defense        ?? 0;
        this.speed          = data.speed          ?? 0;
        this.cargo_capacity = data.cargo_capacity ?? 0;
        this.mass           = data.mass           ?? 0;
        return this;
    }
}
