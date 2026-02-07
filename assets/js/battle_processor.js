// Galaktika Battle Processor

import { ApiClient } from './api.js';

const SVG_NS = 'http://www.w3.org/2000/svg';

export class BattleProcessor {

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
     * Ships that will be processed as destroyed on the next shot
     * @type {Ship[]}
     */
    destroyedShips = [];


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
        this.destroyedShips = [];
        this.shotElements = [];
        this.isAutoPlaying = false;

        try {
            // Load battle from API
            this.battle = await this.apiClient.getBattle('1');

            // Build ship groups
            this.battle.side_a.fillShipGroupMap();
            this.battle.side_b.fillShipGroupMap();

            // Position and draw ship groups
            this.positionShipGroups();
            this.drawBattle();
        } catch (error) {
            console.error('Failed to load battle:', error);
            alert('Failed to load battle. Please try again.');
        }
    }

    positionShipGroups() {
        const leftX = 100;
        const rightX = 700;
        const startY = 100;
        const groupPadding = 20;

        // Position side A ship groups (left column)
        if (this.battle.side_a) {
            let currentY = startY;
            for (const group of this.battle.side_a.shipGroupMap.values()) {
                group.battleX = leftX;
                group.battleY = currentY;
                for (const ship of group.shipList) {
                    ship.battleX = group.battleX;
                    ship.battleY = group.battleY;
                }
                currentY += group.calculateHeight() + groupPadding;
            }
        }

        // Position side B ship groups (right column)
        if (this.battle.side_b) {
            let currentY = startY;
            for (const group of this.battle.side_b.shipGroupMap.values()) {
                group.battleX = rightX;
                group.battleY = currentY;
                for (const ship of group.shipList) {
                    ship.battleX = group.battleX;
                    ship.battleY = group.battleY;
                }
                currentY += group.calculateHeight() + groupPadding;
            }
        }
    }

    drawBattle() {
        // Draw all ship groups from side A
        if (this.battle.side_a) {
            for (const group of this.battle.side_a.shipGroupMap.values()) {
                this.drawShipGroup(group, 'blue', 'a');
            }
        }

        // Draw all ship groups from side B
        if (this.battle.side_b) {
            for (const group of this.battle.side_b.shipGroupMap.values()) {
                this.drawShipGroup(group, 'red', 'b');
            }
        }
    }

    drawShipGroup(shipGroup, color, side) {
        const svgElement = shipGroup.createGroupSvg(color, side);

        // Add hover effect
        svgElement.style.cursor = 'pointer';

        this.svg.appendChild(svgElement);
    }

    clearAllShotsFromDrawing() {
        // Clear previous shot visualization
        this.shotElements.forEach(element => element.remove());
        this.shotElements = [];
    }

    clearDestroyedShipsFromDrawing() {
        for (const ship of this.destroyedShips) {
            ship.destroyed = true;

            // Find the ship's group and notify
            const group = this.battle.side_a.findShipGroup(ship)
                || this.battle.side_b.findShipGroup(ship);

            if (group) {
                // Find the fleet to call notifyDestroyed
                const fleet = this.battle.side_a.findShipGroup(ship)
                    ? this.battle.side_a
                    : this.battle.side_b;
                fleet.notifyDestroyed(ship);

                // Remove group SVG if amount is 0
                if (group.amount === 0 && group.svgElement) {
                    group.svgElement.remove();
                }
            }
        }
        this.destroyedShips = [];
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

        if ( shot.result ) {
            this.destroyedShips.push(shot.destinationShip);
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

}


