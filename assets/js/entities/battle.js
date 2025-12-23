class Battle {
    /**
     * Battle identifier
     * @type {string}
     */
    id = '';

    /**
     * First fleet participating in the battle
     * @type {Fleet|null}
     */
    side_a = null;

    /**
     * Second fleet participating in the battle
     * @type {Fleet|null}
     */
    side_b = null;

    /**
     * Array of all shots fired during the battle
     * @type {Shot[]}
     */
    shots = [];

    findShip(shipId) {
        let ship = this.side_a.findShip(shipId)
        if ( ship != null ) {
            return ship;
        }

        return this.side_b.findShip(shipId)
    }

    fixShotsReferences () {
        for ( let shot of  this.shots) {
            shot.sourceShip = this.findShip(shot.source)
            shot.destinationShip = this.findShip(shot.destination)
        }
    }
}
