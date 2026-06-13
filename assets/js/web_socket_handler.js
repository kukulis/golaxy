import {Dispatcher} from "./dispatcher.js";

export class WebSocketHandler {
    /** @type {WebSocket|null} */
    conn = null;
    /** @type {Dispatcher|null} */
    dispatcher = null;

    /**
     * @param {Dispatcher} dispatcher
     */
    constructor(dispatcher) {
        this.dispatcher = dispatcher;
    }

    /** @returns {void} */
    connect() {
        if (!window["WebSocket"]) {
            this.dispatcher.dispatch("ws:unsupported");
            return;
        }
        this.conn = new WebSocket("ws://" + document.location.host + "/ws");
        this.conn.onclose = (evt) => {
            this.dispatcher.dispatch("ws:close", evt);
        };
        this.conn.onmessage = (evt) => {
            var messages = evt.data.split('\n');
            for (var i = 0; i < messages.length; i++) {
                this.dispatcher.dispatch("ws:message", messages[i]);
            }
        };
    }

    /**
     * @param {string} msg
     * @returns {void}
     */
    send(msg) {
        if (this.conn) {
            this.conn.send(msg);
        }
    }
}
