import {NewE, NewT} from '/assets/js/helper.js'
import {ApiClient} from './api.js'
import {Dispatcher} from './dispatcher.js'
import {formatDateTime} from './date_format.js'

export default class ChallengeView {

    /**
     * @type {ApiClient}
     */
    apiClient = null

    /**
     * @type {Dispatcher}
     */
    dispatcher = null

    /**
     * @param {ApiClient} apiClient
     * @param {Dispatcher} dispatcher
     */
    constructor(apiClient, dispatcher) {
        this.apiClient = apiClient
        this.dispatcher = dispatcher
    }

    /**
     * @param {string} challengeId
     * @return {HTMLElement}
     */
    async generateEditView(challengeId) {
        const container = NewE('div')

        try {
            const [ch, race] = await Promise.all([
                this.apiClient.getChallenge(challengeId),
                this.apiClient.getCurrentRace(),
            ])

            const isChallenger = race.id === ch.challenger_race_id
            const isChallengee = race.id === ch.challengee_race_id
            const isAdmin = race.role === 'admin'

            const canEditA = isChallenger || isAdmin
            const canEditB = isChallengee || isAdmin

            const [buildsA, buildsB, fleetBuildA, fleetBuildB] = await Promise.all([
                canEditA ? this.apiClient.getFleetBuilds(ch.division_id, false, isAdmin ? ch.challenger_race_id : '') : Promise.resolve([]),
                canEditB ? this.apiClient.getFleetBuilds(ch.division_id, false, isAdmin ? ch.challengee_race_id : '') : Promise.resolve([]),
                !canEditA && ch.fleet_build_a_id ? this.apiClient.getFleetBuild(ch.fleet_build_a_id) : Promise.resolve(null),
                !canEditB && ch.fleet_build_b_id ? this.apiClient.getFleetBuild(ch.fleet_build_b_id) : Promise.resolve(null),
            ])

            // Read-only info
            const infoTable = NewE('table')
            for (const [label, value] of [
                ['ID',         ch.id],
                ['Division',   ch.division_id],
                ['Challenger', ch.challenger_race_id],
                ['Challengee', ch.challengee_race_id],
                ['Status',     ch.status],
                ['Created At', formatDateTime(ch.created_at)],
            ]) {
                const tr = NewE('tr')
                const tdLabel = NewE('td')
                tdLabel.appendChild(NewT(label))
                tr.appendChild(tdLabel)
                const tdValue = NewE('td')
                tdValue.appendChild(NewT(value))
                tr.appendChild(tdValue)
                infoTable.appendChild(tr)
            }
            container.appendChild(infoTable)

            // Editable form
            let selectA = null
            let checkboxA = null
            let selectB = null
            let checkboxB = null

            const editTable = NewE('table')

            // Fleet Build A
            const trFbA = NewE('tr')
            const tdFbALabel = NewE('td')
            tdFbALabel.appendChild(NewT('Fleet Build A'))
            trFbA.appendChild(tdFbALabel)
            const tdFbA = NewE('td')
            if (canEditA) {
                selectA = NewE('select')
                selectA.name = 'fleet_build_a_id'
                const emptyOpt = NewE('option')
                emptyOpt.value = ''
                emptyOpt.appendChild(NewT('— select —'))
                selectA.appendChild(emptyOpt)
                for (const b of buildsA) {
                    const opt = NewE('option')
                    opt.value = b.id
                    opt.appendChild(NewT(b.name || b.id))
                    if (b.id === ch.fleet_build_a_id) opt.selected = true
                    selectA.appendChild(opt)
                }
                tdFbA.appendChild(selectA)
            } else {
                const input = NewE('input')
                input.type = 'text'
                input.readOnly = true
                input.value = fleetBuildA ? (fleetBuildA.name || fleetBuildA.id) : ch.fleet_build_a_id
                tdFbA.appendChild(input)
            }
            trFbA.appendChild(tdFbA)
            editTable.appendChild(trFbA)

            // Ready A
            const trReadyA = NewE('tr')
            const tdReadyALabel = NewE('td')
            tdReadyALabel.appendChild(NewT('Ready A'))
            trReadyA.appendChild(tdReadyALabel)
            const tdReadyA = NewE('td')
            if (canEditA) {
                checkboxA = NewE('input')
                checkboxA.type = 'checkbox'
                checkboxA.name = 'ready_a'
                checkboxA.checked = ch.ready_a
                tdReadyA.appendChild(checkboxA)
            } else {
                tdReadyA.appendChild(NewT(ch.ready_a ? '✓' : '✗'))
            }
            trReadyA.appendChild(tdReadyA)
            editTable.appendChild(trReadyA)

            // Fleet Build B
            const trFbB = NewE('tr')
            const tdFbBLabel = NewE('td')
            tdFbBLabel.appendChild(NewT('Fleet Build B'))
            trFbB.appendChild(tdFbBLabel)
            const tdFbB = NewE('td')
            if (canEditB) {
                selectB = NewE('select')
                selectB.name = 'fleet_build_b_id'
                const emptyOpt = NewE('option')
                emptyOpt.value = ''
                emptyOpt.appendChild(NewT('— select —'))
                selectB.appendChild(emptyOpt)
                for (const b of buildsB) {
                    const opt = NewE('option')
                    opt.value = b.id
                    opt.appendChild(NewT(b.name || b.id))
                    if (b.id === ch.fleet_build_b_id) opt.selected = true
                    selectB.appendChild(opt)
                }
                tdFbB.appendChild(selectB)
            } else {
                const input = NewE('input')
                input.type = 'text'
                input.readOnly = true
                input.value = fleetBuildB ? (fleetBuildB.name || fleetBuildB.id) : ch.fleet_build_b_id
                tdFbB.appendChild(input)
            }
            trFbB.appendChild(tdFbB)
            editTable.appendChild(trFbB)

            // Ready B
            const trReadyB = NewE('tr')
            const tdReadyBLabel = NewE('td')
            tdReadyBLabel.appendChild(NewT('Ready B'))
            trReadyB.appendChild(tdReadyBLabel)
            const tdReadyB = NewE('td')
            if (canEditB) {
                checkboxB = NewE('input')
                checkboxB.type = 'checkbox'
                checkboxB.name = 'ready_b'
                checkboxB.checked = ch.ready_b
                tdReadyB.appendChild(checkboxB)
            } else {
                tdReadyB.appendChild(NewT(ch.ready_b ? '✓' : '✗'))
            }
            trReadyB.appendChild(tdReadyB)
            editTable.appendChild(trReadyB)

            const saveBtn = NewE('button')
            saveBtn.type = 'submit'
            saveBtn.appendChild(NewT('Save'))

            const cancelLink = NewE('a')
            cancelLink.href = '/challenges.html'
            cancelLink.className = 'btn'
            cancelLink.style.marginLeft = '12px'
            cancelLink.appendChild(NewT('Cancel'))

            const form = NewE('form')
            form.appendChild(editTable)
            form.appendChild(saveBtn)
            form.appendChild(cancelLink)
            form.addEventListener('submit', async (e) => {
                e.preventDefault()
                const payload = {...ch}
                if (canEditA) {
                    payload.fleet_build_a_id = selectA.value
                    payload.ready_a = checkboxA.checked
                }
                if (canEditB) {
                    payload.fleet_build_b_id = selectB.value
                    payload.ready_b = checkboxB.checked
                }
                try {
                    await this.apiClient.updateChallenge(challengeId, payload)
                    this.dispatcher.dispatch('redirect', '/challenges.html')
                } catch (err) {
                    this.dispatcher.dispatch('displayError', [err.message, true])
                }
            })
            container.appendChild(form)

        } catch (e) {
            this.dispatcher.dispatch('displayError', [e.message, true])
        }

        return container
    }
}
