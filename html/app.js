
let ReportsNode = null;

// This method is called when the document has loaded, and initates
// the fetch of reports
window.addEventListener('DOMContentLoaded', () => {
    ReportsNode = document.querySelector('#reports');
    LoadReports();
});

function LoadReports() {
    fetch("/reports")
        .then((response) => response.json())
        .then(data => {
            let reports = [];
            data.forEach((data) => {
                reports.push(new Report(data));
            });
            RenderTable(reports, (evt, action, report) => {
                switch (action) {
                    case "block":
                        UpdateReport(`/reports/${report.id}`, "BLOCKED");
                        break;
                    case "resolve":
                        UpdateReport(`/reports/${report.id}`, "RESOLVED");
                        break;
                    default:
                        console.warn(`Unknown action: ${action}`);
                }
            });
        });
}

// RenderTable function adds rows to the table and attaches callback function
function RenderTable(reports, callback) {
    // Empty the reports table
    ReportsNode.innerHTML = "";

    // Render the reports within the table, pass callback for each
    // rendered row which is called when the button is pressed
    reports.forEach((report) => {
        ReportsNode.appendChild(report.Render(callback));
    });
}

function UpdateReport(path, action) {
    const req = {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            "ticketState": action
        })
    };

    fetch(path, req)
        .then((response) => response.json())
        .then(data => {
            // Let's simply reload all the reports again
            LoadReports();
        });
}

// Report represents a single spam report
class Report {
    constructor(data) {
        this.$data = data;
    }

    // Properties
    get id() {
        return this.$data.id;
    }

    // Render returns a report as a DOM Node
    Render(callback) {
        let node = document.createElement("tr");
        node.innerHTML = `
        <th scope="row">
            Id: ${this.$data.id}<br>
            State: ${this.$data.state}<br>
            <a href="#">Details</a><br>
        </th>
        <td>
            Type: ${this.$data.payload.reportType}<br>
            Message: ${this.$data.payload.message}<br>
        </td>
        <td>
            <button class="btn btn-primary m-1 block">Block</button><br>
            <button class="btn btn-primary m-1 resolve">Resolve</button><br>
        </td>
        `
        // Attach callback to buttons
        node.querySelector('button.block').addEventListener('click', (evt) => {
            callback(evt, 'block', this);
        });
        node.querySelector('button.resolve').addEventListener('click', (evt) => {
            callback(evt, 'resolve', this);
        });

        return node;
    }
}
