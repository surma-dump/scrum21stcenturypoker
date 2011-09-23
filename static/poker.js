(function($, chantok) {
	var chan = new goog.appengine.Channel(chantok);
	var socket = chan.open();

	function vote(val) {
		alert("Voting for "+val);
	}

	socket.onmessage = function(msg) {
	}

	$("#buttons button").each(function() {
		$(this).click(function() {
			vote(this.value);
		});
	});
})(jQuery, window.CHANNEL_TOKEN);
