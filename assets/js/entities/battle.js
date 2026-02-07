import { Fleet } from './fleet.js';
import { Shot } from './shot.js';

export class Battle {
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

    /**
     * Updates battle properties from DTO data
     * @param {Object} data - Battle data from API
     */
    updateFromDTO(data) {
        this.id = data.id;

        if (data.side_a) {
            this.side_a = new Fleet();
            this.side_a.updateFromDTO(data.side_a);
        } else {
            this.side_a = null;
        }

        if (data.side_b) {
            this.side_b = new Fleet();
            this.side_b.updateFromDTO(data.side_b);
        } else {
            this.side_b = null;
        }

        this.shots = data.shots.map(shotData => {
            const shot = new Shot();
            shot.updateFromDTO(shotData);
            return shot;
        });

        this.fixShotsReferences();

        return this;
    }
}
