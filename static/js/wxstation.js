$(function() {
	function updateStatus() {
		$.getJSON('/status', '', function(data, textStatus, jqXHR) {
			$("#curTemp").html(data['temp'] + '&deg;C');
			$("#curHumidity").html(data['humidity'] + '%');
			if (data['door_status'] == 1) {
				$("#doorStatus").html("Open");
			} else {
				$("#doorStatus").html("Closed");				
			}
		});
	}
	updateStatus();
	// every 30 seconds
	setInterval(updateStatus, 30000);

	function changeGraph(event) {
		$("#gGroup button").removeClass("active");
		$(event.target).addClass("active");
		$("#tempGraph").html('<img src="/g/temp/'+event.data+'">');
		$("#humidityGraph").html('<img src="/g/humidity/'+event.data+'">');
	}
	$("#gDay").click("day", changeGraph);
	$("#gWeek").click("week", changeGraph);
	$("#gMonth").click("month", changeGraph);
	$("#gYear").click("year", changeGraph);
});