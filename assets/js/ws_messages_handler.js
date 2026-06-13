export class WSMessagesHandler {
    /** @type {WebSocket|null} */
    conn = null;
    /** @type {HTMLInputElement|null} */
    msg = null;
    /** @type {HTMLElement|null} */
    log = null;
    /** @type {HTMLElement|null} */
    wrapper = null;
    /** @type {HTMLFormElement|null} */
    form = null;
    /** @type {HTMLElement|null} */
    dragHandle = null;

    /**
     * @param {HTMLElement} wrapper
     */
    constructor(wrapper) {
        this.wrapper = wrapper;
    }

    /** @returns {void} */
    init() {
        this._buildHtml();
        this._initDraggable();
        this._initForm();
        this._connectWs();
    }

    /** @returns {void} */
    _buildHtml() {
        this.dragHandle = document.createElement("div");
        this.dragHandle.className = "task-details-drag-handle";
        this.dragHandle.textContent = "Drag me";

        this.log = document.createElement("div");

        this.form = document.createElement("form");
        const form = this.form;
        const submit = document.createElement("input");
        submit.type = "submit";
        submit.value = "Send";

        this.msg = document.createElement("input");
        this.msg.type = "text";
        this.msg.size = 64;
        this.msg.autofocus = true;

        form.appendChild(submit);
        form.appendChild(this.msg);

        this.wrapper.appendChild(this.dragHandle);
        this.wrapper.appendChild(this.log);
        this.wrapper.appendChild(form);
    }

    /**
     * @param {HTMLElement} item
     * @returns {void}
     */
    _appendLog(item) {
        var doScroll = this.log.scrollTop > this.log.scrollHeight - this.log.clientHeight - 1;
        this.log.appendChild(item);
        if (doScroll) {
            this.log.scrollTop = this.log.scrollHeight - this.log.clientHeight;
        }
    }

    /** @returns {void} */
    _initForm() {
        this.form.onsubmit = () => {
            if (!this.conn) {
                return false;
            }
            if (!this.msg.value) {
                return false;
            }
            this.conn.send(this.msg.value);
            this.msg.value = "";
            return false;
        };
    }

    /** @returns {void} */
    _connectWs() {
        if (!window["WebSocket"]) {
            var item = document.createElement("div");
            item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
            this._appendLog(item);
            return;
        }
        this.conn = new WebSocket("ws://" + document.location.host + "/ws");
        this.conn.onclose = (evt) => {
            var item = document.createElement("div");
            item.innerHTML = "<b>Connection closed.</b>";
            this._appendLog(item);
        };
        this.conn.onmessage = (evt) => {
            var messages = evt.data.split('\n');
            for (var i = 0; i < messages.length; i++) {
                var item = document.createElement("div");
                item.innerText = messages[i];
                this._appendLog(item);
            }
        };
    }

    /** @returns {void} */
    _initDraggable() {
        if (!this.dragHandle) return;

        let isDragging = false;
        let initialX, initialY;

        this.dragHandle.addEventListener('mousedown', (e) => {
            isDragging = true;
            initialX = e.clientX - this.wrapper.offsetLeft;
            initialY = e.clientY - this.wrapper.offsetTop;
            this.dragHandle.style.cursor = 'grabbing';
        });

        // TODO research if the events listeners must be declared globally
        document.addEventListener('mousemove', (e) => {
            if (!isDragging) return;
            e.preventDefault();
            this.wrapper.style.left = (e.clientX - initialX) + 'px';
            this.wrapper.style.top = (e.clientY - initialY) + 'px';
        });

        document.addEventListener('mouseup', () => {
            isDragging = false;
            this.dragHandle.style.cursor = 'move';
        });
    }
}
