export class ChallengesFilter {
    /** @type {string} */
    challenger_id = '';

    /** @type {string} */
    challengee_id = '';

    /** @type {boolean} */
    ready_a = false;

    /** @type {boolean} */
    ready_b = false;

    /** @type {string} */
    status = '';

    /** @type {string} */
    division_id = '';

    toQueryString() {
        const params = new URLSearchParams();
        if (this.challenger_id) params.set('challenger_id', this.challenger_id);
        if (this.challengee_id) params.set('challengee_id', this.challengee_id);
        if (this.ready_a) params.set('ready_a', '1');
        if (this.ready_b) params.set('ready_b', '1');
        if (this.status) params.set('status', this.status);
        if (this.division_id) params.set('division_id', this.division_id);
        return params.size > 0 ? '?' + params.toString() : '';
    }
}
