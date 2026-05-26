import {NewE, NewT} from '/assets/js/helper.js'
import {MenuItem} from './menu_item.js'

export default class UserMenuView {

    /**
     * @return {HTMLElement}
     */
    generateView() {
        const divisionId = document.getElementById('user-menu')?.dataset.divisionId

        const divisionsItem = new MenuItem('Divisions', '/divisions.html')
        if (divisionId) {
            divisionsItem.children.push(new MenuItem(`Division: ${divisionId}`, `/division/${divisionId}/main.html`))
        }

        const items = [
            new MenuItem('Home', '/'),
            divisionsItem,
            new MenuItem('Races', '/races.html'),
            new MenuItem('Ship Models', '/ship-model/list.html'),
        ]

        return this.#buildDom(items)
    }

    /**
     * @param {MenuItem[]} items
     * @param {boolean} root
     * @return {HTMLElement}
     */
    #buildDom(items, root = true) {
        const container = NewE(root ? 'nav' : 'span')

        for (let i = 0; i < items.length; i++) {
            if (i > 0) container.appendChild(NewT(' | '))
            container.appendChild(this.#buildItemDom(items[i]))
        }

        return container
    }

    /**
     * @param {MenuItem} item
     * @return {HTMLElement}
     */
    #buildItemDom(item) {
        const wrapper = NewE('span')

        if (item.link && window.location.pathname !== item.link) {
            const a = NewE('a')
            a.href = item.link
            a.appendChild(NewT(item.name))
            wrapper.appendChild(a)
        } else {
            wrapper.appendChild(NewT(item.name))
        }

        if (item.children.length > 0) {
            wrapper.appendChild(NewT(' ('))
            wrapper.appendChild(this.#buildDom(item.children, false))
            wrapper.appendChild(NewT(')'))
        }

        return wrapper
    }
}
