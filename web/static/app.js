function get_data(url, callback, errorCallback) {
    var req = new XMLHttpRequest();
    req.overrideMimeType("application/json");
    req.open('GET', url, true);
    req.onreadystatechange = function() {
	if (req.readyState === XMLHttpRequest.DONE) {
	    if (req.status === 200) {
		var jsonResponse = JSON.parse(req.responseText);
		callback(jsonResponse);
		return;
	    } else {
		if (errorCallback !== undefined) {
		    errorCallback();
		}
	    }
	}
    };
    req.send();
    return req;
}

function fetchInfo(id, counter) {
    if (counter < 0) {
	return;
    }
    get_data('/api/info/' + id, function(response) {
	if (!('status' in response) || !('overview' in response)) {
	    window.setTimeout(fetchInfo, 1000, id, counter - 1);
	    return;
	}

	let config = [
	    {'id': 'first-aired-span', 'field': 'firstAired' },
	    {'id': 'status-span', 'field': 'status' },
	    {'id': 'network-span', 'field': 'network' },
	    {'id': 'genre-span', 'field': 'genre', 'func': function(x) { return x.join(', ');} },
	    {'id': 'overview-span', 'field': 'overview' },
	];

	config.forEach(function( e ) {
	    let elem = document.getElementById(e.id);
	    let txt = response[e.field];
	    if (e.func) {
		txt = e.func(txt);
	    }
	    
	    elem.innerText = txt;
	});
    });
}

function fetchImage(id, counter) {
    if (counter < 0) {
	return;
    }
    get_data('/api/image/' + id, function(response) {
	if (!('url' in response)) {
	    window.setTimeout(fetchImage, 1000, id, counter - 1);
	    return;
	}
	let img = document.getElementById('showImage');
	img.src = response['url'];
    });
};

function trackOutboundLink(url) {
    gtag('event', 'click', {
	'event_category': 'outbound',
	'event_label': url,
	'transport_type': 'beacon',
	'event_callback': function() {document.location = url;}
    });
};



(function() {
    var xhr;

    let searchBox = document.getElementById('search-results-box');
    let searchList = document.getElementById('search-results-list');
    let nothingFound = document.getElementById('no-results');
    let topSearches = document.getElementById('top-searches-box');

    function hideTitles() {
	searchBox.style.display = 'none';
	nothingFound.style.display = 'none';
    }
    
    function displayTitles(result) {
	while( searchList.firstChild) {
	    searchList.removeChild(searchList.firstChild);
	}
	
	if (result.length == 0) {
	    searchBox.style.display = 'none';
	    nothingFound.style.display = 'block';
	    return;
	}
	
	searchBox.style.display = 'block';
	nothingFound.style.display = 'none';

	result.forEach(function(elem, idx) {
	    let li = document.createElement('li');

	    let a = document.createElement('a');
	    a.href = '/' + elem.t_id;
	    let title = elem.primary_title + ' (' + elem.start_year + ')';
	    let aTitle = document.createTextNode(title);
	    a.appendChild(aTitle);
	    li.appendChild(a);
	    
	    searchList.appendChild(li);
	});
    }
    
    function searchTitles(value) {
	if (!value) {
	    hideTitles();
	    return;
	}
	
	try {
	    xhr.abort();
	} catch(e){};
	xhr = get_data('/api/search?filter=' + value, displayTitles);
    }
    
    document.getElementById('mySearch').onkeypress =
	function(e) {
	    if (e.key === 'Enter') {
		searchTitles(e.target.value);
	    }
	};
    
    
})();
