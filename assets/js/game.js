// Galaktika Game - Main JavaScript

const SVG_NS = 'http://www.w3.org/2000/svg';

class Game {

    /**
     * @type {HTMLElement}
     */
    svg = null;


    /**
     * @type {ApiClient}
     */
    apiClient = null;

    /**
     *
     * @type {Battle}
     */
    battle = null;

    /**
     *
     * @type {boolean}
     */
    isAutoPlaying = false

    /**
     * @type {number}
     */
    currentShotIndex = 0;

    /**
     * Array of SVG elements (lines, circles) representing shot visualizations
     * @type {SVGElement[]}
     */
    shotElements = [];

    /**
     * Set of ship IDs that will be removed on the next shot
     * @type {Set<string|number>}
     */
    destroyedShipsIds = new Set();

    /**
     *
     * TODO remove this property.
     * Array of ship objects containing ship data and their SVG group elements
     * @type {Ship[]}>}
     */
    ships = [];


    /**
     * @param apiClient {ApiClient}
     * @param svg {HTMLElement}
     */
    constructor(apiClient, svg) {
        this.svg = svg;
        this.apiClient = apiClient;
    }

    init() {
        // Set up event listeners
        document.getElementById('new-game').addEventListener('click', () => this.newGame());
        document.getElementById('shot-button').addEventListener('click', () => this.processShot());
        document.getElementById('auto-button').addEventListener('click', () => this.autoPlay());

        // Initial game setup
        this.newGame();
    }

    async newGame() {
        // Clear existing game
        this.svg.innerHTML = '';
        this.ships = [];
        this.currentShotIndex = 0;
        this.destroyedShipsIds.clear();
        this.shotElements = [];
        this.isAutoPlaying = false;

        try {
            // Load battle from API
            this.battle = await this.apiClient.getBattle('1');

            // Position and draw ships
            this.positionShips();
            this.drawBattle();
        } catch (error) {
            console.error('Failed to load battle:', error);
            alert('Failed to load battle. Please try again.');
        }
    }

    positionShips() {
        const leftX = 100;
        const rightX = 700;
        const startY = 100;
        const shipSpacing = 60;

        // Position side A ships (left column)
        if (this.battle.side_a && this.battle.side_a.ships) {
            this.battle.side_a.ships.forEach((ship, index) => {
                ship.battleX = leftX;
                ship.battleY = startY + (index * shipSpacing);
            });
        }

        // Position side B ships (right column)
        if (this.battle.side_b && this.battle.side_b.ships) {
            this.battle.side_b.ships.forEach((ship, index) => {
                ship.battleX = rightX;
                ship.battleY = startY + (index * shipSpacing);
            });
        }
    }

    drawBattle() {
        // Draw all ships from side A
        if (this.battle.side_a && this.battle.side_a.ships) {
            this.battle.side_a.ships.forEach(ship => {
                this.drawShip(ship.battleX, ship.battleY, ship, 'blue', 'a');
            });
        }

        // Draw all ships from side B
        if (this.battle.side_b && this.battle.side_b.ships) {
            this.battle.side_b.ships.forEach(ship => {
                this.drawShip(ship.battleX, ship.battleY, ship, 'red', 'b');
            });
        }
    }

    drawShip(x, y, ship, color, side) {
        this.svg.appendChild(ship.creteShipSvg(color, side));

        // TODO remove this
        this.ships.push(ship);
    }

    clearAllShotsFromDrawing() {
        // Clear previous shot visualization
        this.shotElements.forEach(element => element.remove());
        this.shotElements = [];
    }

    clearDestroyedShipsFromDrawing() {
        // Remove ships marked for removal from previous hit
        this.destroyedShipsIds.forEach(shipId => {
            this.removeShip(shipId);
        });
        this.destroyedShipsIds.clear();
    }



