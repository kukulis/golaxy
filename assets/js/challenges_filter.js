export class ChallengesFilter {
    /** @type {string} */
    challenger_id = '';

    /** @type {string} */
    challengee_id = '';

    /** @type {boolean} */
    accepted_a = false;

    /** @type {boolean} */
    accepted_b = false;

    toQueryString() {
        const params = new URLSearchParams();
        if (this.challenger_id) params.set('challenger_id', this.challenger_id);
        if (this.challengee_id) params.set('challengee_id', this.challengee_id);
        if (this.accepted_a) params.set('accepted_a', '1');
        if (this.accepted_b) params.set('accepted_b', '1');
        return params.size > 0 ? '?' + params.toString() : '';
    }
}
