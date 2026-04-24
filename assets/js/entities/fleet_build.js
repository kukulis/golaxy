export class FleetBuild {
    /** @type {string} */
    id = '';

    /** @type {string} */
    division_id = '';

    /** @type {string} */
    race_id = '';

    /** @type {number} */
    attack_resources = 0;

    /** @type {number} */
    defense_resources = 0;

    /** @type {number} */
    engine_resources = 0;

    /** @type {number} */
    cargo_resources = 0;

    /** @type {number} */
    usedResources = 0;

    /**
     * @param {Object} data
     * @returns {FleetBuild}
     */
    updateFromDTO(data) {
        this.id = data.id;
        this.division_id = data.division_id;
        this.race_id = data.race_id ?? '';
        this.attack_resources = data.attack_resources ?? 0;
        this.defense_resources = data.defense_resources ?? 0;
        this.engine_resources = data.engine_resources ?? 0;
        this.cargo_resources = data.cargo_resources ?? 0;
        this.usedResources = data.usedResources ?? 0;
        return this;
    }
}
