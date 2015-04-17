document.body.onload = function () {
	var iframeEl = document.getElementById('source-code');
	
	iframeEl.onload = function () {
		var counts = iframeEl.contentWindow.counts || {};
		var countsEl = document.getElementsByClassName('counts')[0];

		countsEl.innerHTML = '';

		sourceHeight = iframeEl.contentWindow.document.body.scrollHeight + "px";
		iframeEl.style.height = sourceHeight;

		for (tagName in counts) {
			(function () {
				var count = counts[tagName],
				    li = createListItem(tagName, count);
				
				countsEl.appendChild(li);
			})();
		}
	};

	document.onresize = function (e) {
		sourceHeight = iframeEl.contentWindow.document.body.scrollHeight + "px";
		iframeEl.style.height = sourceHeight;
	};
	
};

function highlight(tagName) {
	var iframeEl = document.getElementById('source-code'),
	    prevEls = iframeEl.contentDocument.getElementsByClassName('highlight'),
	    tagEls = iframeEl.contentDocument.getElementsByClassName(tagName),
	    el;

	// Clear the previously highlighted tagNames
	//
	// A for loop was avoided here due to the `prevEls` actually being an HTMLCollection
	// HTMLCollection will automatically remove items if they no longer fit the className selected by
	while (prevEls.length) {
		el = prevEls[prevEls.length - 1];
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

	td.appendChild(anchor);

	return td;
}

function createTagCountCell(count) {
	var td = document.createElement('td');

	td.innerText = count;
	return td;
}

function createListItem(tagName, count) {
	var li = document.createElement('li'),
	    btnEl = document.createElement('a'),
	    tagNameEl = document.createElement('span'),
	    countEl = document.createElement('span');

	btnEl.className = "button";
	btnEl.href = '#'

	btnEl.onclick = function (e) {
		e.preventDefault();
		highlight(tagName);
	};

	tagNameEl.innerText = tagName;
	tagNameEl.className = 'list-tagname'

	countEl.innerText = count;
	countEl.className = 'list-tagcount'

	btnEl.appendChild(countEl);
	btnEl.appendChild(tagNameEl);

	li.appendChild(btnEl);
	return li;
}
