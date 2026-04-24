import { Battle } from './entities/battle.js';
import { Division } from './entities/division.js';
import { FleetBuild } from './entities/fleet_build.js';
import { ShipModel } from './entities/ship_model.js';

export class ApiClient {

    async _request(method, path, body) {
        const options = { method, headers: {} };
        const token = localStorage.getItem('token');
        if (token) options.headers['Authorization'] = `Bearer ${token}`;
        if (body !== undefined) {
            options.headers['Content-Type'] = 'application/json';
            options.body = JSON.stringify(body);
        }
        const response = await fetch('/api' + path, options);
        if (!response.ok) {
            throw new Error(`${method} ${path} failed: ${response.statusText}`);
        }
        return response.json();
    }

    // Battle

    async getBattle() {
        const data = await this._request('GET', '/battle');
        return (new Battle()).updateFromDTO(data);
    }

    // Divisions

    async getDivisions() {
        const data = await this._request('GET', '/divisions');
        return data.map(d => (new Division()).updateFromDTO(d));
    }

    async getDivision(id) {
        return (new Division()).updateFromDTO(await this._request('GET', `/divisions/${id}`));
    }

    async createDivision(division) {
        return this._request('POST', '/divisions', division);
    }

    async updateDivision(id, division) {
        return this._request('PUT', `/divisions/${id}`, division);
    }

    async deleteDivision(id) {
        return this._request('DELETE', `/divisions/${id}`);
    }

    // Fleet Builds

    async getFleetBuilds(divisionId) {
        const query = divisionId ? `?division_id=${encodeURIComponent(divisionId)}` : '';
        return this._request('GET', `/fleet-builds${query}`);
    }

    async getFleetBuild(id) {
        return (new FleetBuild()).updateFromDTO(await this._request('GET', `/fleet-builds/${id}`));
    }

    async createFleetBuild(fleetBuild) {
        return this._request('POST', '/fleet-builds', fleetBuild);
    }

    async updateFleetBuild(id, fleetBuild) {
        return this._request('PUT', `/fleet-builds/${id}`, fleetBuild);
    }

    async deleteFleetBuild(id) {
        return this._request('DELETE', `/fleet-builds/${id}`);
    }

    async getFleetBuildShipModels(id) {
        return this._request('GET', `/fleet-builds/${id}/ship-models`);
    }

    async assignShipModel(id, assignment) {
        return this._request('POST', `/fleet-builds/${id}/ship-models`, assignment);
    }

    async unassignShipModel(id, shipModelId) {
        return this._request('DELETE', `/fleet-builds/${id}/ship-models/${shipModelId}`);
    }

    // Ship Models

    async getShipModels() {
        const data = await this._request('GET', '/ship-models');
        return data.map(d => (new ShipModel()).updateFromDTO(d));
    }

    async getShipModel(id) {
        return (new ShipModel()).updateFromDTO(await this._request('GET', `/ship-models/${id}`));
    }

    async createShipModel(shipModel) {
        return this._request('POST', '/ship-models', shipModel);
    }

    async updateShipModel(id, shipModel) {
        return this._request('PUT', `/ship-models/${id}`, shipModel);
    }

    async deleteShipModel(id) {
        return this._request('DELETE', `/ship-models/${id}`);
    }
}
