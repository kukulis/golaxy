import {ShipModel} from './entities/ship_model.js';

/**
 * @param {ShipModel[]} models
 * @returns {HTMLTableElement}
 */
export function createShipModelTable(models) {
    const thead = document.createElement('thead');
    const headerRow = document.createElement('tr');
    ['ID', 'Name', 'Guns', 'Gun Mass', 'Defense Mass', 'Engine Mass', 'Cargo Mass'].forEach(label => {
        const th = document.createElement('th');
        th.appendChild(document.createTextNode(label));
        headerRow.appendChild(th);
    });
    thead.appendChild(headerRow);

    const tbody = document.createElement('tbody');

    for (const m of models) {
        const tr = document.createElement('tr');

        const tdId = document.createElement('td');
        const linkId = document.createElement('a');
        linkId.appendChild(document.createTextNode(m['id']));
        linkId.setAttribute('href', '#');
        tdId.appendChild(linkId);
        tr.appendChild(tdId);

        const tdName = document.createElement('td');
        const linkName = document.createElement('a');
        linkName.appendChild(document.createTextNode(m['name']));
        linkName.setAttribute('href', '#');
        tdName.appendChild( linkName);
        tr.appendChild(tdName);

        ['guns', 'one_gun_mass', 'defense_mass', 'engine_mass', 'cargo_mass'].forEach(key => {
            const td = document.createElement('td');
            td.appendChild(document.createTextNode(m[key]));
            tr.appendChild(td);
        });

        tbody.appendChild(tr);
    }

    const table = document.createElement('table');
    table.appendChild(thead);
    table.appendChild(tbody);
    return table;
}
