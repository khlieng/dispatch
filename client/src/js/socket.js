var EventEmitter = require('events').EventEmitter;

class Socket extends EventEmitter {
	constructor() {
		super();
		
		this.ws = new WebSocket('ws://' + window.location.host + '/ws');

		this.ws.onopen = () => this.emit('connect');
		this.ws.onclose = () => this.emit('disconnect');
		this.ws.onmessage = (e) => {
			var msg = JSON.parse(e.data);

			this.emit(msg.type, msg.response);
		}
	}

	send(type, data) {
		this.ws.send(JSON.stringify({ type, request: data }));
	}
}

module.exports = new Socket();