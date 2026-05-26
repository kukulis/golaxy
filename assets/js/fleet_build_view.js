import {NewE, NewT} from '/assets/js/helper.js'
import {ApiClient} from './api.js'
import {Dispatcher} from './dispatcher.js'
import {createFleetBuildStatisticsTable} from './fleet_build_table.js'
import {createDivisionTable} from './division_table.js'
import {createFleetBuildShipModelTable} from './fleet_build_ship_model_table.js'

export default class FleetBuildView {

    /**
     * @type {ApiClient}
     */
    apiClient = null

    /**
     * @type {Dispatcher}
     */
    dispatcher = null

    /**
     * @type {string}
     */
    buildId = null

    /**
     * @type {string}
     */
    divisionId = null

    /**
     * @param {ApiClient} apiClient
     * @param {Dispatcher} dispatcher
     */
    constructor(apiClient, dispatcher) {
        this.apiClient = apiClient
        this.dispatcher = dispatcher
    }

    /**
     * @param {string} buildId
     */
    setBuildId(buildId) {
        this.buildId = buildId
    }

    /**
     * @param {string} divisionId
     */
    setDivisionId(divisionId) {
        this.divisionId = divisionId
    }

    /**
     * @return {HTMLElement}
     */
    async generateView() {
        const container = NewE('div')

        try {
            const [b, division, assignments] = await Promise.all([
                this.apiClient.getFleetBuild(this.buildId),
                this.apiClient.getDivision(this.divisionId),
                this.apiClient.getFleetBuildShipModels(this.buildId),
            ])

            const leftCol = NewE('div')

            const fleetBuildTitle = NewE('h2')
            fleetBuildTitle.appendChild(NewT('Fleet Build'))
            leftCol.appendChild(fleetBuildTitle)
            leftCol.appendChild(createFleetBuildStatisticsTable(b))

            const divisionTitle = NewE('h2')
            divisionTitle.appendChild(NewT('Division'))
            leftCol.appendChild(divisionTitle)
            leftCol.appendChild(createDivisionTable(division))

            const rightCol = NewE('div')

            const shipModelsTitle = NewE('h2')
            shipModelsTitle.appendChild(NewT('Assigned Ship Models'))
            rightCol.appendChild(shipModelsTitle)
            rightCol.appendChild(createFleetBuildShipModelTable(assignments))

            container.className = 'two-col-layout'
            container.appendChild(leftCol)
            container.appendChild(rightCol)
        } catch (e) {
            this.dispatcher.dispatch('displayError', [e.message, true])
        }

        return container
    }
}
