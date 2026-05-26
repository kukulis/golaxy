import {NewE, NewT} from './helper.js';
import {Division} from './entities/division.js';

/**
 * @param {Division} division
 * @returns {HTMLElement}
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

    const tbody = NewE('tbody');
    for (const [label, value] of rows) {
        const tr = NewE('tr');
        const th = NewE('th');
        const td = NewE('td');
        th.appendChild(NewT(label));
        td.appendChild(NewT(value));
        tr.appendChild(th);
        tr.appendChild(td);
        tbody.appendChild(tr);
    }

    const table = NewE('table');
    table.appendChild(tbody);
    table.style.display = 'none';

    const button = NewE('button');
    button.appendChild(NewT('▶'));
    button.addEventListener('click', () => {
        if (table.style.display === 'none') {
            table.style.display = '';
            button.textContent = '▼';
        } else {
            table.style.display = 'none';
            button.textContent = '▶';
        }
    });

    const wrapper = NewE('div');
    wrapper.appendChild(button);
    wrapper.appendChild(table);
    return wrapper;
}
