document.body.onload = function () {
	var iframeEl = document.getElementById('source-code')
	
	iframeEl.onload = function () {
		console.log('lol');
		sourceHeight = iframeEl.contentWindow.document.body.offsetHeight + "px";
		iframeEl.style.height = sourceHeight;
	};
	
};
