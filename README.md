# urlScan
A command line url service validation tool using golang. Support tcp, udp, http, https, ws, wss.

go run main.go -c 100 < url_list 1>ok 2>err

urlScan reads the url list from stdin, output normal url to stdout,  faulty url and error to stderr.
-c is concurrency limits. (default 1)

url_list is a list of \n separated URLs, for example:
```
https://example.com
http://example.com
udp://example.com:443
tcp://example.com:443
ws://example.com
wss://example.com
```
