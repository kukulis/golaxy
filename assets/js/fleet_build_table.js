import {NewE, NewT} from './helper.js';
import {FleetBuild} from './entities/fleet_build.js';

/**
 * @param {FleetBuild} b
 * @param {Object|null} technologies
 * @returns {HTMLElement}
 */
export function createFleetBuildStatisticsTable(b, technologies = null) {
    const tech = (value, key) => technologies ? `${value} (${technologies[key]})` : value;
    const rows = [
        ['ID',                    b.id],
        ['Division',              b.division_id],
        ['Race',                  b.race_id],
        ['Used Attack Resources', tech(b.attack_resources, 'attack')],
        ['Used Defense Resources',tech(b.defense_resources, 'defense')],
        ['Used Engine Resources', tech(b.engine_resources, 'engine')],
        ['Used Cargo Resources',  tech(b.cargo_resources, 'cargo')],
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
