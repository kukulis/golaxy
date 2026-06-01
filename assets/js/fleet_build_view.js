import {ClearE, NewE, NewT} from '/assets/js/helper.js'
import {ApiClient} from './api.js'
import {Dispatcher} from './dispatcher.js'
import {createFleetBuildStatisticsTable} from './fleet_build_table.js'
import {createDivisionTable} from './division_table.js'
import {ShipTech} from './entities/ship_tech.js'

export default class FleetBuildView {

    /**
     * @type {ApiClient}
     */
    apiClient = null

    /**
     * @type {Dispatcher}
     */
    dispatcher = null

    /**
     * @type {string}
     */
    buildId = null

    /**
     * @type {string}
     */
    divisionId = null

    // /**
    //  *
    //  * @type {HTMLElement}
    //  */
    // rightCold = null;

    /**
     *
     * @type {HTMLElement}
     */
    assignmentsWrapper = null;

    /**
     * @param {ApiClient} apiClient
     * @param {Dispatcher} dispatcher
     */
    constructor(apiClient, dispatcher) {
        this.apiClient = apiClient
        this.dispatcher = dispatcher
    }

    /**
     * @param {string} buildId
     */
    setBuildId(buildId) {
        this.buildId = buildId
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
        const container = NewE('div')

        try {
            const [b, division, statistics] = await Promise.all([
                this.apiClient.getFleetBuild(this.buildId),
                this.apiClient.getDivision(this.divisionId),
                this.apiClient.getFleetBuildStatistics(this.buildId),
            ])

            const leftCol = NewE('div')

            const editLink = NewE('a')
            editLink.href = `/fleet-build/${this.buildId}/edit.html`
            editLink.className = 'btn'
            editLink.appendChild(NewT('✏ Edit'))
            leftCol.appendChild(editLink)

            const fleetBuildTitle = NewE('h2')
            fleetBuildTitle.appendChild(NewT(b.name || 'Fleet Build'))
            leftCol.appendChild(fleetBuildTitle)
            leftCol.appendChild(createFleetBuildStatisticsTable(b, statistics))

            const challengeBtn = NewE('button')
            challengeBtn.appendChild(NewT('Challenge'))
            challengeBtn.style.display = b.ready_for_battle ? '' : 'none'
            challengeBtn.addEventListener('click', async () => {
                try {
                    await this.apiClient.createChallenge({id: crypto.randomUUID(), fleet_build_a_id: b.id})
                } catch (e) {
                    this.dispatcher.dispatch('displayError', [e.message, true])
                }
            })
            leftCol.appendChild(challengeBtn)

            const divisionTitle = NewE('h2')
            divisionTitle.appendChild(NewT('Division'))
            leftCol.appendChild(divisionTitle)
            leftCol.appendChild(createDivisionTable(division))

            const rightCol = NewE('div')

            const shipModelsTitle = NewE('h2')
            shipModelsTitle.appendChild(NewT('Assigned Ship Models'))
            rightCol.appendChild(shipModelsTitle)

            const addBtn = NewE('button')
            addBtn.appendChild(NewT('+ Add'))
            addBtn.addEventListener('click', async () => { await this._addAssignment() })
            rightCol.appendChild(addBtn)

            this.assignmentsWrapper = NewE('div')
            rightCol.appendChild(this.assignmentsWrapper)

            await this.reloadAssignmentsTable()


            container.className = 'two-col-layout'
            container.appendChild(leftCol)
            container.appendChild(rightCol)
        } catch (e) {
            this.dispatcher.dispatch('displayError', [e.message, true])
        }

        return container
    }

    async reloadAssignmentsTable() {
        const [assignments, shipsTechs] = await Promise.all([
            this.apiClient.getFleetBuildShipModelsAssignments(this.buildId),
            this.apiClient.calculateAssignmentsShipTechs(this.buildId),
        ])

        ClearE(this.assignmentsWrapper)
        this.assignmentsWrapper.appendChild(this.createFleetBuildShipModelTable(assignments, shipsTechs))
    }

