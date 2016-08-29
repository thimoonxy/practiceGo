# practiceGo


### curl.go

- Practice using standard lib such as net/http, regexp etc. to mimic a common cmd //curl -ksvo $url//

```
$ go run curl.go www.cbs.com
GET http://www.cbs.com HTTP/1.1
Host: www.cbs.com
Address: [64.30.228.50 64.30.228.49]
User-Agent: Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8

HTTP/1.1 200 OK
X-Varnish: 810930990
Via: 1.1 varnish
X-Cache: MISS
X-Hit-Count: 0
Server: Apache
Cache-Control: no-cache
Content-Type: text/html; charset=utf-8
Date: Mon, 29 Aug 2016 10:35:38 GMT
Set-Cookie: graph=%7B%22sv_campaign%22%3A%7B%22ftag%22%3Anull%2C%22siteID%22%3Anull%2C%22clickID%22%3Anull%2C%22subID1%22%3Anull%2C%22subID2%22%3Anull%2C%22subID3%22%3Anull%2C%22cbsClick%22%3Anull%2C%22sharedID%22%3Anull%2C%22promo%22%3Anull%2C%22cbscidmt%22%3Anull%7D%2C%22cookiePath%22%3A%22%5C%2F%22%7D; path=/
Accept-Ranges: bytes
Age: 0
Expires: Sat, 26 Jul 1997 05:00:00 GMT
X-Frame-Options: SAMEORIGIN
Vary: Accept-Encoding
X-Real-Server: cbscom_www_php_vip1
Connection: keep-alive

```

```
$ go run curl.go www.cnet.com
GET http://www.cnet.com HTTP/1.1
Host: www.cnet.com
Address: [104.76.0.105]
User-Agent: Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8

HTTP/1.1 200 OK
Server: nginx
Cache-Control: max-age=300, private
X-Tx-Id: 8dd10f12-3b15-4463-9d1d-3ef4c2b0000a
Access-Control-Allow-Origin: http://www.cnet.com
Accept-Ranges: bytes
Connection: keep-alive
Set-Cookie: fly_geo={"countryCode": "cn"}; expires=Mon, 05-Sep-2016 10:37:39 GMT; path=/; domain=.cnet.com
Expires: Mon, 29 Aug 2016 10:40:47 GMT
Content-Type: text/html; charset=UTF-8
Last-Modified: Mon, 29 Aug 2016 10:35:47 GMT
Content-Security-Policy: frame-ancestors 'self' *.cnet.com *.stumbleupon.com;
X-Frame-Options: SAMEORIGIN
Vary: Accept-Encoding
Date: Mon, 29 Aug 2016 10:37:39 GMT

```

### bytesize.go

- Awesome const definition sample
