class ApiClient {

    /**
     * @param id
     * @returns {Promise<Battle>}
     */
    async getBattle(id) {
        const response = await fetch('/api/battle');

        if (!response.ok) {
            throw new Error(`Failed to fetch battle: ${response.statusText}`);
        }

        const data = await response.json();

        return this.convertToBattle(data);
    }

    convertToBattle(data) {
        const battle = new Battle();
        battle.id = data.id;
        battle.side_a = this.convertToFleet(data.side_a);
        battle.side_b = this.convertToFleet(data.side_b);
        battle.shots = data.shots.map(shot => this.convertToShot(shot));
        battle.fixShotsReferences()

        return battle;
    }

    convertToFleet(data) {
        if (!data) return null;

        const fleet = new Fleet();
        fleet.ships = data.ships.map(ship => this.convertToShip(ship));
        fleet.fillShipsMap(fleet.ships)
        fleet.owner = data.owner;

        return fleet;
    }

    convertToShip(data) {
        const ship = new Ship();
        ship.id = data.id;
        ship.tech = {
            attack: data.tech.attack,
            guns: data.tech.guns,
            defense: data.tech.defense,
            speed: data.tech.speed,
            cargo_capacity: data.tech.cargo_capacity,
            mass: data.tech.mass
        };
        ship.destroyed = data.destroyed;
        ship.name = data.name;
        ship.owner = data.owner;

        return ship;
    }

    convertToShot(data) {
        const shot = new Shot();
        shot.source = data.source;
        shot.destination = data.destination;
        shot.result = data.result;

        return shot;
    }
}