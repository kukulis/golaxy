import {ClearE, NewE, NewT} from '/assets/js/helper.js'
import {ApiClient} from "./api.js";
import {Dispatcher} from "./dispatcher.js";

export default class FleetBuildsView {

    /**
     *
     * @type {ApiClient}
     */
    apiClient = null;

    /**
     *
     * @type {Dispatcher}
     */
    dispatcher = null;

    /**
     *
     * @type {string}
     */
    divisionId = null;

    /**
     *
     * @type {HTMLElement}
     */
    tBody = null;

    /**
     *
     * @param {ApiClient} apiClient
     * @param {Dispatcher} dispatcher
     */
    constructor(apiClient, dispatcher) {
        this.apiClient = apiClient
        this.dispatcher = dispatcher
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
        let viewDiv = NewE('div')

        let table = NewE('table')

        let thead = NewE('thead')
        let tr = NewE('tr')
        for (const header of ['ID', 'Race', '']) {
            let th = NewE('th')
            th.appendChild(NewT(header))
            tr.appendChild(th)
        }
        thead.appendChild(tr)
        table.appendChild(thead)

        this.tBody = NewE('tbody')
        table.appendChild(this.tBody)

        viewDiv.appendChild(table)

        await this.reloadTableBody()

        const addButton = NewE('button')
        addButton.appendChild(NewT('Add fleet build'))
        addButton.addEventListener('click', this.addFleetBuild)
        viewDiv.appendChild(addButton)

        return viewDiv
    }

    addFleetBuild = async () => {
        try {
            const race = await this.apiClient.getCurrentRace()
            await this.apiClient.createFleetBuild({
                id: crypto.randomUUID(),
                race_id: race.id,
                division_id: this.divisionId,
            })
            await this.reloadTableBody(true)
        } catch (e) {
            this.dispatcher.dispatch('displayError', [e.message, true])
        }
    }

    deleteFleetBuild = async (build) => {
        if (!confirm(`Delete fleet build ${build.id}?`)) return
        try {
            await this.apiClient.deleteFleetBuild(build.id)
            await this.reloadTableBody(true)
        } catch (e) {
            this.dispatcher.dispatch('displayError', [e.message, true])
        }
    }

    async reloadTableBody() {
        ClearE(this.tBody)

        try {
            const builds = await this.apiClient.getFleetBuilds(this.divisionId);
            for (const b of builds) {
                const tr = NewE('tr');

                const tdEdit = NewE('td');
                const editLink = NewE('a');
                editLink.href = `/fleet-build/${b['id']}/edit.html`;
                editLink.appendChild(NewT('✏ Edit'));
                tdEdit.appendChild(editLink);
                tr.appendChild(tdEdit);

                const tdId = NewE('td');
                const a = NewE('a');
                a.href = `/fleet-build/${b.id}/main.html`;
                a.appendChild(NewT(b.id));
                tdId.appendChild(a);
                tr.appendChild(tdId);

                const tdRace = NewE('td');
                tdRace.appendChild(NewT(b.race_id ?? ''));
                tr.appendChild(tdRace);

                const tdActions = NewE('td')
                const buttonDelete = NewE('button')
                buttonDelete.appendChild(NewT('Delete'))
                buttonDelete.addEventListener('click', () => this.deleteFleetBuild(b))
                tdActions.appendChild(buttonDelete)

                tr.appendChild(tdActions)

                this.tBody.appendChild(tr);
            }
        } catch (e) {
            this.dispatcher.dispatch("displayError", [e.message, true])
        }
    }
}
