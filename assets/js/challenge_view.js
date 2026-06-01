import {NewE, NewT} from '/assets/js/helper.js'
import {ApiClient} from './api.js'
import {Dispatcher} from './dispatcher.js'

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
            const ch = await this.apiClient.getChallenge(challengeId)
            const [ownBuilds, allBuilds] = await Promise.all([
                this.apiClient.getFleetBuilds(ch.division_id, false),
                this.apiClient.getFleetBuilds(ch.division_id, true),
            ])

            const readOnlyTable = NewE('table')
            for (const [label, value] of [
                ['ID', ch.id],
                ['Challenger Race', ch.challenger_race_id],
            ]) {
                const tr = NewE('tr')
                const tdLabel = NewE('td')
                tdLabel.appendChild(NewT(label))
                tr.appendChild(tdLabel)
                const tdValue = NewE('td')
                tdValue.appendChild(NewT(value))
                tr.appendChild(tdValue)
                readOnlyTable.appendChild(tr)
            }
            container.appendChild(readOnlyTable)

            const editTable = NewE('table')

            for (const [label, name, builds, currentId] of [
                ['Fleet Build A', 'fleet_build_a_id', ownBuilds,  ch.fleet_build_a_id],
                ['Fleet Build B', 'fleet_build_b_id', allBuilds,  ch.fleet_build_b_id],
            ]) {
                const tr = NewE('tr')
                const tdLabel = NewE('td')
                tdLabel.appendChild(NewT(label))
                tr.appendChild(tdLabel)
                const tdInput = NewE('td')
                const select = NewE('select')
                select.name = name
                const emptyOption = NewE('option')
                emptyOption.value = ''
                emptyOption.appendChild(NewT('— select —'))
                select.appendChild(emptyOption)
                for (const b of builds) {
                    const option = NewE('option')
                    option.value = b.id
                    option.appendChild(NewT(b.name || b.id))
                    if (b.id === currentId) option.selected = true
                    select.appendChild(option)
                }
                tdInput.appendChild(select)
                tr.appendChild(tdInput)
                editTable.appendChild(tr)
            }

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
                try {
                    const data = Object.fromEntries(new FormData(form))
                    await this.apiClient.updateChallenge(challengeId, {...ch, ...data})
                    this.dispatcher.dispatch('redirect', '/challenges.html')
                } catch (e) {
                    this.dispatcher.dispatch('displayError', [e.message, true])
                }
            })
            container.appendChild(form)
        } catch (e) {
            this.dispatcher.dispatch('displayError', [e.message, true])
        }

        return container
    }
}
