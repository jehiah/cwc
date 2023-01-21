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
        if (request.Address != undefined) {
            if (document.getElementById("n311_portalcustomeraddressline1") !== null) {
                // TODO
                // document.getElementById("contactEmailAddress").value = request.Address.Email;
                // document.getElementById("contactFirstName").value = request.Address.FirstName;
                // document.getElementById("contactLastName").value = request.Address.LastName;
                // document.getElementById("contactDaytimePhone").value = request.Address.PhoneNumber;
                document.getElementById("n311_portalcustomeraddressline1").value = request.Address.AddressLine1;
                document.getElementById("n311_portalcustomeraddressline2").value = request.Address.AddressLine2;
                document.getElementById("n311_portalcustomeraddressborough").selectedIndex = 3; // manhattan // TODO
                document.getElementById("n311_portalcustomeraddresscity").value = request.Address.City;
                document.getElementById("n311_portalcustomeraddressstate").value = request.Address.State;
                document.getElementById("n311_portalcustomeraddresszip").value = request.Address.ZipCode;
           }
            return true
        }
        
        if (document.getElementById("n311_attendhearing_1") !== null) {
            document.getElementById("n311_attendhearing_1").click();
            if (request.Complaint.vehicle_type == "FHV") {
                document.getElementById("n311_coloroftaxi_2").click() // aka other
                document.getElementById("n311_licensenumber").value = request.Complaint.license_plate;
            } else {
                document.getElementById("n311_coloroftaxi_1").click() // yellow
                document.getElementById("n311_taximedallionnumber_name").value = request.Complaint.license_plate;
                // TODO: delay
                // TODO: set n311_problemdetailid_select first?
                document.getElementById("n311_additionaldetailsid_select").options[3].selected = true // Unsafe Driving - Non-Passenger
            }
            document.getElementById("n311_description").value = request.Complaint.description;
            document.getElementById("n311_datetimeobserved_datepicker_description").value = request.DateTimeOfIncident;
        } else if (document.getElementById("n311_locationtypeid_select") !== null) {
            document.getElementById("n311_locationtypeid_select").options[1].selected = true // street
            document.getElementById("SelectAddressWhere").click()
            document.getElementById("address-search-box-input").value = request.Street + " at " + request.CrossStreet
            document.getElementById("n311_additionallocationdetails").value = request.Complaint.location
        }
        // sendResponse({status: "goodbye"});
        return true
  });