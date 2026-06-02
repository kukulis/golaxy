export class Challenge {
    /** @type {string} */
    id = '';

    /** @type {string} */
    division_id = '';

    /** @type {string} */
    challenger_race_id = '';

    /** @type {string} */
    challengee_race_id = '';

    /** @type {string} */
    fleet_build_a_id = '';

    /** @type {string} */
    fleet_build_b_id = '';

    /** @type {boolean} */
    ready_a = false;

    /** @type {boolean} */
    ready_b = false;

    /** @type {string|null} */
    created_at = null;

    /** @type {string|null} */
    accepted_a_at = null;

    /** @type {string|null} */
    accepted_b_at = null;

    /** @type {string} */
    status = '';

    /** @type {string|null} */
    executed_at = null;

    /** @type {string} */
    battle_report_id = '';

    /**
     * @param {Object} data
     * @returns {Challenge}
     */
    updateFromDTO(data) {
        this.id                = data.id                ?? '';
        this.division_id       = data.division_id       ?? '';
        this.challenger_race_id = data.challenger_race_id ?? '';
        this.challengee_race_id = data.challengee_race_id ?? '';
        this.fleet_build_a_id  = data.fleet_build_a_id  ?? '';
        this.fleet_build_b_id  = data.fleet_build_b_id  ?? '';
        this.ready_a           = data.ready_a           ?? false;
        this.ready_b           = data.ready_b           ?? false;
        this.created_at        = data.created_at        ?? null;
        this.accepted_a_at     = data.accepted_a_at     ?? null;
        this.accepted_b_at     = data.accepted_b_at     ?? null;
        this.status            = data.status            ?? '';
        this.executed_at       = data.executed_at       ?? null;
        this.battle_report_id  = data.battle_report_id  ?? '';
        return this;
    }
}
