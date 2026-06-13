export class WSMessagesHandler {
    /** @type {WebSocket|null} */
    conn = null;
    /** @type {HTMLInputElement|null} */
    msg = null;
    /** @type {HTMLElement|null} */
    log = null;
    /** @type {HTMLElement|null} */
    container = null;
    /** @type {HTMLElement|null} */
    wrapper = null;
    /** @type {HTMLFormElement|null} */
    form = null;
    /** @type {HTMLElement|null} */
    dragHandle = null;
    /** @type {HTMLButtonElement|null} */
    showChatBtn = null;

    constructor() {
    }

    /**
     * @param {HTMLElement} container
     * @returns {void}
     */
    init(container) {
        this.container = container;
        this._buildHtml();
        this._initDraggable();
        this._initForm();
        this._connectWs();
    }

    /** @returns {void} */
    _buildHtml() {
        this.wrapper = document.createElement("div");
        this.wrapper.className = "chat";

        this.showChatBtn = document.createElement("button");
        this.showChatBtn.className = "chat-show-btn";
        this.showChatBtn.textContent = "Chat";

        const hideBtn = document.createElement("button");
        hideBtn.textContent = "_";
        hideBtn.addEventListener("click", () => { this.wrapper.style.display = "none"; });

        this.dragHandle = document.createElement("div");
        this.dragHandle.className = "chat-drag-handle";
        this.dragHandle.textContent = "Global chat";
        this.dragHandle.appendChild(hideBtn);

        this.log = document.createElement("div");
        this.log.className = "chat-log";

        this.form = document.createElement("form");
        const submit = document.createElement("input");
        submit.type = "submit";
        submit.value = "Send";

        this.msg = document.createElement("input");
        this.msg.type = "text";
        this.msg.size = 64;
        this.msg.autofocus = true;

        this.form.appendChild(submit);
        this.form.appendChild(this.msg);

        const formWrapper = document.createElement("div");
        formWrapper.className = "chat-form";
        formWrapper.appendChild(this.form);

        this.wrapper.appendChild(this.dragHandle);
        this.wrapper.appendChild(this.log);
        this.wrapper.appendChild(formWrapper);

        this.showChatBtn.addEventListener("click", () => { this.wrapper.style.display = "flex"; });

        this.container.appendChild(this.wrapper);
        this.container.appendChild(this.showChatBtn);
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
