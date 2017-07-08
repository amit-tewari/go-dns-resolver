# go-dns-resolver
go based dns resolver for bulk lookups

based on https://github.com/majek/goplayground -> resolve

```
$ go get github.com/amit-tewari/go-dns-resolver
$ echo -en "google.com\n" | $GOPATH/bin/go-dns-resolver -soa -server="8.8.8.8:53"
$ for i in " "  "-6" "-soa"; do echo ; echo "==== ${i} ===="; echo -en "google.com\nyahoo.com\nreddit.com\n" | $GOPATH/bin/go-dns-resolver -server="8.8.8.8:53" ${i}; done; echo 

====   ====
Server: 8.8.8.8:53, sending delay: 8.333333ms (120 pps), retry delay: 1s
google.com 74.125.28.100 74.125.28.101 74.125.28.102 74.125.28.113 74.125.28.138 74.125.28.139
yahoo.com 206.190.36.45 98.138.253.109 98.139.180.149
reddit.com 151.101.1.140 151.101.129.140 151.101.193.140 151.101.65.140
Resolved 3 domains in 0.018s. Average retries 1.000. Domains per second: 169.292

==== -6 ====
Server: 8.8.8.8:53, sending delay: 8.333333ms (120 pps), retry delay: 1s
google.com 2607:f8b0:400e:c04::64
yahoo.com 2001:4998:44:204::a7 2001:4998:58:c02::a9 2001:4998:c:a06::2:4008
reddit.com 
Resolved 3 domains in 0.018s. Average retries 1.000. Domains per second: 169.528

==== -soa ====
Server: 8.8.8.8:53, sending delay: 8.333333ms (120 pps), retry delay: 1s
google.com 161287261 ns4.google.com. dns-admin.google.com.
yahoo.com 2017070808 ns1.yahoo.com. hostmaster.yahoo-inc.com.
reddit.com 1 ns-557.awsdns-05.net. awsdns-hostmaster.amazon.com.
Resolved 3 domains in 0.018s. Average retries 1.000. Domains per second: 164.595
```

