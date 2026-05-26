export class Race {
    /**
     * Race identifier
     * @type {string}
     */
    id = '';

    /**
     * Race name
     * @type {string}
     */
    name = '';

    /**
     * @type {string}
     */
    role = '';

    // TODO: temporal property
    /**
     * @type {string}
     */
    token = '';

    /**
     * @param {Object} data
     * @returns {Race}
     */
    updateFromDTO(data) {
        this.id = data.ID;
        this.name = data.Name;
        this.role = data.Role;
        this.token = data.Token;
        return this;
    }
}
