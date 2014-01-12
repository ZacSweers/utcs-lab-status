var socket = io.connect("http://www.utcslabs.org");
socket.on("status", function (data) {

	var status = JSON.parse(data);

	if (status["No"] != null) {
		for (var i = 0; i < status["No"].length; i++) {
			var host = status["No"][i];
			$("#" + host).animate({backgroundColor: "#cf948c"}, "slow").attr("title", "Unavailable (" + host + ")");
		}
	}

	if (status["Yes"] != null) {
		for (var i = 0; i < status["Yes"].length; i++) {
			var host = status["Yes"][i];
			$("#" + host).animate({backgroundColor: "#b7d6a5"}, "slow").attr("title", "Available (" + host + ")");
		}
	}

	if (status["Offline"] != null) {
		for (var i = 0; i < status["Offline"].length; i++) {
			var host = status["Offline"][i];
			$("#" + host).animate({backgroundColor: "#beb7b7"}, "slow").attr("title", "Offline (" + host + ")");
		}
	}
	
});

$.get("lab_1.tsv", function(data) {

	var lines = data.split("\n");
	var height = 0;
	for (var i = 0; i < lines.length; i++) {
		var fields = lines[i].split("\t");
		var desk = $("<div class='desk'></div>").appendTo("#lab_1")
		var top = (fields[2] - (fields[3] == "0" ? 10 : 15)) * 1.5 + 28;
		var left = (fields[1] - (fields[3] == "0" ? 15 : 10)) * 1.5 + 26;
		if (top > height) height = top;
		desk.addClass(fields[3] == "0" ? "horizontal" : "vertical");
		desk.css("top", top).css("left", left);
		if (fields.length == 5) {
			desk.attr("id", fields[4]);
			desk.css("border-color", "#7d7474");
		}
		$("#lab_1").css("height", height + 66);
	}

});

$.get("lab_2.tsv", function(data) {

	var lines = data.split("\n");
	var height = 0;
	for (var i = 0; i < lines.length; i++) {
		var fields = lines[i].split("\t");
		var desk = $("<div class='desk'></div>").appendTo("#lab_2")
		var top = (fields[2] - (fields[3] == "0" ? 10 : 15)) * 1.5 + 28;
		var left = (fields[1] - (fields[3] == "0" ? 15 : 10)) * 1.5 + 26 + 70;
		if (top > height) height = top;
		desk.addClass(fields[3] == "0" ? "horizontal" : "vertical");
		desk.css("top", top).css("left", left);
		if (fields.length == 5) {
			desk.attr("id", fields[4]);
			desk.css("border-color", "#7d7474");
		}
		$("#lab_2").css("height", height + 66);
	}

});

$(function() {

	$("#lab_1").show();
	$("#lab_1_toggle").css("background-color", "#706868");

	$("#lab_1_toggle").click(function() {
		$("#lab_1").show();
		$("#lab_2").hide();
		$("#toggles a").css("background-color", "#978f90");
		$(this).css("background-color", "#706868");
		$("h1").html("Third Floor Lab &mdash; GDC");
	});

	$("#lab_2_toggle").click(function() {
		$("#lab_2").show();
		$("#lab_1").hide();
		$("#toggles a").css("background-color", "#978f90");
		$(this).css("background-color", "#706868");
		$("h1").html("Basement Lab &mdash; GDC");
	});

});