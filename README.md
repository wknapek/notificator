Library implements an HTTP notification client. A client is
configured with a URL to which notifications are sent. It implements a
function that takes messages and notifies about them by sending HTTP
POST requests to the configured URL with the message content in the request
body.
