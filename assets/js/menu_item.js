export class MenuItem {

    /**
     * @type {string}
     */
    name = '';

    /**
     * @type {string}
     */
    link = '';

    /**
     * @type {MenuItem[]}
     */
    children = [];

    /**
     * @param {string} name
     * @param {string} link
     * @param {MenuItem[]} children
     */
    constructor(name, link, children = []) {
        this.name = name
        this.link = link
        this.children = children
    }
}
