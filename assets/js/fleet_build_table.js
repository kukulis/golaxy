import {NewE, NewT} from './helper.js';
import {FleetBuild} from './entities/fleet_build.js';

/**
 * @param {FleetBuild} b
 * @returns {HTMLElement}
 */
export function createFleetBuildStatisticsTable(b) {
    const rows = [
        ['ID',                    b.id],
        ['Division',              b.division_id],
        ['Race',                  b.race_id],
        ['Used Attack Resources', b.attack_resources],
        ['Used Defense Resources',b.defense_resources],
        ['Used Engine Resources', b.engine_resources],
        ['Used Cargo Resources',  b.cargo_resources],
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
