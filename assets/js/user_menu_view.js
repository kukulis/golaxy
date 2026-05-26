import {NewE, NewT} from '/assets/js/helper.js'

export default class UserMenuView {

    /**
     * @return {HTMLElement}
     */
    generateView() {
        const menu = NewE('nav')

        const links = [
            {label: 'Ship Models', href: '/ship-model/list.html'},
            {label: 'Divisions', href: '/divisions.html'},
        ]

        for (let i = 0; i < links.length; i++) {
            if (i > 0) {
                menu.appendChild(NewT(' | '))
            }
            const a = NewE('a')
            a.href = links[i].href
            a.appendChild(NewT(links[i].label))
            menu.appendChild(a)
        }

        return menu
    }
}
