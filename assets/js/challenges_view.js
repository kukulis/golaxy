import {ClearE, NewE, NewT} from '/assets/js/helper.js'
import {ApiClient} from './api.js'
import {Dispatcher} from './dispatcher.js'

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
        } catch (e) {
            this.dispatcher.dispatch('displayError', [e.message, true])
        }

        const table = NewE('table')

        const thead = NewE('thead')
        const headerRow = NewE('tr')
        for (const label of ['ID', 'Challenger Race', 'Fleet Build A', 'Fleet Build B', 'Battle Report', '']) {
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
            const challenges = await this.apiClient.getChallenges()
            for (const ch of challenges) {
                const tr = NewE('tr')

                for (const val of [ch.id, ch.challenger_race_id, ch.fleet_build_a_id, ch.fleet_build_b_id, ch.battle_report_id]) {
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
}
