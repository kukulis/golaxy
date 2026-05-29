import {Technologies} from './technologies.js';

export class FleetBuildStatistics {
    /** @type {number} */
    max_resources = 0;

    /** @type {number} */
    used_resources = 0;

    /** @type {number} */
    used_resources_for_ships = 0;

    /** @type {number} */
    used_resources_for_technologies = 0;

    /** @type {number} */
    remaining_resources = 0;

    /** @type {number} */
    exceeding_resources = 0;

    /** @type {Technologies|null} */
    technologies = null;

    /**
     * @param {Object} data
     * @returns {FleetBuildStatistics}
     */
    updateFromDTO(data) {
        this.max_resources = data.max_resources ?? 0;
        this.used_resources = data.used_resources ?? 0;
        this.used_resources_for_ships = data.used_resources_for_ships ?? 0;
        this.used_resources_for_technologies = data.used_resources_for_technologies ?? 0;
        this.remaining_resources = data.remaining_resources ?? 0;
        this.exceeding_resources = data.exceeding_resources ?? 0;
        this.technologies = data.technologies ? (new Technologies()).updateFromDTO(data.technologies) : null;
        return this;
    }
}
