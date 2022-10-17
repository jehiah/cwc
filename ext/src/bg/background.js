// if you checked "fancy-settings" in extensionizr.com, uncomment this lines

// var settings = new Store("settings", {
//     "sample_setting": "This is how you use Store.js to remember values"
// });


// When the extension is installed or upgraded ...
chrome.runtime.onInstalled.addListener(function() {
  // Replace all rules ...
  chrome.declarativeContent.onPageChanged.removeRules(undefined, function() {
    // With a new rule ...
    chrome.declarativeContent.onPageChanged.addRules([
      {
        // That fires when a page's URL contains a 'g' ...
        conditions: [
          new chrome.declarativeContent.PageStateMatcher({
            pageUrl: { urlContains: 'https://portal.311.nyc.gov/sr-step/' },
          })
        ],
        // And shows the extension's page action.
        actions: [ new chrome.declarativeContent.ShowPageAction() ]
      }
    ]);
  });
});

// from https://mathiasbynens.be/notes/xhr-responsetype-json
var getJSON = function(url) {
	return new Promise(function(resolve, reject) {
		var xhr = new XMLHttpRequest();
		xhr.open('get', url, true);
		xhr.responseType = 'json';
		xhr.onload = function() {
			var status = xhr.status;
			if (status == 200) {
				resolve(xhr.response);
			} else {
				reject(status);
			}
		};
		xhr.send();
	});
};

//example of using a message handler from the inject scripts
// chrome.extension.onMessage.addListener(
//   function(request, sender, sendResponse) {
//       chrome.pageAction.show(sender.tab.id);
//     sendResponse();
//   });

chrome.pageAction.onClicked.addListener(function(tab){
    console.log("pageAction.onClicked")
    getJSON("http://[::]:53000/complaint/latest.json").then(function(data) {
        console.log("got latest.json", data)
        chrome.tabs.sendMessage(tab.id, data, function(response) {
          console.log(response);
        });
    }, function(status) {
    	alert('Something went wrong.');
    });
})
