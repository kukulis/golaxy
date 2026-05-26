import {ClearE, NewE, NewT} from '/assets/js/helper.js'
import {ApiClient} from "./api.js";
import {Dispatcher} from "./dispatcher.js";

export default class DummyLoginView {

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
        for (const header of ['Name', 'Role']) {
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
                tr.style.cursor = 'pointer';
                tr.addEventListener('click', () => {
                    this.dispatcher.dispatch('storeLoginData', [r.name, r.token]);
                    this.dispatcher.dispatch('redirect', '/');
                });

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
