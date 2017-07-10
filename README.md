# go-dns-resolver
go based dns resolver for bulk lookups

based on https://github.com/majek/goplayground -> resolve

```
Format:
DOMAIN QueryType QueueTime TimeToResolve(in milliseconds) Answer(r)/Reply-Error
------ --------- --------- ------------------------------ ---------------------
$ for i in " " "-soa"  "-mx" "-txt" "-6"; do echo ; echo "==== ${i} ===="; echo -en "google.com\n115.email\nreddit.com\n" | $GOPATH/bin/go-dns-resolver -serverPool=8.8.8.8,8.8.4.4 -retry=8s ${i}; done; echo

====   ====
Servers: [8.8.8.8 8.8.4.4], sending delay: 8.333333ms (120 pps), retry delay: 8s
google.com 1 0 2 74.125.28.100 74.125.28.101 74.125.28.102 74.125.28.113 74.125.28.138 74.125.28.139
reddit.com 1 17 0 151.101.1.140 151.101.129.140 151.101.193.140 151.101.65.140
115.email 1 8 1081 -ERR-SERVFAIL
Resolved 3 domains in 1.090s. Average retries 1.000. Domains per second: 2.752

==== -soa ====
Servers: [8.8.8.8 8.8.4.4], sending delay: 8.333333ms (120 pps), retry delay: 8s
google.com 6 0 2 161379310 ns2.google.com. dns-admin.google.com.
reddit.com 6 17 0 1 ns-557.awsdns-05.net. awsdns-hostmaster.amazon.com.
115.email 6 8 1063 -ERR-SERVFAIL
Resolved 3 domains in 1.072s. Average retries 1.000. Domains per second: 2.798

==== -mx ====
Servers: [8.8.8.8 8.8.4.4], sending delay: 8.333333ms (120 pps), retry delay: 8s
google.com 15 0 3 alt1.aspmx.l.google.com.:20 aspmx.l.google.com.:10 alt4.aspmx.l.google.com.:50 alt2.aspmx.l.google.com.:30 alt3.aspmx.l.google.com.:40 
reddit.com 15 17 9 aspmx.l.google.com.:1 aspmx2.googlemail.com.:10 aspmx3.googlemail.com.:10 alt1.aspmx.l.google.com.:5 alt2.aspmx.l.google.com.:5 
115.email 15 8 1051 -ERR-SERVFAIL
Resolved 3 domains in 1.060s. Average retries 1.000. Domains per second: 2.830

==== -txt ====
Servers: [8.8.8.8 8.8.4.4], sending delay: 8.333333ms (120 pps), retry delay: 8s
google.com 16 0 2 "v=spf1 include:_spf.google.com ~all", 
reddit.com 16 17 0 "v=spf1 include:_spf.google.com include:mailgun.org a:mail.reddit.com ip4:174.129.203.189 ip4:52.205.61.79 ip4:54.172.97.247 ~all", 
115.email 16 8 1118 -ERR-SERVFAIL
Resolved 3 domains in 1.128s. Average retries 1.000. Domains per second: 2.661

==== -6 ====
Servers: [8.8.8.8 8.8.4.4], sending delay: 8.333333ms (120 pps), retry delay: 8s
google.com 28 0 4 2607:f8b0:400e:c04::8a
reddit.com 28 17 0 
115.email 28 8 1089 -ERR-SERVFAIL
Resolved 3 domains in 1.098s. Average retries 1.000. Domains per second: 2.731

```

