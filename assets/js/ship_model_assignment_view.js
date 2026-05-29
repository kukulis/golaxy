import {NewE, NewT} from '/assets/js/helper.js'
import {ApiClient} from './api.js'
import {Dispatcher} from './dispatcher.js'

export class ShipModelAssignmentView {

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
    assignmentId = null

    /**
     * @type {string}
     */
    fleetBuildId = null

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
     * @param {string} assignmentId
     */
    setAssignmentId(assignmentId) {
        this.assignmentId = assignmentId
    }

    /**
     * @param {string} fleetBuildId
     */
    setFleetBuildId(fleetBuildId) {
        this.fleetBuildId = fleetBuildId
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
            const assignment = await this.apiClient.getShipModelAssignment(this.assignmentId)

            // Menu used instead.
            // const backLink = NewE('a')
            // backLink.href = `/fleet-build/${this.fleetBuildId}/main.html`
            // backLink.className = 'btn'
            // backLink.appendChild(NewT('← Fleet Build'))
            // container.appendChild(backLink)

            const title = NewE('h2')
            title.appendChild(NewT(assignment.shipModel.name))
            container.appendChild(title)

            const table = NewE('table')
            for (const [label, value] of [
                ['Assignment ID', assignment.id],
                ['Fleet Build', assignment.fleet_build_id],
                ['Ship Model ID', assignment.ship_model_id],
                ['Amount', assignment.amount],
                ['Result Mass', assignment.result_mass],
                ['Guns', assignment.shipModel.guns],
                ['Gun Mass', assignment.shipModel.one_gun_mass],
                ['Defense Mass', assignment.shipModel.defense_mass],
                ['Engine Mass', assignment.shipModel.engine_mass],
                ['Cargo Mass', assignment.shipModel.cargo_mass],
            ]) {
                const tr = NewE('tr')
                const tdLabel = NewE('td')
                tdLabel.appendChild(NewT(label))
                tr.appendChild(tdLabel)
                const tdValue = NewE('td')
                tdValue.appendChild(NewT(String(value)))
                tr.appendChild(tdValue)
                table.appendChild(tr)
            }
            container.appendChild(table)
        } catch (e) {
            this.dispatcher.dispatch('displayError', [e.message, true])
        }

        return container
    }
}
