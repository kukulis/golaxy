export class WSMessagesHandler {
    constructor() {
        this.conn = null;
        this.msg = document.getElementById("msg");
        this.log = document.getElementById("log");
        this.wrapper = document.getElementById("wrapper");
    }

    init() {
        this._initDraggable();
        this._initForm();
        this._connectWs();
    }

    _appendLog(item) {
        var doScroll = this.log.scrollTop > this.log.scrollHeight - this.log.clientHeight - 1;
        this.log.appendChild(item);
        if (doScroll) {
            this.log.scrollTop = this.log.scrollHeight - this.log.clientHeight;
        }
    }

    _initForm() {
        document.getElementById("form").onsubmit = () => {
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

    _initDraggable() {
        let dragHandle = this.wrapper.querySelector('.task-details-drag-handle');
        if (!dragHandle) return;

        let isDragging = false;
        let initialX, initialY;

        dragHandle.addEventListener('mousedown', (e) => {
            isDragging = true;
            initialX = e.clientX - this.wrapper.offsetLeft;
            initialY = e.clientY - this.wrapper.offsetTop;
            dragHandle.style.cursor = 'grabbing';
        });

        document.addEventListener('mousemove', (e) => {
            if (!isDragging) return;
            e.preventDefault();
            this.wrapper.style.left = (e.clientX - initialX) + 'px';
            this.wrapper.style.top = (e.clientY - initialY) + 'px';
        });

        document.addEventListener('mouseup', () => {
            isDragging = false;
            dragHandle.style.cursor = 'move';
        });
    }
}
