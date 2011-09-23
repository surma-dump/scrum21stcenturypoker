(function($, token, room_name) {
	var chan = new goog.appengine.Channel(token);
	var socket = chan.open();

	function vote(val) {
		$.post("/vote", {room: room_name, vote: val});
	}

	socket.onmessage = function(msg) {
		$("<div>").text("Got: "+msg.data).appendTo("#result");
	}

	$("#buttons button").each(function() {
		$(this).click(function() {
			vote(this.value);
		});
	});
})(jQuery, window.CHANNEL_TOKEN, window.ROOM);
