import {NewE, NewT} from '/assets/js/helper.js'
import {Dispatcher} from "./dispatcher.js";

export default class UserBarView {

    /**
     *
     * @type {Dispatcher}
     */
    dispatcher = null;

    /**
     *
     * @param {Dispatcher} dispatcher
     */
    constructor(dispatcher) {
        this.dispatcher = dispatcher
    }

    /**
     * @return {HTMLElement}
     */
    generateView() {
        const bar = NewE('div')
        bar.id = 'user-bar'

        const loginName = this.dispatcher.dispatch('getLoginName')

        if (loginName) {
            bar.appendChild(NewT(loginName))

            const logoutLink = NewE('a')
            logoutLink.href = '#'
            logoutLink.appendChild(NewT('Logout'))
            logoutLink.addEventListener('click', (e) => {
                e.preventDefault()
                this.dispatcher.dispatch('logout')
            })

            bar.appendChild(NewT(' | '))
            bar.appendChild(logoutLink)
        } else {
            const loginLink = NewE('a')
            loginLink.href = '/dummy_login.html'
            loginLink.appendChild(NewT('Login'))
            bar.appendChild(loginLink)
        }

        return bar
    }
}
