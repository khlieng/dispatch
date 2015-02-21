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

exports.wrapMessages = function(messages, width, charWidth, indent = 0) {
	for (var j = 0, llen = messages.length; j < llen; j++) {
		var message = messages[j];
		var words = message.message.split(' ');
		var line = '';
		var wrapped = [];
		var lineWidth = (6 + (message.from ? message.from.length + 1 : 0)) * charWidth;
		var wordCount = 0;
		var hasWrapped = false;

		// Add empty line if first word after timestamp + sender wraps
		if (words.length > 0 && message.from && lineWidth + words[0].length * charWidth >= width) {
			wrapped.push(line);
			lineWidth = 0;
		}

		for (var i = 0, wlen = words.length; i < wlen; i++) {
			var word = words[i];

			if (hasWrapped) {
				hasWrapped = false;
				lineWidth += indent;
			}

			lineWidth += word.length * charWidth;
			wordCount++;

			if (lineWidth >= width) {
				if (wordCount !== 1) {
					wrapped.push(line);

					if (i !== wlen - 1) {
						line = word + ' ';
						lineWidth = (word.length + 1) * charWidth;
						wordCount = 1;
					} else {
						wrapped.push(word);
					}
				} else {
					wrapped.push(word);
					lineWidth = 0;
					wordCount = 0;
				}

				hasWrapped = true;
			} else if (i !== wlen - 1) {
				line += word + ' ';
				lineWidth += charWidth;
			} else {
				line += word;
				wrapped.push(line);
			}
		}

		message.lines = wrapped;
	}
};

var canvas = document.createElement('canvas');
var ctx = canvas.getContext('2d');

exports.stringWidth = function(str, font) {	
	ctx.font = font;
	return ctx.measureText(str).width;
};

exports.scrollbarWidth = function() {
    var outer = document.createElement('div');
    outer.style.visibility = 'hidden';
    outer.style.width = '100px';

    document.body.appendChild(outer);

    var widthNoScroll = outer.offsetWidth;
    outer.style.overflow = 'scroll';

    var inner = document.createElement('div');
    inner.style.width = '100%';
    outer.appendChild(inner);        

    var widthWithScroll = inner.offsetWidth;

    outer.parentNode.removeChild(outer);

    return widthNoScroll - widthWithScroll;
};