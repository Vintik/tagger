document.body.onload = function () {
	var iframeEl = document.getElementById('source-code');
	
	iframeEl.onload = function () {
		var counts = iframeEl.contentWindow.counts || {};
		var countsEl = document.getElementsByClassName('counts')[0];
		console.log(counts);

		sourceHeight = iframeEl.contentWindow.document.body.offsetHeight + "px";
		iframeEl.style.height = sourceHeight;

		for (tagName in counts) {
			(function () {
				var count = counts[tagName],
				    tr = createRow(tagName, count);
				
				countsEl.appendChild(tr);
			})();
		}
		
		
	};
	
};

function highlight(tagName) {
	var iframeEl = document.getElementById('source-code'),
	    prevEls = iframeEl.contentDocument.getElementsByClassName('highlight'),
	    tagEls = iframeEl.contentDocument.getElementsByClassName(tagName),
	    el;

	// Clear the previously highlighted tagNames
	while (prevEls.length) {
		el = prevEls[i];
		el.className = el.className.replace('highlight', '');
	}

	// Highlight the current tagNames
	for (var i = 0; i < tagEls.length; i++) {
		tagEls[i].className += " highlight";
	}
};

function createTagNameCell(tagName) {
	var td = document.createElement('td'),
	    anchor = document.createElement('a');

	anchor.innerText = tagName;
	anchor.href = "#";
	anchor.onclick = function (e) {
		e.preventDefault();
		highlight(tagName);
	};

	td.appendChild(anchor);

	return td;
}

function createTagCountCell(count) {
	var td = document.createElement('td');

	td.innerText = count;
	return td;
}

function createRow(tagName, count) {
	var tr = document.createElement('tr'),
	    tagNameEl = createTagNameCell(tagName),
	    countEl = createTagCountCell(count);

	tr.appendChild(tagNameEl);
	tr.appendChild(countEl);

	return tr;
}
