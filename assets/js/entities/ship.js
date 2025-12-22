class Ship {
    id = '';
    tech = {
        attack: 0,
        guns: 0,
        defense: 0,
        speed: 0,
        cargo_capacity: 0,
        mass: 0
    };
    destroyed = false;
    name = '';
    owner = '';

    // not received from API
    battleX = 0;
    battleY = 0;
}