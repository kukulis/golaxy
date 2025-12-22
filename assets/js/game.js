// Galaktika Game - Main JavaScript

const SVG_NS = 'http://www.w3.org/2000/svg';

class Game {
    constructor() {
        this.svg = document.getElementById('galaxy');
        this.ships = [];
        this.init();
    }

    init() {
        // Set up event listeners
        document.getElementById('new-game').addEventListener('click', () => this.newGame());

        // Initial game setup
        this.newGame();
    }

    newGame() {
        // Clear existing game
        this.svg.innerHTML = '';
        this.ships = [];

        // Draw some example ships
        this.drawShip(100, 100, { name: 'Ship 1', speed: 10 });
        this.drawShip(300, 200, { name: 'Ship 2', speed: 8 });
        this.drawShip(500, 150, { name: 'Ship 3', speed: 12 });
    }

    drawShip(x, y, shipData) {
        // Create ship group
        const group = document.createElementNS(SVG_NS, 'g');
        group.setAttribute('class', 'ship');

        // Draw ship as a circle
        const circle = document.createElementNS(SVG_NS, 'circle');
        circle.setAttribute('cx', x);
        circle.setAttribute('cy', y);
        circle.setAttribute('r', 10);
        circle.setAttribute('fill', 'blue');
        circle.setAttribute('stroke', 'white');
        circle.setAttribute('stroke-width', 2);

        // Add ship label
        const text = document.createElementNS(SVG_NS, 'text');
        text.setAttribute('x', x);
        text.setAttribute('y', y - 15);
        text.setAttribute('text-anchor', 'middle');
        text.setAttribute('fill', 'white');
        text.setAttribute('font-size', '12');
        text.textContent = shipData.name;

        group.appendChild(circle);
        group.appendChild(text);

        // Add click handler
        group.addEventListener('click', () => this.handleShipClick(shipData));

        // Add hover effect
        group.style.cursor = 'pointer';

        this.svg.appendChild(group);
        this.ships.push({ data: shipData, element: group });
    }

    handleShipClick(shipData) {
        console.log('Ship clicked:', shipData);
        alert(`Ship: ${shipData.name}\nSpeed: ${shipData.speed}`);
    }
}

// Initialize game when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    const game = new Game();
});
