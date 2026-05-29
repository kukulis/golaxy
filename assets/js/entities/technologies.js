export class Technologies {
    /** @type {number} */
    attack = 1;

    /** @type {number} */
    defense = 1;

    /** @type {number} */
    engine = 1;

    /** @type {number} */
    cargo = 1;

    /**
     * @param {Object} data
     * @returns {Technologies}
     */
    updateFromDTO(data) {
        this.attack = data.attack ?? 1;
        this.defense = data.defense ?? 1;
        this.engine = data.engine ?? 1;
        this.cargo = data.cargo ?? 1;
        return this;
    }
}
