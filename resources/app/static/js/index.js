let index = {
  about: function(html) {
    let c = document.createElement("div");
    c.innerHTML = html;
    asticode.modaler.setContent(c);
    asticode.modaler.show();
  },
  init: function() {
    // Init
    asticode.loader.init();
    asticode.modaler.init();
    asticode.notifier.init();
    // Wait for astilectron to be ready
    document.addEventListener('astilectron-ready', function() {
      // Listen
      index.listen();
    })
  },

  makeNewTrackerDiv: function(payload) {
    let idAttributeVal = "tracker-" + payload.Symbol;
    let c = document.getElementById(idAttributeVal)
    if(c==null){
      c = document.createElement("div");
      c.setAttribute("id", idAttributeVal);
      c.setAttribute("class", "panel");
      let parent = document.getElementById("tracker-rows");
      parent.appendChild(c);
    }
    c.innerHTML = payload.Text;
  },
  listen: function() {
    astilectron.onMessage(function(message) {
      switch (message.name) {
        case "about":
          index.about(message.payload);
          return {payload: "payload"};            
        case "time.test":
          document.getElementById("timetext").innerHTML = message.payload;
          break;
        case "check.out.menu":
          asticode.notifier.info(message.payload);
          break;
        case "test.struct":
          console.log(message);
          let test = JSON.parse(message.payload);
          console.log(test);
          break;
        case "ticker.track":
          index.makeNewTrackerDiv(message.payload);
          break;
      }
    });
  }  
};