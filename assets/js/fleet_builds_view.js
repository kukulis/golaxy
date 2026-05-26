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
        table.style.display = 'none'

        let thead = NewE('thead')
        let tr = NewE('tr')
        for (const header of ['ID', 'Race']) {
            let th = NewE('th')
            th.appendChild(NewT(header))
            tr.appendChild(th)
        }
        thead.appendChild(tr)
        table.appendChild(thead)

        this.tBody = NewE('tbody')
        table.appendChild(this.tBody)
        table.style.display = ''

        viewDiv.appendChild(table)

        await this.reloadTableBody()

        return viewDiv
    }

    async reloadTableBody() {
        ClearE(this.tBody)

        try {
            const builds = await this.apiClient.getFleetBuilds(this.divisionId);
            for (const b of builds) {
                const tr = NewE('tr');

                const tdId = NewE('td');
                const a = NewE('a');
                a.href = `/fleet-build/${b.id}/main.html`;
                a.appendChild(NewT(b.id));
                tdId.appendChild(a);
                tr.appendChild(tdId);

                const tdRace = NewE('td');
                tdRace.appendChild(NewT(b.race_id ?? ''));
                tr.appendChild(tdRace);

                this.tBody.appendChild(tr);
            }
        } catch (e) {
            this.dispatcher.dispatch("displayError", [e.message, true])
        }
    }
}
