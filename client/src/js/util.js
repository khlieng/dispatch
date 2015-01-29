var _ = require('lodash');

exports.UUID = function() {
	return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
	    var r = Math.random()*16|0, v = c == 'x' ? r : (r&0x3|0x8);
	    return v.toString(16);
	});
};

exports.timestamp = function(date) {
	date = date || new Date();
	
	var h = _.padLeft(date.getHours(), 2, '0')
	var m = _.padLeft(date.getMinutes(), 2, '0');

	return h + ':' + m;
};