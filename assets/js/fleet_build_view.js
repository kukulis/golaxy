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
            const [b, division, assignments, statistics] = await Promise.all([
                this.apiClient.getFleetBuild(this.buildId),
                this.apiClient.getDivision(this.divisionId),
                this.apiClient.getFleetBuildShipModels(this.buildId),
                this.apiClient.getFleetBuildStatistics(this.buildId),
            ])

            const leftCol = NewE('div')

            const editLink = NewE('a')
            editLink.href = `/fleet-build/${this.buildId}/edit.html`
            editLink.className = 'btn'
            editLink.appendChild(NewT('✏ Edit'))
            leftCol.appendChild(editLink)

            const fleetBuildTitle = NewE('h2')
            fleetBuildTitle.appendChild(NewT('Fleet Build'))
            leftCol.appendChild(fleetBuildTitle)
            leftCol.appendChild(createFleetBuildStatisticsTable(b, statistics))

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

    /**
     * @param {string} buildId
     * @return {HTMLElement}
     */
    async generateEditView(buildId) {
        const container = NewE('div')

        try {
            const b = await this.apiClient.getFleetBuild(buildId)

            const readOnlyTable = NewE('table')
            for (const [label, value] of [
                ['ID', b.id],
                ['Division', b.division_id],
                ['Race', b.race_id],
                ['Used Resources', b.usedResources],
            ]) {
                const tr = NewE('tr')
                const tdLabel = NewE('td')
                tdLabel.appendChild(NewT(label))
                tr.appendChild(tdLabel)
                const tdValue = NewE('td')
                tdValue.appendChild(NewT(String(value)))
                tr.appendChild(tdValue)
                readOnlyTable.appendChild(tr)
            }
            container.appendChild(readOnlyTable)

            const editTable = NewE('table')
            const editFields = [
                {label: 'Name', key: 'name', type: 'text', value: b.name},
                {label: 'Attack Resources', key: 'attack_resources', type: 'number', value: b.attack_resources},
                {label: 'Defense Resources', key: 'defense_resources', type: 'number', value: b.defense_resources},
                {label: 'Engine Resources', key: 'engine_resources', type: 'number', value: b.engine_resources},
                {label: 'Cargo Resources', key: 'cargo_resources', type: 'number', value: b.cargo_resources},
            ]
            for (const f of editFields) {
                const tr = NewE('tr')
                const tdLabel = NewE('td')
                tdLabel.appendChild(NewT(f.label))
                tr.appendChild(tdLabel)
                const tdInput = NewE('td')
                const input = NewE('input')
                input.type = f.type
                input.name = f.key
                input.value = String(f.value)
                tdInput.appendChild(input)
                tr.appendChild(tdInput)
                editTable.appendChild(tr)
            }
            const saveBtn = NewE('button')
            saveBtn.type = 'submit'
            saveBtn.appendChild(NewT('Save'))

            const cancelLink = NewE('a')
            cancelLink.href = `/fleet-build/${buildId}/main.html`
            cancelLink.className = 'btn'
            cancelLink.style.marginLeft = '12px'
            cancelLink.appendChild(NewT('Cancel'))

            const form = NewE('form')
            form.appendChild(editTable)
            form.appendChild(saveBtn)
            form.appendChild(cancelLink)
            form.addEventListener('submit', async (e) => {
                e.preventDefault()
                await this.submitFleetBuildEdit(form, buildId, b)
            })
            container.appendChild(form)
        } catch (e) {
            this.dispatcher.dispatch('displayError', [e.message, true])
        }

        return container
    }

    /**
     * @param {HTMLFormElement} form
     * @param {string} buildId
     * @param {FleetBuild} b
     */
    async submitFleetBuildEdit(form, buildId, b) {
        try {
            const data = Object.fromEntries(new FormData(form))
            data.attack_resources = parseFloat(data.attack_resources) || 0
            data.defense_resources = parseFloat(data.defense_resources) || 0
            data.engine_resources = parseFloat(data.engine_resources) || 0
            data.cargo_resources = parseFloat(data.cargo_resources) || 0
            await this.apiClient.updateFleetBuild(buildId, {...b, ...data})
            this.dispatcher.dispatch('redirect', `/fleet-build/${buildId}/main.html`)
        } catch (e) {
            this.dispatcher.dispatch('displayError', [e.message, true])
        }
    }
}
