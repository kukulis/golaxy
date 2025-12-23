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
    }

    clearAllShotsFromDrawing() {
        // Clear previous shot visualization
        this.shotElements.forEach(element => element.remove());
        this.shotElements = [];
    }

    clearDestroyedShipsFromDrawing() {
        // Remove ships marked for removal from previous hit
        this.destroyedShipsIds.forEach(shipId => {
            this.removeShipFromDrawing(shipId);
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

        this.shotElements.push( shot.buildSvg())
        this.svg.appendChild(shot.svgElement);

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

    removeShipFromDrawing(shipId) {
        const foundShip = this.battle.findShip(shipId)
        if (foundShip && foundShip.svgElement) {
            foundShip.svgElement.remove();
        }
    }

}


