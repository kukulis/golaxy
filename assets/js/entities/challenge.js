export class Challenge {
    /** @type {string} */
    id = '';

    /** @type {string} */
    division_id = '';

    /** @type {string} */
    challenger_race_id = '';

    /** @type {string} */
    fleet_build_a_id = '';

    /** @type {string} */
    fleet_build_b_id = '';

    /** @type {string} */
    battle_report_id = '';

    /**
     * @param {Object} data
     * @returns {Challenge}
     */
    updateFromDTO(data) {
        this.id                 = data.id                  ?? '';
        this.division_id        = data.division_id         ?? '';
        this.challenger_race_id = data.challenger_race_id  ?? '';
        this.fleet_build_a_id   = data.fleet_build_a_id    ?? '';
        this.fleet_build_b_id   = data.fleet_build_b_id    ?? '';
        this.battle_report_id   = data.battle_report_id    ?? '';
        return this;
    }
}
