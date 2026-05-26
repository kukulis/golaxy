import {Dispatcher} from "./dispatcher.js";
import {ClearE, GetE, NewT} from "./helper.js";
import DivisionsView from "./divisions_view.js";
import FleetBuildsView from "./fleet_builds_view.js";
import RacesView from "./races_view.js";
import {ApiClient} from "./api.js";

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
     * @type {RacesView}
     */
    racesView = null;

    /**
     *
     * @type {ApiClient}
     */
    apiClient = null;

    /**
     * @returns {Dispatcher}
     */
    getDispatcher() {
        if ( this.dispatcher == null ) {
            this.dispatcher = new Dispatcher()

            // error handling
            this.dispatcher.addListener("displayError", ([msg, clear]) => {
                console.error( "Error received", msg)

                let errorMsg = GetE("error-msg")
                if ( clear) {
                    ClearE(errorMsg)
                }

                errorMsg.appendChild(NewT(msg))
                errorMsg.style.display = 'block'
            })
        }

        return this.dispatcher
    }


    /**
     * @returns {DivisionsView}
     */
    getDivisionsView() {
        if ( this.divisionsView == null ) {
            this.divisionsView = new DivisionsView(this.getApiClient(), this.getDispatcher())
        }

        return this.divisionsView
    }

    /**
     * @returns {FleetBuildsView}
     */
    getFleetBuildsView() {
        if ( this.fleetBuildsView == null ) {
            this.fleetBuildsView = new FleetBuildsView(this.getApiClient(), this.getDispatcher())
        }

        return this.fleetBuildsView
    }

    /**
     * @returns {RacesView}
     */
    getRacesView() {
        if ( this.racesView == null ) {
            this.racesView = new RacesView(this.getApiClient(), this.getDispatcher())
        }

        return this.racesView
    }

    /**
     * @returns {ApiClient}
     */
    getApiClient() {
        if ( this.apiClient == null ) {
            this.apiClient = new ApiClient()
        }

        return this.apiClient
    }

}
