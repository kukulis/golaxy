import {ClearE, GetE, NewE, NewT} from '/assets/js/helper.js'
import {ApiClient} from "./api.js";

export default class DivisionsView {

    /**
     *
     * @type {ApiClient}
     */
    apiClient = null;

    /**
     *
     * @type {HTMLElement}
     */
    tBody = null;

    /**
     *
     * @param {ApiClient} apiClient
     */
    constructor(apiClient) {
        this.apiClient = apiClient
    }

    /**
     * @return {HTMLElement}
     */
    async generateView() {
        let viewDiv = NewE('div')

        let table = NewE('table')
        table.id = 'divisions-table'
        table.style.display = 'none'

        let thead = NewE('thead')
        let tr = NewE('tr')
        for (const header of ['ID', 'Resources', 'Attack', 'Defense', 'Engines', 'Cargo']) {
            let th = NewE('th')
            th.appendChild(NewT(header))
            tr.appendChild(th)
        }
        thead.appendChild(tr)
        table.appendChild(thead)

        this.tBody = NewE('tbody')
        this.tBody.id = 'divisions-body'
        table.appendChild(this.tBody)
        table.style.display = ''

        viewDiv.appendChild(table)

        await this.reloadTableBody()

        // testing
        const refreshButton = NewE('button')
        refreshButton.appendChild(NewT('Refresh'))
        refreshButton.addEventListener('click', async (event) => {
            console.log('TODO implement click')
            await this.reloadTableBody()
        })

        viewDiv.appendChild(refreshButton)

        return viewDiv
    }

    async reloadTableBody() {

        ClearE(this.tBody)

        try {
            const divisions = await this.apiClient.getDivisions();
            for (const d of divisions) {
                const tr = NewE('tr');

                const tdDivisionId = NewE('td');
                const aDivisionView = NewE('a');
                aDivisionView.href = `/division/${d.id}/main.html`;

                // color defined in a css file
                // aDivisionView.style.color = '#4a9eff';
                aDivisionView.appendChild(NewT(d.id));
                tdDivisionId.appendChild(aDivisionView);
                tr.appendChild(tdDivisionId);

                for (const val of [d.resources_amount, d.tech_attack, d.tech_defense, d.tech_engines, d.tech_cargo]) {
                    const td = NewE('td');
                    td.appendChild(NewT(val));
                    tr.appendChild(td);
                }

                this.tBody.appendChild(tr);
            }
        } catch (e) {
            // TODO use own event dispatcher
            const err = GetE('error-msg');
            err.textContent = e.message;
            err.style.display = '';
        }
    }

}