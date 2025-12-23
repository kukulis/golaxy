class Shot {
    /**
     * Source ship ID that fired the shot
     * @type {string|number}
     */
    source = '';

    /**
     * Destination ship ID that is the target of the shot
     * @type {string|number}
     */
    destination = '';

    /**
     * Whether the shot hit the target (true) or missed (false)
     * @type {boolean}
     */
    result = false;
}
