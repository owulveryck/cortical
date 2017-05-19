'use strict';
window.addEventListener("load", function(evt) {
  var output = document.getElementById("output");
  var input = document.getElementById("input");
  var ws;

  var video = document.querySelector('video');
  var canvas;

  var print = function(message) {

    //var d = document.createElement("div");
    //d.innerHTML = message;
    document.getElementById("result").innerHTML= message;
    //output.appendChild(d);
  };

  var loc = window.location, new_uri;
  if (loc.protocol === "https:") {
    new_uri = "wss:";
  } else {
    new_uri = "ws:";
  }
  new_uri += "//" + loc.host + "/ws";
  //new_uri += loc.pathname + "ws";
  ws = new WebSocket(new_uri);
  ws.onopen = function(evt) {
    print("Connected");
    //takeSnapshot();
  }
  ws.onclose = function(evt) {
    print("CLOSE");
    ws = null;
  }
  ws.onmessage = function(evt) {
    print(evt.data);
     var msg = new SpeechSynthesisUtterance(evt.data);
     msg.lang = 'en-US';
     window.speechSynthesis.speak(msg);        
  }
  ws.onerror = function(evt) {
    print("ERROR: " + evt.data);
  }

  /**
   *  generates a still frame image from the stream in the <video>
   *  appends the image to the <body>
   */
  function takeSnapshot() {
    //var img = document.querySelector('img') || document.createElement('img');
    var context;
    var width = video.offsetWidth
      , height = video.offsetHeight;

    canvas = canvas || document.createElement('canvas');
    canvas.width = width;
    canvas.height = height;

    context = canvas.getContext('2d');
    context.drawImage(video, 0, 0, width, height);

    //img.src = canvas.toDataURL('image/jpeg');
    var dataURI = canvas.toDataURL('image/jpeg')
    var byteString = dataURI.split(',')[1];

    // separate out the mime component
    var mimeString = dataURI.split(',')[0].split(':')[1].split(';')[0]
    //
    var message = {"dataURI":{}};
    message.dataURI.contentType = mimeString;
    message.dataURI.content = byteString;
    var json = JSON.stringify(message);
    ws.send(json);
    console.log("message sent");
    //ws.send(canvas.toDataURL('image/jpeg'));
    //ws.send(img.src);

    //document.body.appendChild(img);
  }

  // use MediaDevices API
  // docs: https://developer.mozilla.org/en-US/docs/Web/API/MediaDevices/getUserMedia
  if (navigator.mediaDevices) {
    // access the web cam
    var front = false;
    document.getElementById('flip-button').onclick = function() { front = !front; };

    var constraints = { video: { facingMode: (front? "user" : "environment") } };
    navigator.mediaDevices.getUserMedia(constraints)
    // permission granted:
      .then(function(stream) {
        video.src = window.URL.createObjectURL(stream);
        video.addEventListener('click', takeSnapshot);
        setInterval(takeSnapshot,3000);
      })
    // permission denied:
      .catch(function(error) {
        document.body.textContent = 'Could not access the camera. Error: ' + error.name;
      });
  }

});
