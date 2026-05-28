import {Dispatcher} from "./dispatcher.js";
import {ClearE, GetE, NewT} from "./helper.js";
import DivisionsView from "./divisions_view.js";
import FleetBuildsView from "./fleet_builds_view.js";
import FleetBuildView from "./fleet_build_view.js";
import RacesView from "./races_view.js";
import DummyLoginView from "./dummy_login_view.js";
import UserBarView from "./user_bar_view.js";
import UserMenuView from "./user_menu_view.js";
import {ApiClient} from "./api.js";
import {ShipModelComponent} from "./ship_model_component.js";

export class App {
    /**
     *
     * @type {Dispatcher}
     */
    dispatcher = null;

    /**
     *
     * @type {DivisionsView}
     */
    divisionsView = null;

    /**
     *
     * @type {FleetBuildsView}
     */
    fleetBuildsView = null;

    /**
     *
     * @type {FleetBuildView}
     */
    fleetBuildView = null;

    /**
     *
     * @type {RacesView}
     */
    racesView = null;

    /**
     *
     * @type {DummyLoginView}
     */
    dummyLoginView = null;

    /**
     *
     * @type {ShipModelComponent}
     */
    shipModelComponent = null;

    /**
     *
     * @type {ApiClient}
     */
    apiClient = null;

    /**
     * @returns {Dispatcher}
     */
    getDispatcher() {
        if (this.dispatcher == null) {
            this.dispatcher = new Dispatcher()

            // error handling
            this.dispatcher.addListener("displayError", ([msg, clear]) => {
                console.error("Error received", msg)

                let errorMsg = GetE("error-msg")
                if (errorMsg === undefined) {
                    alert(msg)
                    return
                }
                if (clear) {
                    ClearE(errorMsg)
                }

                errorMsg.appendChild(NewT(msg))
                errorMsg.style.display = 'block'
            })

            this.dispatcher.addListener("getToken", () => this.getToken())
            this.dispatcher.addListener("getLoginName", () => this.getLoginName())
            this.dispatcher.addListener("storeLoginData", ([name, token]) => {
                this.setLoginName(name)
                this.setToken(token)
            })
            this.dispatcher.addListener("redirect", (url) => {
                window.location.href = url
            })
            this.dispatcher.addListener("afterShipModelEdit", ({shipModelId, fleetBuildId}) => {
                const query = fleetBuildId ? `?FleetBuildId=${fleetBuildId}` : ''
                location.href = `/ship-model/${shipModelId}/details.html${query}`
            })
            this.dispatcher.addListener("logout", () => {
                this.setLoginName('')
                this.setToken('')
                this.dispatcher.dispatch('redirect', '/')
            })
        }

        return this.dispatcher
    }


    /**
     * @returns {DivisionsView}
     */
    getDivisionsView() {
        if (this.divisionsView == null) {
            this.divisionsView = new DivisionsView(this.getApiClient(), this.getDispatcher())
        }

        return this.divisionsView
    }

    /**
     * @returns {FleetBuildsView}
     */
    getFleetBuildsView() {
        if (this.fleetBuildsView == null) {
            this.fleetBuildsView = new FleetBuildsView(this.getApiClient(), this.getDispatcher())
        }

        return this.fleetBuildsView
    }

    /**
     * @returns {FleetBuildView}
     */
    getFleetBuildView() {
        if (this.fleetBuildView == null) {
            this.fleetBuildView = new FleetBuildView(this.getApiClient(), this.getDispatcher())
        }

        return this.fleetBuildView
    }

    /**
     * @returns {RacesView}
     */
    getRacesView() {
        if (this.racesView == null) {
            this.racesView = new RacesView(this.getApiClient(), this.getDispatcher())
        }

        return this.racesView
    }

    renderMenu() {
        const container = document.getElementById('user-menu')
        if (container == null) {
            return
        }

        container.appendChild(new UserMenuView().generateView())
        container.appendChild(new UserBarView(this.getDispatcher()).generateView())
    }

    /**
     * @returns {DummyLoginView}
     */
    getDummyLoginView() {
        if (this.dummyLoginView == null) {
            this.dummyLoginView = new DummyLoginView(this.getApiClient(), this.getDispatcher())
        }

        return this.dummyLoginView
    }

    /**
     * @returns {string}
     */
    getLoginName() {
        return localStorage.getItem('loginName')
    }

    /**
     * @param {string} name
     */
    setLoginName(name) {
        localStorage.setItem('loginName', name)
    }

    /**
     * @returns {string}
     */
    getToken() {
        return localStorage.getItem('token')
    }

    /**
     * @param {string} token
     */
    setToken(token) {
        localStorage.setItem('token', token)
        this.apiClient.setToken(token)
    }

    /**
     * @returns {ShipModelComponent}
     */
    getShipModelComponent() {
        if (this.shipModelComponent == null) {
            this.shipModelComponent = new ShipModelComponent(this.getApiClient(), this.getDispatcher())
        }

        return this.shipModelComponent
    }

    /**
     * @returns {ApiClient}
     */
    getApiClient() {
        if (this.apiClient == null) {
            this.apiClient = new ApiClient()
        }
        this.apiClient.setToken(this.getToken())

        return this.apiClient
    }

}
