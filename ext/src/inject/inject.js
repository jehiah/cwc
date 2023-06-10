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

function setIsDirty(target) {
    var el = document.getElementById(target)
    el.classList.add("dirty")
}
function clickElement(target) {
    const clickEvent = new MouseEvent('click', {
        bubbles: true,
        cancelable: true,
        view: window
      });
      target.dispatchEvent(clickEvent)
}

chrome.runtime.onMessage.addListener(
    function(request, sender, sendResponse) {

        console.log("in inject.js onMessage ", request)

        // on start page
        if (document.querySelector("h1.entry-title").innerText == "Taxi Complaint") {
            let a = document.querySelectorAll('a')
            for (var i = 0; i < a.length; i++) {
                if (a[i].innerText == "Report a problem with a driver, if you were NOT a passenger.") {
                    clickElement(a[i])
                }
            }
        }


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
                // n311_additionaldetailsid == eb4e791a-374e-e811-a94d-000d3a360e00
            }
            document.getElementById("n311_description").value = request.Complaint.description;
            // '2023-05-26T21:08:00.0000000Z'
            document.getElementById("n311_datetimeobserved").value = request.DateTimeOfIncidentISO
            document.getElementById("n311_datetimeobserved_datepicker_description").value = request.DateTimeOfIncident;
            setIsDirty("n311_datetimeobserved_datepicker_description");

            console.log("request.Complaint.videos.length", request.Complaint.videos, request.Complaint.videos.length)
            // attachments-addbutton
            if (request.Complaint.videos != null  || request.Complaint.photos != null) {
                clickElement(document.getElementById("attachments-addbutton"))
                clickElement(document.getElementsByName("file")[0])
            }

        } else if (document.getElementById("n311_locationtypeid_select") !== null) {
            document.getElementById("n311_locationtypeid_select").options[1].selected = true // street
            document.getElementById("n311_additionallocationdetails").value = request.Complaint.location
            document.getElementById("SelectAddressWhere").click()
            document.getElementById("address-search-box-input").value = request.Street + " at " + request.CrossStreet
        }
        // sendResponse({status: "goodbye"});
        return true
  });