class Fleet {
    /**
     * Array of ships in this fleet
     * @type {Ship[]}
     */
    ships = [];

    /**
     * Fleet owner identifier
     * @type {string}
     */
    owner = '';

    /**
     *
     * @type {Map<string, Ship>}
     */
    shipsMap = new Map();

    /**
     *
     * @param ships {Ship[]}
     */
    fillShipsMap(ships) {
        for ( let ship of ships ) {
            this.shipsMap.set(ship.id, ship)
        }
    }

    findShip(shipId) {
        return this.shipsMap.get(shipId)
    }
}
