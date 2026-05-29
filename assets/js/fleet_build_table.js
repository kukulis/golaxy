import {NewE, NewT} from './helper.js';
import {FleetBuild} from './entities/fleet_build.js';

/**
 * @param {FleetBuild} b
 * @param {Object|null} technologies
 * @returns {HTMLElement}
 */
export function createFleetBuildStatisticsTable(b, technologies = null) {
    if (technologies === null) {
        technologies = {
            attack: 1,
            defense: 1,
            engine: 1,
            cargo: 1,
        }
    }
    // const tech = (value, key) => technologies ? `${value} (${technologies[key]})` : value;
    const rows = [
        ['ID', b.id],
        ['Division', b.division_id],
        ['Race', b.race_id],
        ['Used Attack Resources', b.attack_resources + ' (' + technologies.attack + ')'],
        ['Used Defense Resources', b.defense_resources + ' (' + technologies.defense + ')'],
        ['Used Engine Resources', b.engine_resources + ' (' + technologies.engine + ')'],
        ['Used Cargo Resources', b.cargo_resources + ' (' + technologies.cargo + ')'],
        ['Resources Used', b.usedResources],
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
