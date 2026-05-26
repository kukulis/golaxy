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
        bar.style.cssText = 'position:fixed;top:12px;right:20px;color:#4a9eff;font-size:14px'

        const loginName = this.dispatcher.dispatch('getLoginName')
        bar.appendChild(NewT(loginName ?? ''))

        return bar
    }
}
