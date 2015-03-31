# geoip-hoge

An IP address search tool.

This library reads MaxMind [GeoLite2](http://dev.maxmind.com/geoip/geoip2/geolite2/)
and [GeoIP2](http://www.maxmind.com/en/geolocation_landing) databases.


# Usage


```

# echo -e "aaaa bbbb 183.79.71.173 ccc\nxxxx yyyy 173.194.126.134 zzz\n" | geoip-hoge

aaaa bbb JP 183.79.71.173 ccc
xxx yyyy US 173.194.126.134 zzz

```

