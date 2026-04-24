import { Battle } from './entities/battle.js';

export class ApiClient {

    async _request(method, path, body) {
        const options = { method, headers: {} };
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
        return this._request('GET', '/divisions');
    }

    async getDivision(id) {
        return this._request('GET', `/divisions/${id}`);
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

    async getFleetBuilds() {
        return this._request('GET', '/fleet-builds');
    }

    async getFleetBuild(id) {
        return this._request('GET', `/fleet-builds/${id}`);
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
        return this._request('GET', '/ship-models');
    }

    async getShipModel(id) {
        return this._request('GET', `/ship-models/${id}`);
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