    /**
     * @param {string} buildId
     * @return {HTMLElement}
     */
    async generateEditView(buildId) {
        const container = NewE('div')

        try {
            const b = await this.apiClient.getFleetBuild(buildId)

            const readOnlyTable = NewE('table')
            for (const [label, value] of [
                ['ID', b.id],
                ['Division', b.division_id],
                ['Race', b.race_id],
                ['Used Resources', b.usedResources],
            ]) {
                const tr = NewE('tr')
                const tdLabel = NewE('td')
                tdLabel.appendChild(NewT(label))
                tr.appendChild(tdLabel)
                const tdValue = NewE('td')
                tdValue.appendChild(NewT(String(value)))
                tr.appendChild(tdValue)
                readOnlyTable.appendChild(tr)
            }
            container.appendChild(readOnlyTable)

            const editTable = NewE('table')
            const editFields = [
                {label: 'Name', key: 'name', type: 'text', value: b.name},
                {label: 'Attack Resources', key: 'attack_resources', type: 'number', value: b.attack_resources},
                {label: 'Defense Resources', key: 'defense_resources', type: 'number', value: b.defense_resources},
                {label: 'Engine Resources', key: 'engine_resources', type: 'number', value: b.engine_resources},
                {label: 'Cargo Resources', key: 'cargo_resources', type: 'number', value: b.cargo_resources},
            ]
            for (const f of editFields) {
                const tr = NewE('tr')
                const tdLabel = NewE('td')
                tdLabel.appendChild(NewT(f.label))
                tr.appendChild(tdLabel)
                const tdInput = NewE('td')
                const input = NewE('input')
                input.type = f.type
                input.name = f.key
                input.value = String(f.value)
                tdInput.appendChild(input)
                tr.appendChild(tdInput)
                editTable.appendChild(tr)
            }

            const readyTr = NewE('tr')
            const readyTdLabel = NewE('td')
            readyTdLabel.appendChild(NewT('Ready for battle'))
            readyTr.appendChild(readyTdLabel)
            const readyTdInput = NewE('td')
            const readyCheckbox = NewE('input')
            readyCheckbox.type = 'checkbox'
            readyCheckbox.name = 'ready_for_battle'
            readyCheckbox.checked = b.ready_for_battle
            readyTdInput.appendChild(readyCheckbox)
            readyTr.appendChild(readyTdInput)
            editTable.appendChild(readyTr)
            const saveBtn = NewE('button')
            saveBtn.type = 'submit'
            saveBtn.appendChild(NewT('Save'))

            const cancelLink = NewE('a')
            cancelLink.href = `/fleet-build/${buildId}/main.html`
            cancelLink.className = 'btn'
            cancelLink.style.marginLeft = '12px'
            cancelLink.appendChild(NewT('Cancel'))

            const form = NewE('form')
            form.appendChild(editTable)
            form.appendChild(saveBtn)
            form.appendChild(cancelLink)
            form.addEventListener('submit', async (e) => {
                e.preventDefault()
                await this.submitFleetBuildEdit(form, buildId, b)
            })
            container.appendChild(form)
        } catch (e) {
            this.dispatcher.dispatch('displayError', [e.message, true])
        }

        return container
    }

    /**
     * @param {HTMLFormElement} form
     * @param {string} buildId
     * @param {FleetBuild} b
     */
    async submitFleetBuildEdit(form, buildId, b) {
        try {
            const data = Object.fromEntries(new FormData(form))
            data.attack_resources = parseFloat(data.attack_resources) || 0
            data.defense_resources = parseFloat(data.defense_resources) || 0
            data.engine_resources = parseFloat(data.engine_resources) || 0
            data.cargo_resources = parseFloat(data.cargo_resources) || 0
            data.ready_for_battle = 'ready_for_battle' in data

            console.log ( "fleet_build_view.js [248] data: ", data )
            await this.apiClient.updateFleetBuild(buildId, {...b, ...data})
            this.dispatcher.dispatch('redirect', `/fleet-build/${buildId}/main.html`)
            // alert("fleet build saved")
        } catch (e) {
            this.dispatcher.dispatch('displayError', [e.message, true])
        }
    }

