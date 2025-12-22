// Galaktika Game - Main JavaScript

const SVG_NS = 'http://www.w3.org/2000/svg';

class Game {
    constructor() {
        this.svg = document.getElementById('galaxy');
        this.ships = [];
        this.apiClient = new ApiClient();
        this.battle = null;
        this.currentShotIndex = 0;
        this.shipsToRemove = new Set();
        this.shotElements = [];
        this.isAutoPlaying = false;
        this.init();
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
        this.shipsToRemove.clear();
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
        // Create ship group
        const group = document.createElementNS(SVG_NS, 'g');
        group.setAttribute('class', 'ship');

        // Draw ship as a triangle
        const triangle = document.createElementNS(SVG_NS, 'polygon');

        // Create triangle points based on side
        let points;
        if (side === 'a') {
            // Right-pointing triangle for side A
            points = `${x-10},${y-10} ${x+15},${y} ${x-10},${y+10}`;
        } else {
            // Left-pointing triangle for side B
            points = `${x+10},${y-10} ${x-15},${y} ${x+10},${y+10}`;
        }

        triangle.setAttribute('points', points);
        triangle.setAttribute('fill', ship.destroyed ? 'gray' : color);
        triangle.setAttribute('stroke', 'white');
        triangle.setAttribute('stroke-width', 2);

        // Add ship label
        const text = document.createElementNS(SVG_NS, 'text');
        text.setAttribute('x', x);
        text.setAttribute('y', y - 20);
        text.setAttribute('text-anchor', 'middle');
        text.setAttribute('fill', 'white');
        text.setAttribute('font-size', '12');
        text.textContent = ship.name;

        group.appendChild(triangle);
        group.appendChild(text);

        // Add click handler
        group.addEventListener('click', () => this.handleShipClick(ship));

        // Add hover effect
        group.style.cursor = 'pointer';

        this.svg.appendChild(group);
        this.ships.push({ data: ship, element: group });
    }

    handleShipClick(ship) {
        console.log('Ship clicked:', ship);
        const info = `Ship: ${ship.name}\nID: ${ship.id}\nSpeed: ${ship.tech.speed}\nAttack: ${ship.tech.attack}\nDefense: ${ship.tech.defense}\nDestroyed: ${ship.destroyed}`;
        alert(info);
    }

    processShot() {
        // Clear previous shot visualization
        this.shotElements.forEach(element => element.remove());
        this.shotElements = [];

        // Remove ships marked for removal from previous hit
        this.shipsToRemove.forEach(shipId => {
            this.removeShip(shipId);
        });
        this.shipsToRemove.clear();

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
            this.shipsToRemove.add(shot.destination);
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
        const shipEntry = this.ships.find(s => s.data.id === shipId);
        if (shipEntry && shipEntry.element) {
            shipEntry.element.remove();
            this.ships = this.ships.filter(s => s.data.id !== shipId);
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

// Initialize game when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    const game = new Game();
});
