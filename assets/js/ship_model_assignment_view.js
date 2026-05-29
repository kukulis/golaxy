import {ClearE, NewE, NewT} from '/assets/js/helper.js'
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
     * @type {string|null}
     */
    shipModelId = null

    /**
     * @type {Object|null}
     */
    shipModel = null

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
     * Re-renders into #assignment-container, preserving view state.
     */
    async refresh() {
        const root = document.getElementById('assignment-container')
        if (root) {
            ClearE(root).appendChild(await this.generateView())
        }
    }

    /**
     * @return {HTMLElement}
     */
    async generateView() {
        const container = NewE('div')

        try {
            const assignment = await this.apiClient.getShipModelAssignment(this.assignmentId)

            const title = NewE('h2')
            title.appendChild(NewT(assignment.shipModel.name))
            container.appendChild(title)

            const amountInput = NewE('input')
            amountInput.type = 'number'
            amountInput.name = 'amount'
            amountInput.value = String(assignment.amount)

            const selectionContainer = NewE('div')

            const table = NewE('table')
            table.style.width = 'max-content'
            table.style.margin = '0 auto'

            const sm = this.shipModel ?? assignment.shipModel
            for (const [label, content] of [
                ['Assignment ID', NewT(assignment.id)],
                ['Fleet Build',   NewT(assignment.fleet_build_id)],
                ['Ship Model ID', this._shipModelIdCell(assignment, selectionContainer)],
                ['Amount',        amountInput],
                ['Result Mass',   NewT(String(assignment.result_mass))],
                ['Guns',          NewT(String(sm.guns))],
                ['Gun Mass',      NewT(String(sm.one_gun_mass))],
                ['Defense Mass',  NewT(String(sm.defense_mass))],
                ['Engine Mass',   NewT(String(sm.engine_mass))],
                ['Cargo Mass',    NewT(String(sm.cargo_mass))],
            ]) {
                const tr = NewE('tr')
                const tdLabel = NewE('td')
                tdLabel.appendChild(NewT(label))
                tr.appendChild(tdLabel)
                const tdValue = NewE('td')
                tdValue.appendChild(content)
                tr.appendChild(tdValue)
                table.appendChild(tr)
            }

            const saveBtn = NewE('button')
            saveBtn.type = 'submit'
            saveBtn.appendChild(NewT('Save'))

            const cancelLink = NewE('a')
            cancelLink.href = `/fleet-build/${this.fleetBuildId}/main.html`
            cancelLink.className = 'btn'
            cancelLink.style.marginLeft = '8px'
            cancelLink.appendChild(NewT('Cancel'))

            const form = NewE('form')
            form.appendChild(table)
            form.appendChild(selectionContainer)
            form.appendChild(saveBtn)
            form.appendChild(cancelLink)
            form.addEventListener('submit', async (e) => {
                e.preventDefault()
                await this._save(
                    this.shipModelId ?? assignment.ship_model_id,
                    parseInt(amountInput.value) || 0,
                    assignment.result_mass,
                )
            })

            container.appendChild(form)
        } catch (e) {
            this.dispatcher.dispatch('displayError', [e.message, true])
        }

        return container
    }

    /**
     * @param {Object} assignment
     * @param {HTMLElement} selectionContainer
     * @return {HTMLElement}
     */
    _shipModelIdCell(assignment, selectionContainer) {
        const span = NewE('span')
        span.appendChild(NewT(this.shipModelId ?? assignment.ship_model_id))

        const btn = NewE('button')
        btn.type = 'button'
        btn.appendChild(NewT('...'))
        btn.style.marginLeft = '8px'
        btn.addEventListener('click', () => this.renderShipModelSelection(selectionContainer))

        const wrapper = NewE('span')
        wrapper.appendChild(span)
        wrapper.appendChild(btn)
        return wrapper
    }

    /**
     * @param {string} shipModelId
     * @param {number} amount
     * @param {number} resultMass
     */
    async _save(shipModelId, amount, resultMass) {
        try {
            await this.apiClient.updateShipModelAssignment(this.assignmentId, {
                ship_model_id: shipModelId,
                amount:        amount,
                result_mass:   resultMass,
            })
            this.dispatcher.dispatch('redirect', `/fleet-build/${this.fleetBuildId}/main.html`)
        } catch (e) {
            this.dispatcher.dispatch('displayError', [e.message, true])
        }
    }

    /**
     * @param {HTMLElement} container
     */
    async renderShipModelSelection(container) {
        try {
            const shipModels = await this.apiClient.getShipModels()

            ClearE(container)

            const table = NewE('table')
            const thead = NewE('thead')
            const headerRow = NewE('tr')
            for (const label of ['ID', 'Name', 'Guns', 'Gun Mass', 'Defense Mass', 'Engine Mass', 'Cargo Mass', '']) {
                const th = NewE('th')
                th.appendChild(NewT(label))
                headerRow.appendChild(th)
            }
            thead.appendChild(headerRow)
            table.appendChild(thead)

            const tbody = NewE('tbody')
            for (const sm of shipModels) {
                const tr = NewE('tr')

                for (const value of [sm.id, sm.name]) {
                    const td = NewE('td')
                    td.appendChild(NewT(String(value)))
                    td.style.cursor = 'pointer'
                    td.addEventListener('click', () => {
                        this.shipModelId = sm.id
                        this.shipModel = sm
                        this.refresh()
                    })
                    tr.appendChild(td)
                }
                for (const value of [sm.guns, sm.one_gun_mass, sm.defense_mass, sm.engine_mass, sm.cargo_mass]) {
                    const td = NewE('td')
                    td.appendChild(NewT(String(value)))
                    tr.appendChild(td)
                }

                const tdEdit = NewE('td')
                const editLink = NewE('a')
                editLink.href = `/ship-model/${sm.id}/edit.html?DivisionId=${this.divisionId}&FleetBuildId=${this.fleetBuildId}`
                editLink.className = 'btn'
                editLink.appendChild(NewT('Edit'))
                tdEdit.appendChild(editLink)
                tr.appendChild(tdEdit)

                tbody.appendChild(tr)
            }
            table.appendChild(tbody)
            container.appendChild(table)
        } catch (e) {
            this.dispatcher.dispatch('displayError', [e.message, true])
        }
    }
}
