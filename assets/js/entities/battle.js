class Battle {
    /**
     * Battle identifier
     * @type {string}
     */
    id = '';

    /**
     * First fleet participating in the battle
     * @type {Fleet|null}
     */
    side_a = null;

    /**
     * Second fleet participating in the battle
     * @type {Fleet|null}
     */
    side_b = null;

    /**
     * Array of all shots fired during the battle
     * @type {Shot[]}
     */
    shots = [];
}
