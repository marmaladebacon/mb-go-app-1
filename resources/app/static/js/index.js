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

      // Explore default path
      index.explore();
    })
  },
  listen: function() {
    astilectron.onMessage(function(message) {
      switch (message.name) {
        case "about":
            index.about(message.payload);
            return {payload: "payload"};            
        case "check.out.menu":
            asticode.notifier.info(message.payload);
            break;
      }
    });
  }  
};