chrome.extension.sendMessage({}, function(response) {
	var readyStateCheckInterval = setInterval(function() {
	if (document.readyState === "complete") {
		clearInterval(readyStateCheckInterval);

		// ----------------------------------------------------------
		// This part of the script triggers when page is done loading
		console.log("Hello. This message was sent from scripts/inject.js");
		// ----------------------------------------------------------

	}
	}, 10);
});

chrome.runtime.onMessage.addListener(
    function(request, sender, sendResponse) {

        console.log("in inject.js onMessage ")
        
        if (document.getElementById("char1_1___No") !== null) { // Taxi
            document.getElementById("char1_1___No").click();
            document.getElementById("vehicletype___Yellow").click();
            document.getElementById("affidavit___Yes").click();
            document.getElementById("attendHearing___Yes").click();
        } else if (document.getElementById("char1_1___Withinthe5BoroughsofNewYorkCity") !== null) { // FHV
            document.getElementById("char1_1___Withinthe5BoroughsofNewYorkCity").click();
            document.getElementById("licenseType_1-B3X-3").click();
            document.getElementById("affidavit___Yes").click();
            document.getElementById("attendHearing___Yes").click();
        } else if (document.getElementById("vehicletype_1-B3X-7") !== null) { // FHV
            document.getElementById("taxiLicenseNumber").value = request.Complaint.license_plate;
            document.getElementById("vehicletype_1-B3X-7").click();
            document.getElementById("complaintDetails").value = request.Complaint.description;
            document.getElementById("dateTimeOfIncident").value = request.DateTimeOfIncident;
        } else if (document.getElementById("taxiMedallioNum") !== null) {
            document.getElementById("taxiMedallioNum").value = request.Complaint.license_plate;
            document.getElementById("complaintDetails").value = request.Complaint.description;
            document.getElementById("dateTimeOfIncident").value = request.DateTimeOfIncident;
        } else if (document.getElementById("addressType___Intersection") !== null) {
            document.getElementById("addressType___Intersection").click()
            document.getElementById("incidentBorough5").options[3].selected = true  // manhattan // TODO
            document.getElementById("incidentOnStreet").value = request.Street
            document.getElementById("incidentStreet1Name").value = request.CrossStreet
            document.getElementById("locationDetails").value = request.Complaint.location
        } else if (document.getElementById("contactEmailAddress") !== null) {
            document.getElementById("contactEmailAddress").value = request.Address.Email;
            document.getElementById("contactFirstName").value = request.Address.FirstName;
            document.getElementById("contactLastName").value = request.Address.LastName;
            document.getElementById("contactDaytimePhone").value = request.Address.PhoneNumber;
            document.getElementById("contactBorough").selectedIndex = 3; // manhattan // TODO
            document.getElementById("contactAddressNumber").value = request.Address.StreetNumber;
            document.getElementById("contactStreetName").value = request.Address.StreetName;
            document.getElementById("contactApartment").value = request.Address.Apartment;
        }
        // sendResponse({status: "goodbye"});
        return true
  });