# go-dns-resolver
go based dns resolver for bulk lookups

based on https://github.com/majek/goplayground -> resolve

```
Format:
DOMAIN QueryType QueueTime TimeToResolve(in milliseconds) Answer(r)/Reply-Error
------ --------- --------- ------------------------------ ---------------------

$ for i in " " "-soa"  "-mx" "-txt" "-6"; do echo ; echo "==== ${i} ===="; echo -en "google.com\n115.email\nreddit.com\n" | $GOPATH/bin/go-dns-resolver -server="127.0.0.1:53" -retry=8s ${i}; done; echo

====   ====
Server: 127.0.0.1:53, sending delay: 8.333333ms (120 pps), retry delay: 8s
google.com 1 0 164 74.125.28.100 74.125.28.101 74.125.28.102 74.125.28.113 74.125.28.138 74.125.28.139
reddit.com 1 17 168 151.101.1.140 151.101.129.140 151.101.193.140 151.101.65.140
115.email 1 8 3733 -ERR-SERVFAIL
Resolved 3 domains in 3.742s. Average retries 1.000. Domains per second: 0.802

==== -soa ====
Server: 127.0.0.1:53, sending delay: 8.333333ms (120 pps), retry delay: 8s
google.com 6 0 35 161359078 ns4.google.com. dns-admin.google.com.
reddit.com 6 17 39 1 ns-557.awsdns-05.net. awsdns-hostmaster.amazon.com.
115.email 6 8 1284 -ERR-SERVFAIL
Resolved 3 domains in 1.294s. Average retries 1.000. Domains per second: 2.319

==== -mx ====
Server: 127.0.0.1:53, sending delay: 8.333333ms (120 pps), retry delay: 8s
google.com 15 0 1 alt1.aspmx.l.google.com.:20 alt4.aspmx.l.google.com.:50 alt3.aspmx.l.google.com.:40 aspmx.l.google.com.:10 alt2.aspmx.l.google.com.:30 
115.email 15 8 0 -ERR-SERVFAIL
reddit.com 15 17 40 aspmx.l.google.com.:1 aspmx2.googlemail.com.:10 aspmx3.googlemail.com.:10 alt1.aspmx.l.google.com.:5 alt2.aspmx.l.google.com.:5 
Resolved 3 domains in 0.057s. Average retries 1.000. Domains per second: 52.265

==== -txt ====
Server: 127.0.0.1:53, sending delay: 8.333333ms (120 pps), retry delay: 8s
115.email 16 11 0 -ERR-SERVFAIL
google.com 16 2 34 "v=spf1 include:_spf.google.com ~all", 
reddit.com 16 20 39 "v=spf1 include:_spf.google.com include:mailgun.org a:mail.reddit.com ip4:174.129.203.189 ip4:52.205.61.79 ip4:54.172.97.247 ~all", 
Resolved 3 domains in 0.060s. Average retries 1.000. Domains per second: 50.111

==== -6 ====
Server: 127.0.0.1:53, sending delay: 8.333333ms (120 pps), retry delay: 8s
google.com 28 0 1 2607:f8b0:400e:c04::65
115.email 28 8 0 -ERR-SERVFAIL
reddit.com 28 17 39 
Resolved 3 domains in 0.057s. Average retries 1.000. Domains per second: 52.781
```

