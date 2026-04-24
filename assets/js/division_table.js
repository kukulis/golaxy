import {Division} from './entities/division.js';

/**
 * @param {Division} division
 * @returns {HTMLTableElement}
 */
export function createDivisionTable(division) {
    const rows = [
        ['ID',          division.id],
        ['Resources',   division.resources_amount],
        ['Tech Attack', division.tech_attack],
        ['Tech Defense',division.tech_defense],
        ['Tech Engines',division.tech_engines],
        ['Tech Cargo',  division.tech_cargo],
    ];

    const tbody = document.createElement('tbody');
    for (const [label, value] of rows) {
        const tr = document.createElement('tr');
        const th = document.createElement('th');
        const td = document.createElement('td');
        th.appendChild(document.createTextNode(label));
        td.appendChild(document.createTextNode(value));
        tr.appendChild(th);
        tr.appendChild(td);
        tbody.appendChild(tr);
    }

    const table = document.createElement('table');
    table.appendChild(tbody);
    return table;
}
