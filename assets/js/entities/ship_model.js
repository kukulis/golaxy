export class ShipModel {
    /** @type {string} */
    id = '';

    /** @type {string} */
    name = '';

    /** @type {string} */
    owner_id = '';

    /** @type {number} */
    guns = 0;

    /** @type {number} */
    one_gun_mass = 0;

    /** @type {number} */
    defense_mass = 0;

    /** @type {number} */
    engine_mass = 0;

    /** @type {number} */
    cargo_mass = 0;

    /**
     * @param {Object} data
     * @returns {ShipModel}
     */
    updateFromDTO(data) {
        this.id = data.id;
        this.name = data.name ?? '';
        this.owner_id = data.owner_id ?? '';
        this.guns = data.guns ?? 0;
        this.one_gun_mass = data.one_gun_mass ?? 0;
        this.defense_mass = data.defense_mass ?? 0;
        this.engine_mass = data.engine_mass ?? 0;
        this.cargo_mass = data.cargo_mass ?? 0;
        return this;
    }
}
