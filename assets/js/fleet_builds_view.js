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
     * @type {boolean}
     */
    showAll = false;

    /**
     * @type {string}
     */
    filterRaceId = '';

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

        const allCheckbox = NewE('input')
        allCheckbox.type = 'checkbox'
        allCheckbox.id = 'filter-all'
        const allLabel = NewE('label')
        allLabel.htmlFor = 'filter-all'
        allLabel.appendChild(NewT('All'))
        filterDiv.appendChild(allCheckbox)
        filterDiv.appendChild(allLabel)

        const raceSelect = NewE('select')
        raceSelect.disabled = true
        const defaultOption = NewE('option')
        defaultOption.value = ''
        defaultOption.appendChild(NewT('All races'))
        raceSelect.appendChild(defaultOption)
        try {
            const races = await this.apiClient.getRaces()
            for (const r of races) {
                const option = NewE('option')
                option.value = r.id
                option.appendChild(NewT(r.name || r.id))
                raceSelect.appendChild(option)
            }
        } catch (e) {
            this.dispatcher.dispatch('displayError', [e.message, true])
        }
        filterDiv.appendChild(raceSelect)

        allCheckbox.addEventListener('change', async () => {
            this.showAll = allCheckbox.checked
            raceSelect.disabled = !this.showAll
            if (!this.showAll) {
                this.filterRaceId = ''
                raceSelect.value = ''
            }
            await this.reloadTableBody()
        })
        raceSelect.addEventListener('change', async () => {
            this.filterRaceId = raceSelect.value
            await this.reloadTableBody()
        })

        viewDiv.appendChild(filterDiv)

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
            const builds = await this.apiClient.getFleetBuilds(this.divisionId, this.showAll, this.filterRaceId);
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