    processShot() {
        this.clearAllShotsFromDrawing();
        this.clearDestroyedShipsFromDrawing();

        // Check if there are more shots
        if (this.currentShotIndex >= this.battle.shots.length) {
            if (!this.isAutoPlaying) {
                alert('No more shots available!');
            }
            return false;
        }

        // Get the next shot
        const shot = this.battle.shots[this.currentShotIndex];
        this.currentShotIndex++;

        // Find source and destination ships
        const sourceShip = this.findShip(shot.source);
        const destShip = this.findShip(shot.destination);

        if (!sourceShip || !destShip) {
            console.error('Could not find ships for shot:', shot);
            return false;
        }

        // Draw the shot line
        this.drawLine(sourceShip.battleX, sourceShip.battleY, destShip.battleX, destShip.battleY);

        // Draw miss or hit indicator
        if (shot.result) {
            // Hit - draw cross
            this.drawCross(destShip.battleX, destShip.battleY);
            // Mark ship for removal on next shot
            this.destroyedShipsIds.add(shot.destination);
        } else {
            // Miss - draw circle
            this.drawCircle(destShip.battleX, destShip.battleY);
        }

        return true;
    }

    autoPlay() {
        // Prevent multiple auto-plays
        if (this.isAutoPlaying) {
            return;
        }

        // Check if there are shots to play
        if (this.currentShotIndex >= this.battle.shots.length) {
            alert('No more shots available!');
            return;
        }

        this.isAutoPlaying = true;

        // Disable buttons during auto-play
        document.getElementById('shot-button').disabled = true;
        document.getElementById('auto-button').disabled = true;

        this.executeNextShot();
    }

    executeNextShot() {
        const hasMoreShots = this.processShot();

        if (hasMoreShots && this.currentShotIndex < this.battle.shots.length) {
            // Schedule next shot after 1 second
            setTimeout(() => this.executeNextShot(), 1000);
        } else {
            // Auto-play finished
            this.isAutoPlaying = false;
            document.getElementById('shot-button').disabled = false;
            document.getElementById('auto-button').disabled = false;
        }
    }

    // TODO move this method to the Fleet class
    findShip(shipId) {
        // Search in side A
        if (this.battle.side_a && this.battle.side_a.ships) {
            const ship = this.battle.side_a.ships.find(s => s.id === shipId);
            if (ship) return ship;
        }

        // Search in side B
        if (this.battle.side_b && this.battle.side_b.ships) {
            const ship = this.battle.side_b.ships.find(s => s.id === shipId);
            if (ship) return ship;
        }

        return null;
    }

    removeShip(shipId) {
        const foundShip = this.ships.find(s => s.id === shipId);
        if (foundShip && foundShip.svgElement) {
            foundShip.svgElement.remove();
            this.ships = this.ships.filter(s => s.id !== shipId);
        }
    }

    drawLine(x1, y1, x2, y2) {
        const line = document.createElementNS(SVG_NS, 'line');
        line.setAttribute('x1', x1);
        line.setAttribute('y1', y1);
        line.setAttribute('x2', x2);
        line.setAttribute('y2', y2);
        line.setAttribute('stroke', 'yellow');
        line.setAttribute('stroke-width', 2);
        this.svg.appendChild(line);
        this.shotElements.push(line);
    }

    drawCircle(x, y) {
        const circle = document.createElementNS(SVG_NS, 'circle');
        circle.setAttribute('cx', x);
        circle.setAttribute('cy', y);
        circle.setAttribute('r', 20);
        circle.setAttribute('fill', 'none');
        circle.setAttribute('stroke', 'white');
        circle.setAttribute('stroke-width', 2);
        this.svg.appendChild(circle);
        this.shotElements.push(circle);
    }

    drawCross(x, y) {
        // Draw X with two lines
        const line1 = document.createElementNS(SVG_NS, 'line');
        line1.setAttribute('x1', x - 15);
        line1.setAttribute('y1', y - 15);
        line1.setAttribute('x2', x + 15);
        line1.setAttribute('y2', y + 15);
        line1.setAttribute('stroke', 'red');
        line1.setAttribute('stroke-width', 3);

        const line2 = document.createElementNS(SVG_NS, 'line');
        line2.setAttribute('x1', x + 15);
        line2.setAttribute('y1', y - 15);
        line2.setAttribute('x2', x - 15);
        line2.setAttribute('y2', y + 15);
        line2.setAttribute('stroke', 'red');
        line2.setAttribute('stroke-width', 3);

        this.svg.appendChild(line1);
        this.svg.appendChild(line2);
        this.shotElements.push(line1);
        this.shotElements.push(line2);
    }
}


