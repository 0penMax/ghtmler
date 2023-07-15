package builder

import "strings"

const jsScript = `<script type="text/javascript">

function tryConnectToReload(address) {
  var conn = new WebSocket(address);

  conn.onclose = function() {
    setTimeout(function() {
      tryConnectToReload(address);
    }, 2000);
  };

  conn.onmessage = function(evt) {
    location.reload()
  };
}

try {
  if (window["WebSocket"]) {
      tryConnectToReload("ws://localhost:12450/reload");
  } else {
    console.log("Your browser does not support WebSockets, cannot connect to the Reload service.");
  }
} catch (ex) {
  console.error('Exception during connecting to Reload:', ex);
}
</script>`

func injectLiveReloadScript(content []string) []string {
	for i := range content {
		content[i] = strings.Replace(content[i], "</body>", jsScript+"</body>", -1)
	}
	return content
}
