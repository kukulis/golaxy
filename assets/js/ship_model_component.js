import {NewE, NewT} from '/assets/js/helper.js';
import {ApiClient} from "/assets/js/api.js";
import {Dispatcher} from './dispatcher.js';
import {ShipModel} from './entities/ship_model.js';

export class ShipModelComponent {

    /**
     * @type {ApiClient}
     */
    apiClient = null;

    /**
     * @type {Dispatcher}
     */
    dispatcher = null;

    /**
     * @param {ApiClient} apiClient
     * @param {Dispatcher} dispatcher
     */
    constructor(apiClient, dispatcher) {
        this.apiClient = apiClient;
        this.dispatcher = dispatcher;
    }

    /**
     * @returns {Promise<HTMLElement>}
     */
    async renderList() {
        const thead = NewE('thead');
        const headerRow = NewE('tr');
        ['', 'ID', 'Name', 'Guns', 'Gun Mass', 'Defense Mass', 'Engine Mass', 'Cargo Mass'].forEach(label => {
            const th = NewE('th');
            th.appendChild(NewT(label));
            headerRow.appendChild(th);
        });
        thead.appendChild(headerRow);

        const tbody = NewE('tbody');

        try {
            const models = await this.apiClient.getShipModels();
            for (const m of models) {
                const tr = NewE('tr');

                const tdEdit = NewE('td');
                const linkEdit = NewE('a');
                linkEdit.appendChild(NewT('✏'));
                linkEdit.setAttribute('href', `/ship-model/${m['id']}/edit.html`);
                tdEdit.appendChild(linkEdit);
                tr.appendChild(tdEdit);

                const tdId = NewE('td');
                const linkId = NewE('a');
                linkId.appendChild(NewT(m['id']));
                linkId.setAttribute('href', `/ship-model/${m['id']}/details.html`);
                tdId.appendChild(linkId);
                tr.appendChild(tdId);

                const tdName = NewE('td');
                const linkName = NewE('a');
                linkName.appendChild(NewT(m['name']));
                linkName.setAttribute('href', `/ship-model/${m['id']}/details.html`);
                tdName.appendChild(linkName);
                tr.appendChild(tdName);

                ['guns', 'one_gun_mass', 'defense_mass', 'engine_mass', 'cargo_mass'].forEach(key => {
                    const td = NewE('td');
                    td.appendChild(NewT(m[key]));
                    tr.appendChild(td);
                });

                tbody.appendChild(tr);
            }
        } catch (e) {
            this.dispatcher.dispatch('displayError', [e.message, true]);
        }

        const table = NewE('table');
        table.appendChild(thead);
        table.appendChild(tbody);
        return table;
    }

    /**
     * @param {string} shipModelId
     * @returns {Promise<HTMLElement>}
     */
    async renderDetails(shipModelId) {
        try {
            const m = await this.apiClient.getShipModel(shipModelId);
            const tbody = NewE('tbody');
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
                const tr = NewE('tr');
                const th = NewE('th');
                th.appendChild(NewT(label));
                const td = NewE('td');
                td.appendChild(NewT(value));
                tr.appendChild(th);
                tr.appendChild(td);
                tbody.appendChild(tr);
            });

            const table = NewE('table');
            table.appendChild(tbody);

            const editLink = NewE('a');
            editLink.href = `/ship-model/${shipModelId}/edit.html`;
            editLink.appendChild(NewT('✏ Edit'));

            const wrapper = NewE('div');
            wrapper.appendChild(table);
            wrapper.appendChild(editLink);
            return wrapper;
        } catch (e) {
            this.dispatcher.dispatch('displayError', [e.message, true]);
            return NewE('div');
        }
    }

    /**
     * @param {string} shipModelId
     * @returns {Promise<HTMLFormElement>}
     */
    async renderEdit(shipModelId) {
        try {
            const m = await this.apiClient.getShipModel(shipModelId);

            const tbody = NewE('tbody');
            [
                ['Name', 'name', 'text', m.name],
                ['Guns', 'guns', 'number', m.guns],
                ['Gun Mass', 'one_gun_mass', 'number', m.one_gun_mass],
                ['Defense Mass', 'defense_mass', 'number', m.defense_mass],
                ['Engine Mass', 'engine_mass', 'number', m.engine_mass],
                ['Cargo Mass', 'cargo_mass', 'number', m.cargo_mass],
            ].forEach(([label, name, type, value]) => {
                const tr = NewE('tr');
                const th = NewE('th');
                th.appendChild(NewT(label));
                const td = NewE('td');
                const input = NewE('input');
                input.type = type;
                input.name = name;
                input.value = value;
                td.appendChild(input);
                tr.appendChild(th);
                tr.appendChild(td);
                tbody.appendChild(tr);
            });

            const table = NewE('table');
            table.appendChild(tbody);

            const button = NewE('button');
            button.type = 'submit';
            button.appendChild(NewT('Save'));

            const cancelLink = NewE('a');
            cancelLink.href = `/ship-model/${shipModelId}/details.html`;
            cancelLink.className = 'btn';
            cancelLink.style.marginLeft = '12px';
            cancelLink.appendChild(NewT('Cancel'));

            const form = NewE('form');
            form.appendChild(table);
            form.appendChild(button);
            form.appendChild(cancelLink);
            form.addEventListener('submit', async (e) => {
                e.preventDefault();
                try {
                    const data = Object.fromEntries(new FormData(form));
                    data.guns = parseInt(data.guns);
                    ['one_gun_mass', 'defense_mass', 'engine_mass', 'cargo_mass'].forEach(k => data[k] = parseFloat(data[k]));
                    await this.apiClient.updateShipModel(shipModelId, data);
                    this.dispatcher.dispatch('afterShipModelEdit', shipModelId);
                } catch (e) {
                    this.dispatcher.dispatch('displayError', [e.message, true]);
                }
            });

            return form;
        } catch (e) {
            this.dispatcher.dispatch('displayError', [e.message, true]);
            return NewE('div');
        }
    }
}
