# cclogconv (former name is geoip-hoge)

This tool is looking up country-code from MaxMind GeoIP2/GeoLite2 databases.

* MaxMind [GeoLite2](http://dev.maxmind.com/geoip/geoip2/geolite2/)
* MaxMind [GeoIP2](http://www.maxmind.com/en/geolocation_landing)

# Usage


```

# echo -e "aaaa bbbb 183.79.71.173 ccc\nxxxx yyyy 173.194.126.134 zzz\n" | cclogconv --data /usr/share/GeoIP/GeoLite2-Country.mmdb

aaaa bbb JP 183.79.71.173 ccc
xxx yyyy US 173.194.126.134 zzz

```

