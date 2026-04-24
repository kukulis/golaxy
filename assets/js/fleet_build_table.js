import {FleetBuild} from './entities/fleet_build.js';

/**
 * @param {FleetBuild} b
 * @returns {HTMLTableElement}
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
