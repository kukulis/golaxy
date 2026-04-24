import {FleetBuildShipModel} from './entities/fleet_build_ship_model.js';

/**
 * @param {FleetBuildShipModel[]} assignments
 * @returns {HTMLTableElement}
 */
export function createFleetBuildShipModelTable(assignments) {
    const thead = document.createElement('thead');
    const headerRow = document.createElement('tr');
    ['Name', 'Guns', 'Gun Mass', 'Defense Mass', 'Engine Mass', 'Cargo Mass', 'Amount', 'Result Mass'].forEach(label => {
        const th = document.createElement('th');
        th.appendChild(document.createTextNode(label));
        headerRow.appendChild(th);
    });
    thead.appendChild(headerRow);

    const tbody = document.createElement('tbody');
    for (const a of assignments) {
        const tr = document.createElement('tr');
        [
            a.shipModel.name,
            a.shipModel.guns,
            a.shipModel.one_gun_mass,
            a.shipModel.defense_mass,
            a.shipModel.engine_mass,
            a.shipModel.cargo_mass,
            a.amount,
            a.result_mass,
        ].forEach(value => {
            const td = document.createElement('td');
            td.appendChild(document.createTextNode(value));
            tr.appendChild(td);
        });
        tbody.appendChild(tr);
    }

    const table = document.createElement('table');
    table.appendChild(thead);
    table.appendChild(tbody);
    return table;
}
