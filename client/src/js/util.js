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

exports.wrap = function(lines, width, charWidth) {
	var wrapped = [];
	var lineWidth;
	var wordCount;

	for (var j = 0, llen = lines.length; j < llen; j++) {
		var words = lines[j].split(' ');
		var line = '';
		lineWidth = 0;
		wordCount = 0;

		for (var i = 0, wlen = words.length; i < wlen; i++) {
			var word = words[i];

			lineWidth += word.length * charWidth;
			wordCount++;

			if (lineWidth >= width) {
				if (wordCount !== 1) {
					wrapped.push(line);

					line = word;
					lineWidth = word.length * charWidth;
					wordCount = 1;
					
					if (i !== wlen - 1) {
						line += ' ';
						lineWidth += charWidth;
					}
				} else {
					wrapped.push(word);
					lineWidth = 0;
					wordCount = 0;
				}
			} else if (i !== wlen - 1) {
				line += word + ' ';
				lineWidth += charWidth;
			} else {
				line += word;
				wrapped.push(line);
			}
		}
	}

	return wrapped;
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