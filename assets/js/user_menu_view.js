import {NewE, NewT} from '/assets/js/helper.js'
import {MenuItem} from './menu_item.js'

export default class UserMenuView {

    /**
     * @return {HTMLElement}
     */
    generateView() {
        const divisionId = document.getElementById('user-menu')?.dataset.divisionId

        const fleetBuildId = document.getElementById('user-menu')?.dataset.fleetBuild

        const divisionsItem = new MenuItem('Divisions', '/divisions.html')
        if (divisionId) {
            const divisionItem = new MenuItem(`Division: ${divisionId}`, `/division/${divisionId}/main.html`)
            if (fleetBuildId !== undefined) {
                const fleetBuildsItem = new MenuItem('Fleet builds', `/division/${divisionId}/fleet-builds.html`)
                if (fleetBuildId && fleetBuildId !== '0') {
                    fleetBuildsItem.children.push(new MenuItem(`Fleet Build: ${fleetBuildId}`, `/fleet-build/${fleetBuildId}/main.html`))
                }
                divisionItem.children.push(fleetBuildsItem)
            }
            divisionsItem.children.push(divisionItem)
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
        if (!root) container.className = 'nav-submenu'

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
        wrapper.className = 'nav-item'

        if (item.link && window.location.pathname !== item.link) {
            const a = NewE('a')
            a.href = item.link
            a.appendChild(NewT(item.name))
            wrapper.appendChild(a)
        } else {
            wrapper.appendChild(NewT(item.name))
        }

        if (item.children.length > 0) {
            wrapper.appendChild(this.#buildDom(item.children, false))
        }

        return wrapper
    }
}
