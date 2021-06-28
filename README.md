# Demo SSE: Sever Sent Event

It's demo Server Sent Event (golang) with simple Javascript client.
By default, the server sends every second a message(time event) with time of now in UTC(Coordinated Universal Time).

## Message Structure

```
res.write('event: time\n');
res.write('data: 2021-06-28T14:29:17Z\n'); // time in UTC
res.write('data: \n\n');
```

```javascript
let source = new EventSource('/clocktimes');
source.addEventListener('time', timeHandler, false);
// Client can subscribe on different event.  
source.addEventListener('otherEvent', otherEventHandler, false); 
```

## CURL client
curl -N -H "Accept:text/event-stream" http://localhost:3000/clocktimes

## JavaScript client
The file client/index.html can directly consumer SSE messages with Access-Control-Allow-Origin parameter.

## Resources
[Javascript API](https://developer.mozilla.org/en-US/docs/Web/API/EventSource)

[Spec SSE](https://html.spec.whatwg.org/multipage/server-sent-events.html)