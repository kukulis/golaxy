import {NewE, NewT} from '/assets/js/helper.js';
import {ApiClient} from "/assets/js/api.js";
import {Dispatcher} from './dispatcher.js';
import {ShipModel} from './entities/ship_model.js';
import {createFleetBuildStatisticsTable} from './fleet_build_table.js';
import {createDivisionTable} from './division_table.js';

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
        ['', 'ID', 'Name', 'Guns', 'Gun Mass', 'Defense Mass', 'Engine Mass', 'Cargo Mass', ''].forEach(label => {
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

                const tdDelete = NewE ('td');
                const buttonDelete = NewE('button');
                tdDelete.appendChild(buttonDelete);
                buttonDelete.appendChild(NewT('Delete'));
                const thisShipModel = this;
                tdDelete.addEventListener('click', () =>  thisShipModel.deleteShipModel(m['id']));

                tr.appendChild(tdDelete);

                tbody.appendChild(tr);
            }
        } catch (e) {
            this.dispatcher.dispatch('displayError', [e.message, true]);
        }

        const table = NewE('table');
        table.appendChild(thead);
        table.appendChild(tbody);

        const addButton = NewE('button');
        addButton.appendChild(NewT('Add ship model'));
        addButton.addEventListener('click', this.addShipModel);

        const wrapper = NewE('div');
        wrapper.appendChild(table);
        wrapper.appendChild(addButton);

        return wrapper;
    }

    async addShipModel(e) {
        console.log("ShipModelComponent.addShipModel called");

        const shipModel = new ShipModel();

        // TODO generate ID, generate Name and other parameters
        // send object to the backend through apiClient

    }

    async deleteShipModel(id) {
        console.log("ShipModelComponent.deleteShipModel id="+id)

        // TODO call apiClient
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

            const fleetBuildId = new URLSearchParams(window.location.search).get('FleetBuildId')
            const editLink = NewE('a');
            editLink.href = `/ship-model/${shipModelId}/edit.html` + (fleetBuildId ? `?FleetBuildId=${fleetBuildId}` : '');
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
            const fleetBuildId = new URLSearchParams(window.location.search).get('FleetBuildId')

            const requests = [this.apiClient.getShipModel(shipModelId)]
            if (fleetBuildId) {
                requests.push(this.apiClient.getFleetBuild(fleetBuildId))
            }
            const [m, fleetBuild] = await Promise.all(requests)

            const wrapper = NewE('div')
            let formCol = wrapper

            if (fleetBuild) {
                const [division, technologies] = await Promise.all([
                    this.apiClient.getDivision(fleetBuild.division_id),
                    this.apiClient.getFleetBuildTechnologies(fleetBuildId),
                ])

                const leftCol = NewE('div')

                const fleetBuildTitle = NewE('h2')
                fleetBuildTitle.appendChild(NewT('Fleet Build'))
                leftCol.appendChild(fleetBuildTitle)
                leftCol.appendChild(createFleetBuildStatisticsTable(fleetBuild, technologies))

                const divisionTitle = NewE('h2')
                divisionTitle.appendChild(NewT('Division'))
                leftCol.appendChild(divisionTitle)
                leftCol.appendChild(createDivisionTable(division))

                formCol = NewE('div')
                wrapper.className = 'two-col-layout'
                wrapper.appendChild(leftCol)
                wrapper.appendChild(formCol)
            }

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
            cancelLink.href = `/ship-model/${shipModelId}/details.html` + (fleetBuildId ? `?FleetBuildId=${fleetBuildId}` : '');
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
                    this.dispatcher.dispatch('afterShipModelEdit', {shipModelId, fleetBuildId});
                } catch (e) {
                    this.dispatcher.dispatch('displayError', [e.message, true]);
                }
            });

            formCol.appendChild(form)
            return wrapper;
        } catch (e) {
            this.dispatcher.dispatch('displayError', [e.message, true]);
            return NewE('div');
        }
    }
}
