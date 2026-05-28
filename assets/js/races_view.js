import {ClearE, NewE, NewT} from '/assets/js/helper.js'
import {ApiClient} from "./api.js";
import {Dispatcher} from "./dispatcher.js";

export default class RacesView {

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

        let thead = NewE('thead')
        let tr = NewE('tr')
        for (const header of ['ID', 'Name', 'Role', '']) {
            let th = NewE('th')
            th.appendChild(NewT(header))
            tr.appendChild(th)
        }
        thead.appendChild(tr)
        table.appendChild(thead)

        this.tBody = NewE('tbody')
        table.appendChild(this.tBody)

        viewDiv.appendChild(table)

        await this.reloadTableBody(isAdmin)

        if (isAdmin) {
            const addButton = NewE('button')
            addButton.appendChild(NewT('Add race'))
            addButton.addEventListener('click', this.addRace)
            viewDiv.appendChild(addButton)
        }

        return viewDiv
    }

    addRace = async () => {
        try {
            await this.apiClient.createRace({
                ID: crypto.randomUUID(),
                Name: 'New Race',
                Role: 'player',
                Token: crypto.randomUUID(),
            })
            await this.reloadTableBody(true)
        } catch (e) {
            this.dispatcher.dispatch('displayError', [e.message, true])
        }
    }

    deleteRace = async (race) => {
        if (!confirm(`Delete race ${race.id} ${race.name}?`)) return
        try {
            await this.apiClient.deleteRace(race.id)
            await this.reloadTableBody(true)
        } catch (e) {
            this.dispatcher.dispatch('displayError', [e.message, true])
        }
    }

    async reloadTableBody(isAdmin) {
        ClearE(this.tBody)

        try {
            const races = await this.apiClient.getRaces()
            for (const r of races) {
                const tr = NewE('tr')

                const tdId = NewE('td')
                const a = NewE('a')
                a.href = `/race/${r.id}/main.html`
                a.appendChild(NewT(r.id))
                tdId.appendChild(a)
                tr.appendChild(tdId)

                for (const val of [r.name, r.role]) {
                    const td = NewE('td')
                    td.appendChild(NewT(val))
                    tr.appendChild(td)
                }

                const tdActions = NewE('td')
                if (isAdmin) {
                    const buttonDelete = NewE('button')
                    buttonDelete.appendChild(NewT('Delete'))
                    buttonDelete.addEventListener('click', () => this.deleteRace(r))
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
