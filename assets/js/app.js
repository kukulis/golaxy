import {Dispatcher} from "./dispatcher";
import {GetE, NewT} from "./helper";

export class App {


    /**
     *
     * @type {Dispatcher}
     */
    dispatcher = null;

    /**
     * @returns {Dispatcher}
     */
    getDispatcher() {
        if ( this.dispatcher == null ) {
            this.dispatcher = new Dispatcher()

            // error handling

            this.dispatcher.addListener("displayError", (msg, clear) => {
                let errorMsg = GetE("error-msg")
                if ( clear) {
                    errorMsg.appendChild(NewT(msg))
                }
            })
        }

        return this.dispatcher
    }



}
