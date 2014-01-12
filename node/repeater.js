var net = require("net");
var tls = require("tls");
var fs = require("fs");

var Name = "Lab Status Repeater";
var Version = "0.0.1";
var Identifier = Name + " " + Version;

var app = require("http").createServer(handler),
	io = require("./Modules/socket.io").listen(app, {log: false}),
	fs = require("fs")

app.listen(80);
var status

function handler (req, res) {
	if (req.url == "" || req.url == "/") {
		req.url = "/Index.html";
	}
	fs.readFile(__dirname + "/Public" + req.url,
	function (err, data) {
		if (err) {
			res.writeHead(500);
			return res.end();
		}
		res.writeHead(200);
		res.end(data);
	});
}

io.sockets.on("connection", function (socket) {
	socket.emit("status", status);
});

var options = {
	key: fs.readFileSync("Certificates/Server.key"),
	cert: fs.readFileSync("Certificates/Server.pem"),
	ca: [ fs.readFileSync("Certificates/Client.pem") ],
	requestCert: true
};

var server = tls.createServer(options, function(stream) {
	stream.setEncoding("utf8");
	stream.on("data", function(data) {
		status = data
		io.sockets.emit("status", data)
	});
});

server.listen(8000, function() {
	console.log(Identifier);
});