    async _addAssignment() {
        try {
            const id = crypto.randomUUID()
            await this.apiClient.addShipModelAssignment({
                id:            id,
                fleet_build_id: this.buildId,
                ship_model_id: '',
                amount:        1,
            })
            await this.reloadAssignmentsTable()
        } catch (e) {
            this.dispatcher.dispatch('displayError', [e.message, true])
        }
    }

    async _deleteAssignment(assignmentId) {
        if (!confirm('Delete this assignment?')) {
            return
        }
        try {
            await this.apiClient.unassignShipModel(assignmentId)
            await this.reloadAssignmentsTable()
        } catch (e) {
            this.dispatcher.dispatch('displayError', [e.message, true])
        }
    }

    /**
     * @param {FleetBuildShipModel[]} assignments
     * @param {{assignment_id: string, ship_tech: ShipTech}[]} shipTechs
     * @returns {HTMLTableElement}
     */
    createFleetBuildShipModelTable(assignments, shipTechs) {
        const thead = document.createElement('thead');
        const headerRow = document.createElement('tr');
        ['', 'Name', 'Guns', 'Gun Mass (Attack)', 'Defense Mass (Defense)', 'Engine Mass (Speed)', 'Cargo Mass (Cargo Cap.)', 'Amount', 'Result Mass (Ship Mass)', ''].forEach(label => {
            const th = document.createElement('th');
            th.appendChild(document.createTextNode(label));
            headerRow.appendChild(th);
        });
        thead.appendChild(headerRow);

        const techByAssignmentId = new Map()
        for (const entry of shipTechs) {
            techByAssignmentId.set(entry.assignment_id, entry.ship_tech)
        }


        const tbody = document.createElement('tbody');
        for (const assignment of assignments) {
            const rawTech = techByAssignmentId.get(assignment.id) ?? {}
            const tech = {
                attack:         rawTech.attack         ?? 0,
                defense:        rawTech.defense        ?? 0,
                speed:          rawTech.speed          ?? 0,
                cargo_capacity: rawTech.cargo_capacity ?? 0,
                mass:           rawTech.mass           ?? 0,
            }
            const tr = document.createElement('tr');
            const editTd = document.createElement('td');
            const editLink = document.createElement('a');
            editLink.href = `/ship-model-assignment/${assignment.id}/main.html`;
            editLink.appendChild(NewT('View/Edit'));
            editTd.appendChild(editLink);
            tr.appendChild(editTd);

            const displayTech = (v) => ' (' + v.toFixed(2) + ')';

           let displayedData =
            [
                assignment.shipModel.name,
                assignment.shipModel.guns,
                '' + assignment.shipModel.one_gun_mass + displayTech(tech.attack),
                '' + assignment.shipModel.defense_mass + displayTech(tech.defense),
                '' + assignment.shipModel.engine_mass + displayTech(tech.speed),
                '' + assignment.shipModel.cargo_mass + displayTech(tech.cargo_capacity),
                assignment.amount,
                '' + assignment.result_mass + displayTech(tech.mass),
            ];

            displayedData.forEach(value => {
                const td = document.createElement('td');
                td.appendChild(document.createTextNode(value));
                tr.appendChild(td);
            });

            const deleteTd = document.createElement('td');
            const deleteBtn = NewE('button');
            deleteBtn.appendChild(NewT('Delete'));
            deleteBtn.addEventListener('click', async () => { await this._deleteAssignment(assignment.id) });
            deleteTd.appendChild(deleteBtn);
            tr.appendChild(deleteTd);

            tbody.appendChild(tr);
        }

        const table = document.createElement('table');
        table.appendChild(thead);
        table.appendChild(tbody);
        return table;
    }
}
