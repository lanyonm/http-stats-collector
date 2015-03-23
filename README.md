# http-stats-collector [![Build Status](https://travis-ci.org/lanyonm/http-stats-collector.svg)](https://travis-ci.org/lanyonm/http-stats-collector)
This program collects [Navigation Timing API](http://www.html5rocks.com/en/tutorials/webperformance/basics/) data and forwards it along to the specified recorders (like [StatsD](https://github.com/etsy/statsd/)).

The JSON structure the `/nav-timing` endpoint expects is:

```javascript
{
	"page-uri": "/foo/bar",
	"nav-timing": {
		"dns":1,
		"connect":2,
		"ttfb":3,
		"basePage":4,
		"frontEnd":5
	}
}
```

Javascript that will send this information can be as simple as the following:

```javascript
<script type="text/javascript">
  window.onload = function() {
    if (window.performance && window.performance.timing) {
      var timing = window.performance.timing;
      var stats = {
        'page-uri': window.location.pathname,
        'nav-timing': {
          dns: timing.domainLookupEnd - timing.domainLookupStart,
          connect: timing.connectEnd - timing.connectStart,
          ttfb: timing.responseStart - timing.connectEnd,
          basePage: timing.responseEnd - timing.responseStart,
          frontEnd: timing.loadEventStart - timing.responseEnd
        }
      };
      if (window.XMLHttpRequest) {
        xmlhttp=new XMLHttpRequest();
        xmlhttp.open("POST", "//example.com/nav-timing", true);
        xmlhttp.setRequestHeader("Content-type", "application/json");
        xmlhttp.send(JSON.stringify(stats));
      }
    }
  };
</script>
```

The `page-uri` will be converted into the appropriate format for the recorder and the stat pushed to that recorder.

## Building
This will run tests as well.

	make

## Running

	make run

If you want to test the running system, you'll need to send it some stats:

```bash
curl -d '{"page-uri": "/foo/bar", "nav-timing":{"dns":1,"connect":2,"ttfb":3,"basePage":4,"frontEnd":5}}' -H "X-Real-Ip: 192.168.0.1" http://localhost:8080/nav-timing
```

## Test Coverage

	make cover

## TODO

- [ ] Create ability to specify a whitelist nav-timing page-uris
- [ ] Write regex for StatsD validStat
