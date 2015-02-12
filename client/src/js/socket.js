var EventEmitter = require('events').EventEmitter;

var _ = require('lodash');

var ws = new WebSocket('ws://' + window.location.host + '/ws');

var socket = {
	send: function(type, data) {
		ws.send(JSON.stringify({ type: type, request: data }));
	}
};

_.extend(socket, EventEmitter.prototype);

ws.onopen = function() {
	socket.emit('connect');
};

ws.onclose = function() {
	socket.emit('disconnect');
};

ws.onmessage = function(e) {
	var msg = JSON.parse(e.data);

	socket.emit(msg.type, msg.response);
};

module.exports = socket;