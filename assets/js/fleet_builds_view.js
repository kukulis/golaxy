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
     * @type {HTMLElement}
     */
    tBody = null;

    /**
     * @type {string}
     */
    filterRaceId = '';

    /**
     * @type {boolean}
     */
    isAdmin = false;

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

        const filterDiv = NewE('div')

        const loggedInRace = await this.apiClient.getCurrentRace()
        this.isAdmin = loggedInRace.role === 'admin'
        this.filterRaceId = loggedInRace.id

        const raceSelect = NewE('select')

        if (this.isAdmin) {
            const defaultOption = NewE('option')
            defaultOption.value = ''
            defaultOption.appendChild(NewT('All races'))
            raceSelect.appendChild(defaultOption)
        }

        if ( this.isAdmin) {
            try {
                const races = await this.apiClient.getRaces()
                for (const r of races) {
                    const option = NewE('option')
                    option.value = r.id
                    if (r.id === this.filterRaceId) {
                        option.selected = true
                    }
                    option.appendChild(NewT(r.name || r.id))
                    raceSelect.appendChild(option)
                }
            } catch (e) {
                this.dispatcher.dispatch('displayError', [e.message, true])
            }
        }
        else {
            const option = NewE('option')
            option.value = loggedInRace.id
            option.appendChild(NewT(loggedInRace.name || loggedInRace.id))
            option.selected = true
            raceSelect.appendChild(option)
        }
        filterDiv.appendChild(raceSelect)

        raceSelect.addEventListener('change', async () => {
            this.filterRaceId = raceSelect.value
            await this.reloadTableBody()
        })

        viewDiv.appendChild(filterDiv)

        let table = NewE('table')

        let thead = NewE('thead')
        let tr = NewE('tr')
        for (const header of ['ID', 'Name', 'Race', '']) {
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
            const builds = await this.apiClient.getFleetBuilds(this.divisionId, this.filterRaceId);
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

                const tdName = NewE('td');
                tdName.appendChild(NewT(b.name ?? ''));
                tr.appendChild(tdName);

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
