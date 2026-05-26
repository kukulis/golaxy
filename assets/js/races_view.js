import {ClearE, NewE, NewT} from '/assets/js/helper.js'
import {ApiClient} from "./api.js";
import {Dispatcher} from "./dispatcher.js";

export default class RacesView {

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
     * @return {HTMLElement}
     */
    async generateView() {
        let viewDiv = NewE('div')

        let table = NewE('table')
        table.style.display = 'none'

        let thead = NewE('thead')
        let tr = NewE('tr')
        for (const header of ['ID', 'Name', 'Role']) {
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
            const races = await this.apiClient.getRaces();
            for (const r of races) {
                const tr = NewE('tr');

                const tdId = NewE('td');
                const a = NewE('a');
                a.href = `/race/${r.id}/main.html`;
                a.appendChild(NewT(r.id));
                tdId.appendChild(a);
                tr.appendChild(tdId);

                for (const val of [r.name, r.role]) {
                    const td = NewE('td');
                    td.appendChild(NewT(val));
                    tr.appendChild(td);
                }

                this.tBody.appendChild(tr);
            }
        } catch (e) {
            this.dispatcher.dispatch("displayError", [e.message, true])
        }
    }
}
