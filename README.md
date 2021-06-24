# Example Full-Stack Challenge Solution

Time Spent: 2hrs 45 mins

This is a solution for the full-stack challenge, includes golang, javascript and HTML. In order to run the server,
you'll need a recent go compiler. Then, you can compile and run the solution with the following set of commands:

```
[bash] git clone git@github.com:djthorpe/coding-challenge.git
[bash] cd coding-challenge
[bash] go run ./cmd/spamserver data/reports.json
```

The single argument is the path to the data file, and loads the data into the memory. It serves on port `:8080` 
on localhost (so then navigate to http://localhost:8080/html/ to see the frontend). You can run it on a 
different port with the -addr flag, for example:

```
[bash] go run ./cmd/spamserver -addr :9001 data/reports.json
```

The unit tests can be run using:

```
[bash] git clone git@github.com:djthorpe/coding-challenge.git
[bash] cd coding-challenge
[bash] go test -v ./pkg/...
```

## Commentary

The command-line tool is in the `cmd` folder and compiles in the HTML with the go code, so that it can create
a single binary for serving.

There are three packages:

  * `pkg/server` implements a basic HTTP server;
  * `pkg/backend` implements the backend database and handler functions;
  * `pkg/schema` defines the schema for reports and updating a report.

There is also a `html` folder which contains the frontend page and javascript.

## Appendix: Challenge description

This challenge imagines that we have a social media platform that is under attack from spam. We have implemented a reporting system for users that lets them report spam to the platform, and our spam protection team.  

The challenge is to create a small full stack application for our spam protection team that consists of a server and a web based UI in order to manage reported content.

The UI should look something like:

![Reporting listing](images/wireframe.png)

We provide an example listing response ([`data/reports.json`](data/reports.json)) that you can use as the basis of your listing. Please fill the appropriate fields in the wireframe, ignore the "Details" link.

Furthermore we need a way to _block_ the content and _resolve_ those reports. The two buttons in the UI should do a call to your backend service in order to block the content or to resolve the ticket. You are free to implement the blocking as you want, however the resolving should be defined as a `PUT` request to an endpoint with this structure `/reports/:reportId`. An example request for how to update a report is in [`data/update_ticket_request.json`](data/update_ticket_request.json).


- **`Block`:** Means that the content should no longer be available to users
- **`Resolve`:** Means that the report is considered "resolved", and it is no longer visible to the spam protection team
- **`Details`:** Functionality can be ignored.
