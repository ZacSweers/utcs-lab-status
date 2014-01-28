var Name = "Lab Status Repeater";
var Version = "0.0.2";
var Identifier = Name + " " + Version;

var request = require("request");
var app = require("http").createServer(handler),
	io = require("socket.io").listen(app, {log: false}),
	fs = require("fs")

poll();

console.log(Identifier);
app.listen(80);
var status

function handler(req, res) {
	if (req.url == "" || req.url == "/") {
		req.url = "/index.html";
	}
	fs.readFile(__dirname + "/public" + req.url,
	function (err, data) {
		if (err) {
			res.writeHead(500);
			return res.end();
		}
		res.writeHead(200);
		res.end(data);
	});
}

function poll() {
	request({
		uri: "http://www.cs.utexas.edu/~yeh/cgi-bin/poll.scgi",
	}, function(error, response, data) {
		status = data
		io.sockets.emit("status", data)
	});
}

setInterval(function() {
	poll();
}, 60*1000);

io.sockets.on("connection", function(socket) {
	socket.emit("status", status);
});