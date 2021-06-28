/* SSE */
console.log("Demo SSE")

let source = new EventSource('http://localhost:3000/clocktimes');
const content = document.getElementById("content");
const spinner = document.getElementById("spinner");
const time = document.getElementById("time");
const today = document.getElementById("today");

source.onopen = function () {
    console.log("Connection to SSE opened.");
};

source.onerror = function (e) {
    console.log("EventSource failed.");
    if (e.readyState === EventSource.CONNECTING) {
        console.log(`Reconnecting (readyState=${e.readyState})...`);

    } else if (e.readyState === EventSource.CLOSED) {
        console.log("Connection was closed.");

    } else {
        console.log("Error has occured.");
    }
};

source.addEventListener('time', message => {

    if (document.getElementById("spinner") !== null) {
        content.removeChild(spinner);
    }
    // console.log("time event:", message.data);

    let now = new Date(message.data);

    time.innerText = ('0' + now.getUTCHours()).slice(-2) + ":" + ('0' + now.getUTCMinutes()).slice(-2) + ":" + ('0' + now.getUTCSeconds()).slice(-2)
    today.innerText = ('0' + now.getUTCDate()).slice(-2) + "/" + ('0' + (now.getUTCMonth() + 1)).slice(-2) + "/" + now.getUTCFullYear()
})

