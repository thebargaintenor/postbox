# Postbox

## Running Postbox

There isn't a whole lot here.  It's a tiny utility that runs a web server with a form to upload files into a local directory set at startup, like so:

```sh
go run postbox.go ~/.postbox
```

I suppose you could compile it too.  That is also fine, if the machine you want to run postbox doesn't have a go compiler handy.

Did I mention that you should make the directory in advance?  You should, because this won't.  It's a bit simple, but it does its job.

## Accessing Postbox

Postbox assumes you have the necessary port open in your firewall too (it uses 8080).  Once the server is running, point your browser on that device to `serverhostname:8080/deposit` and add your file with the form.  That's it.  There's no retrieval yet, so it's mostly a dead drop for now.