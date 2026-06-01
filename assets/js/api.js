import { Battle } from './entities/battle.js';
import { Challenge } from './entities/challenge.js';
import { Race } from './entities/race.js';
import { Division } from './entities/division.js';
import { FleetBuild } from './entities/fleet_build.js';
import { FleetBuildShipModel } from './entities/fleet_build_ship_model.js';
import { ShipModel } from './entities/ship_model.js';
import { ShipTech } from './entities/ship_tech.js';

export class ApiClient {
    token = '';

    async _request(method, path, body) {
        const options = { method, headers: {} };

        // localStorage.getItem('token');
        let token = this.token;

        if (token) options.headers['Authorization'] = `Bearer ${token}`;
        if (body !== undefined) {
            options.headers['Content-Type'] = 'application/json';
            options.body = JSON.stringify(body);
        }
        const response = await fetch('/api' + path, options);
        if (!response.ok) {
            let message = method + ' ' + path + ' failed: ' + response.statusText;
            try {
                const data = await response.json();
                if (data.error) {
                    message = message + ': ' + data.error;
                }
            } catch (_) {}
            throw new Error(message);
        }
        return response.json();
    }

    setToken(t) {
        this.token = t
    }

    // Races

    async getCurrentRace() {
        return (new Race()).updateFromDTO(await this._request('GET', '/current-race'));
    }

    async getRaces() {
        const data = await this._request('GET', '/races');
        return data.map(d => (new Race()).updateFromDTO(d));
    }

    async getRace(id) {
        return (new Race()).updateFromDTO(await this._request('GET', `/races/${id}`));
    }

    async createRace(race) {
        return (new Race()).updateFromDTO(await this._request('POST', '/races', race));
    }

    async deleteRace(id) {
        return this._request('DELETE', `/races/${id}`);
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

    async getFleetBuilds(divisionId, all = false, raceId = '') {
        const params = new URLSearchParams()
        if (divisionId) params.set('division_id', divisionId)
        if (all) params.set('all', 'true')
        if (raceId) params.set('race_id', raceId)
        const query = params.size > 0 ? '?' + params.toString() : ''
        return this._request('GET', `/fleet-builds${query}`)
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

    async getFleetBuildShipModelsAssignments(id) {
        const data = await this._request('GET', `/fleet-builds/${id}/ship-model-assignments`);
        return data.map(d => (new FleetBuildShipModel()).updateFromDTO(d));
    }

    async getShipModelAssignment(assignmentId) {
        return (new FleetBuildShipModel()).updateFromDTO(await this._request('GET', `/ship-model-assignment/${assignmentId}`));
    }

    async addShipModelAssignment(assignment) {
        return this._request('POST', `/ship-model-assignment`, assignment);
    }

    async updateShipModelAssignment(assignmentId, assignment) {
        return this._request('POST', `/ship-model-assignment/${assignmentId}`, assignment);
    }

    async unassignShipModel(assignmentId) {
        return this._request('DELETE', `/ship-models-assignment/${assignmentId}`);
    }

    async getFleetBuildStatistics(fleetBuildId) {
        return this._request('GET', `/fleet-builds/${fleetBuildId}/statistics`);
    }

    async getFleetBuildTechnologies(fleetBuildId) {
        return this._request('GET', `/fleet-builds/${fleetBuildId}/technologies`);
    }

    async calculateShipTech(fleetBuildId, shipModelId) {
        return this._request('GET', `/fleet-builds/${fleetBuildId}/ship-models/${shipModelId}/calculate-ship-tech`);
    }

    async calculateAssignmentsShipTechs(fleetBuildId) {
        const data = await this._request('GET', `/fleet-builds/${fleetBuildId}/calculate-assignments-ship-techs`);
        return data.map(d => ({ assignment_id: d.assignment_id, ship_tech: (new ShipTech()).updateFromDTO(d.ship_tech) }));
    }

    // Challenges

    async getChallenges() {
        const data = await this._request('GET', '/challenges');
        return data.map(d => (new Challenge()).updateFromDTO(d));
    }

    async getChallenge(id) {
        return (new Challenge()).updateFromDTO(await this._request('GET', `/challenges/${id}`));
    }

    async createChallenge(challenge) {
        return (new Challenge()).updateFromDTO(await this._request('POST', '/challenges', challenge));
    }

    async updateChallenge(id, challenge) {
        return (new Challenge()).updateFromDTO(await this._request('PUT', `/challenges/${id}`, challenge));
    }

    async deleteChallenge(id) {
        return this._request('DELETE', `/challenges/${id}`);
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
