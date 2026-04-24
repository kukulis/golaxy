import {ShipModel} from './ship_model.js';

export class FleetBuildShipModel {
    /** @type {string} */
    fleet_build_id = '';

    /** @type {string} */
    ship_model_id = '';

    /** @type {number} */
    amount = 0;

    /** @type {number} */
    result_mass = 0;

    /** @type {ShipModel} */
    shipModel = new ShipModel();

    /**
     * @param {Object} data
     * @returns {FleetBuildShipModel}
     */
    updateFromDTO(data) {
        this.fleet_build_id = data.fleet_build_id ?? '';
        this.ship_model_id = data.ship_model_id ?? '';
        this.amount = data.amount ?? 0;
        this.result_mass = data.result_mass ?? 0;
        if (data.shipModel) {
            this.shipModel = (new ShipModel()).updateFromDTO(data.shipModel);
        }
        return this;
    }
}
