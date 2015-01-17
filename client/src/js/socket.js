var EventEmitter = require('events').EventEmitter;
var _ = require('lodash');

var sockets = {};

function createSocket(path) {
	if (sockets[path]) {
		return sockets[path];
	} else {
		var ws = new WebSocket('ws://' + window.location.host + path);
	
		var sock = {
			send: function(type, data) {
				ws.send(JSON.stringify({ type: type, request: data }));
			}
		};

		_.extend(sock, EventEmitter.prototype);

		sockets[path] = sock;

		ws.onopen = function() {
			sock.emit('connect');
		};

		ws.onclose = function() {
			sock.emit('disconnect');
		};

		ws.onmessage = function(e) {
			var msg = JSON.parse(e.data);

			sock.emit(msg.type, msg.response);
		};

		return sock;
	}
}

module.exports = createSocket;