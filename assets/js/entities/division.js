export class Division {
    /** @type {string} */
    id = '';

    /** @type {number} */
    resources_amount = 0;

    /** @type {number} */
    tech_attack = 0;

    /** @type {number} */
    tech_defense = 0;

    /** @type {number} */
    tech_engines = 0;

    /** @type {number} */
    tech_cargo = 0;

    /**
     * @param {Object} data
     * @returns {Division}
     */
    updateFromDTO(data) {
        this.id = data.id;
        this.resources_amount = data.resources_amount ?? 0;
        this.tech_attack = data.tech_attack ?? 0;
        this.tech_defense = data.tech_defense ?? 0;
        this.tech_engines = data.tech_engines ?? 0;
        this.tech_cargo = data.tech_cargo ?? 0;
        return this;
    }
}
