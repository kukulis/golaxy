import {NewE, NewT} from './helper.js';
import {FleetBuild} from './entities/fleet_build.js';
import {FleetBuildStatistics} from './entities/fleet_build_statistics.js';

/**
 * @param {FleetBuild} b
 * @param {FleetBuildStatistics} statistics
 * @returns {HTMLElement}
 */
export function createFleetBuildStatisticsTable(b, statistics = null ) {
    // const tech = (value, key) => technologies ? `${value} (${technologies[key]})` : value;
    const rows = [
        ['ID', b.id],
        ['Division', b.division_id],
        ['Race', b.race_id],
        ['Name', b.name],
        ['Used Attack Resources', b.attack_resources + ' (' + statistics.technologies.attack + ')'],
        ['Used Defense Resources', b.defense_resources + ' (' + statistics.technologies.defense + ')'],
        ['Used Engine Resources', b.engine_resources + ' (' + statistics.technologies.engine + ')'],
        ['Used Cargo Resources', b.cargo_resources + ' (' + statistics.technologies.cargo + ')'],
        ['Resources Used', statistics.used_resources],
        ['Resources Used for technologies', statistics.used_resources_for_technologies],
        ['Resources Used for ships', statistics.used_resources_for_ships],
        ['Max resources', statistics.max_resources],
        ['Resources Remaining', statistics.remaining_resources],
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
