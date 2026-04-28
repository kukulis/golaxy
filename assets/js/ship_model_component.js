import {ApiClient} from "/assets/js/api.js";
import {ShipModel} from './entities/ship_model.js';

export class ShipModelComponent {

    /**
     * @type {ApiClient}
     */
    apiClient = null;

    /**
     * @param {ApiClient} apiClient
     */
    constructor(apiClient) {
        this.apiClient = apiClient;
    }

    /**
     * @param {ShipModel[]} models
     * @returns {HTMLTableElement}
     */
    renderList(models) {
        const thead = document.createElement('thead');
        const headerRow = document.createElement('tr');
        ['ID', 'Name', 'Guns', 'Gun Mass', 'Defense Mass', 'Engine Mass', 'Cargo Mass'].forEach(label => {
            const th = document.createElement('th');
            th.appendChild(document.createTextNode(label));
            headerRow.appendChild(th);
        });
        thead.appendChild(headerRow);

        const tbody = document.createElement('tbody');

        for (const m of models) {
            const tr = document.createElement('tr');

            const tdId = document.createElement('td');
            const linkId = document.createElement('a');
            linkId.appendChild(document.createTextNode(m['id']));
            linkId.setAttribute('href', `/ship-model/${m['id']}/details.html`);
            tdId.appendChild(linkId);
            tr.appendChild(tdId);

            const tdName = document.createElement('td');
            const linkName = document.createElement('a');
            linkName.appendChild(document.createTextNode(m['name']));
            linkName.setAttribute('href', `/ship-model/${m['id']}/details.html`);
            tdName.appendChild(linkName);
            tr.appendChild(tdName);

            ['guns', 'one_gun_mass', 'defense_mass', 'engine_mass', 'cargo_mass'].forEach(key => {
                const td = document.createElement('td');
                td.appendChild(document.createTextNode(m[key]));
                tr.appendChild(td);
            });

            tbody.appendChild(tr);
        }

        const table = document.createElement('table');
        table.appendChild(thead);
        table.appendChild(tbody);
        return table;
    }

    /**
     * @param {string} shipModelId
     * @returns {Promise<HTMLTableElement>}
     */
    async renderDetails(shipModelId) {
        const m = await this.apiClient.getShipModel(shipModelId);
        const tbody = document.createElement('tbody');
        [
            ['ID', m.id],
            ['Name', m.name],
            ['Guns', m.guns],
            ['Gun Mass', m.one_gun_mass],
            ['Defense Mass', m.defense_mass],
            ['Engine Mass', m.engine_mass],
            ['Cargo Mass', m.cargo_mass],
            ['Owner ID', m.owner_id],
        ].forEach(([label, value]) => {
            const tr = document.createElement('tr');
            const th = document.createElement('th');
            th.appendChild(document.createTextNode(label));
            const td = document.createElement('td');
            td.appendChild(document.createTextNode(value));
            tr.appendChild(th);
            tr.appendChild(td);
            tbody.appendChild(tr);
        });

        const table = document.createElement('table');
        table.appendChild(tbody);
        return table;
    }
}
