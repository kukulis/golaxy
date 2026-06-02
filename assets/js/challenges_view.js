import {ClearE, NewE, NewT} from '/assets/js/helper.js'
import {ApiClient} from './api.js'
import {ChallengesFilter} from './challenges_filter.js'
import {Dispatcher} from './dispatcher.js'
import {formatDateTime} from './date_format.js'

export default class ChallengesView {

    /**
     * @type {ApiClient}
     */
    apiClient = null

    /**
     * @type {Dispatcher}
     */
    dispatcher = null

    /**
     * @type {HTMLElement}
     */
    tBody = null

    /**
     * @type {string}
     */
    currentRaceId = ''

    /**
     * @type {string}
     */
    currentRaceRole = ''

    /**
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
        const viewDiv = NewE('div')

        try {
            const race = await this.apiClient.getCurrentRace()
            this.currentRaceId = race.id
            this.currentRaceRole = race.role
        } catch (e) {
            this.dispatcher.dispatch('displayError', [e.message, true])
        }

        const table = NewE('table')

        const thead = NewE('thead')
        const headerRow = NewE('tr')
        for (const label of ['', 'ID', 'Status', 'Division', 'Challenger', 'Challengee', 'Fleet Build A', 'Ready A', 'Fleet Build B', 'Ready B', 'Created At', 'Battle Report', '']) {
            const th = NewE('th')
            th.appendChild(NewT(label))
            headerRow.appendChild(th)
        }
        thead.appendChild(headerRow)
        table.appendChild(thead)

        this.tBody = NewE('tbody')
        table.appendChild(this.tBody)
        viewDiv.appendChild(table)

        await this.reloadTableBody()

        return viewDiv
    }

    async reloadTableBody() {
        ClearE(this.tBody)

        try {
            const challenges = await this.loadChallenges()
            for (const ch of challenges) {
                const tr = NewE('tr')

                const tdEdit = NewE('td')
                if (ch.challenger_race_id === this.currentRaceId || ch.challengee_race_id === this.currentRaceId) {
                    const editLink = NewE('a')
                    editLink.href = `/challenge/${ch.id}/edit.html`
                    editLink.appendChild(NewT('✏ Edit'))
                    tdEdit.appendChild(editLink)
                }
                tr.appendChild(tdEdit)

                const columns = [
                    ch.id,
                    ch.status,
                    ch.division_id,
                    ch.challenger_race_id,
                    ch.challengee_race_id,
                    ch.fleet_build_a_id,
                    ch.ready_a ? '✓' : '✗',
                    ch.fleet_build_b_id,
                    ch.ready_b ? '✓' : '✗',
                    formatDateTime(ch.created_at),
                    // ch.created_at ?? '',
                    ch.battle_report_id
                ]
                for (const val of columns) {
                    const td = NewE('td')
                    td.appendChild(NewT(val))
                    tr.appendChild(td)
                }

                const tdActions = NewE('td')
                if (ch.challenger_race_id === this.currentRaceId) {
                    const deleteBtn = NewE('button')
                    deleteBtn.appendChild(NewT('Delete'))
                    deleteBtn.addEventListener('click', async () => {
                        if (!confirm(`Delete challenge ${ch.id}?`)) return
                        try {
                            await this.apiClient.deleteChallenge(ch.id)
                            await this.reloadTableBody()
                        } catch (e) {
                            this.dispatcher.dispatch('displayError', [e.message, true])
                        }
                    })
                    tdActions.appendChild(deleteBtn)
                }
                tr.appendChild(tdActions)

                this.tBody.appendChild(tr)
            }
        } catch (e) {
            this.dispatcher.dispatch('displayError', [e.message, true])
        }
    }

    async loadChallenges() {
        if (this.currentRaceRole === 'admin') {
            return this.apiClient.getChallenges(new ChallengesFilter())
        }

        const asChallenger = new ChallengesFilter()
        asChallenger.challenger_id = this.currentRaceId

        const asChallengee = new ChallengesFilter()
        asChallengee.challengee_id = this.currentRaceId

        const [sent, received] = await Promise.all([
            this.apiClient.getChallenges(asChallenger),
            this.apiClient.getChallenges(asChallengee),
        ])

        const merged = new Map()
        for (const ch of [...sent, ...received]) {
            merged.set(ch.id, ch)
        }
        return [...merged.values()]
    }
}
