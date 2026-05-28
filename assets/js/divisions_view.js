import {ClearE, NewE, NewT} from '/assets/js/helper.js'
import {ApiClient} from "./api.js";
import {Dispatcher} from "./dispatcher.js";

export default class DivisionsView {

    /**
     * @type {ApiClient}
     */
    apiClient = null;

    /**
     * @type {Dispatcher}
     */
    dispatcher = null;

    /**
     * @type {HTMLElement}
     */
    tBody = null;

    /**
     * @param {ApiClient} apiClient
     * @param {Dispatcher} dispatcher
     */
    constructor(apiClient, dispatcher) {
        this.apiClient = apiClient
        this.dispatcher = dispatcher
    }

    /**
     * @return {HTMLElement}
     */
    async generateView() {
        let viewDiv = NewE('div')

        let isAdmin = false
        try {
            const race = await this.apiClient.getCurrentRace()
            isAdmin = race.role === 'admin'
        } catch (e) {
            this.dispatcher.dispatch('displayError', [e.message, true])
        }

        let table = NewE('table')
        table.id = 'divisions-table'

        let thead = NewE('thead')
        let tr = NewE('tr')
        const headers = ['ID', 'Resources', 'Attack', 'Defense', 'Engines', 'Cargo', '']
        for (const header of headers) {
            let th = NewE('th')
            th.appendChild(NewT(header))
            tr.appendChild(th)
        }
        thead.appendChild(tr)
        table.appendChild(thead)

        this.tBody = NewE('tbody')
        this.tBody.id = 'divisions-body'
        table.appendChild(this.tBody)

        viewDiv.appendChild(table)

        await this.reloadTableBody(isAdmin)

        if (isAdmin) {
            const addButton = NewE('button')
            addButton.appendChild(NewT('Add division'))
            addButton.addEventListener('click', this.addDivision)
            viewDiv.appendChild(addButton)
        }

        return viewDiv
    }

    addDivision = async () => {
        try {
            await this.apiClient.createDivision({
                id: crypto.randomUUID(),
                resources_amount: 0,
                tech_attack: 1,
                tech_defense: 1,
                tech_engines: 1,
                tech_cargo: 1,
            })
            await this.reloadTableBody(true)
        } catch (e) {
            this.dispatcher.dispatch('displayError', [e.message, true])
        }
    }

    deleteDivision = async (division) => {
        if (!confirm(`Delete division ${division.id}?`)) return
        try {
            await this.apiClient.deleteDivision(division.id)
            await this.reloadTableBody(true)
        } catch (e) {
            this.dispatcher.dispatch('displayError', [e.message, true])
        }
    }

    async reloadTableBody(isAdmin) {
        ClearE(this.tBody)

        try {
            const divisions = await this.apiClient.getDivisions()
            for (const d of divisions) {
                const tr = NewE('tr')

                const tdDivisionId = NewE('td')
                const aDivisionView = NewE('a')
                aDivisionView.href = `/division/${d.id}/main.html`
                aDivisionView.appendChild(NewT(d.id))
                tdDivisionId.appendChild(aDivisionView)
                tr.appendChild(tdDivisionId)

                for (const val of [d.resources_amount, d.tech_attack, d.tech_defense, d.tech_engines, d.tech_cargo]) {
                    const td = NewE('td')
                    td.appendChild(NewT(val))
                    tr.appendChild(td)
                }

                const tdActions = NewE('td')
                if (isAdmin) {
                    const buttonDelete = NewE('button')
                    buttonDelete.appendChild(NewT('Delete'))
                    buttonDelete.addEventListener('click', () => this.deleteDivision(d))
                    tdActions.appendChild(buttonDelete)
                }
                tr.appendChild(tdActions)

                this.tBody.appendChild(tr)
            }
        } catch (e) {
            this.dispatcher.dispatch('displayError', [e.message, true])
        }
    }
}
