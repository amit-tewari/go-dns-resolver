# go-dns-resolver
go based dns resolver for bulk lookups

based on https://github.com/majek/goplayground -> resolve

```
$ for i in " " "-soa"  "-mx" "-txt" "-6"; do echo ; echo "==== ${i} ===="; echo -en "google.com\nyahoo.com\nreddit.com\n" | $GOPATH/bin/go-dns-resolver -server="8.8.8.8:53" ${i}; done; echo

====   ====
Server: 8.8.8.8:53, sending delay: 8.333333ms (120 pps), retry delay: 1s
google.com 1 74.125.28.100 74.125.28.101 74.125.28.102 74.125.28.113 74.125.28.138 74.125.28.139
yahoo.com 1 206.190.36.45 98.138.253.109 98.139.180.149
reddit.com 1 151.101.1.140 151.101.129.140 151.101.193.140 151.101.65.140
Resolved 3 domains in 0.018s. Average retries 1.000. Domains per second: 167.372

==== -soa ====
Server: 8.8.8.8:53, sending delay: 8.333333ms (120 pps), retry delay: 1s
google.com 6 161314624 ns1.google.com. dns-admin.google.com.
yahoo.com 6 2017070903 ns1.yahoo.com. hostmaster.yahoo-inc.com.
reddit.com 6 1 ns-557.awsdns-05.net. awsdns-hostmaster.amazon.com.
Resolved 3 domains in 0.018s. Average retries 1.000. Domains per second: 162.916

==== -mx ====
Server: 8.8.8.8:53, sending delay: 8.333333ms (120 pps), retry delay: 1s
google.com 15 alt2.aspmx.l.google.com.:30 aspmx.l.google.com.:10 alt1.aspmx.l.google.com.:20 alt4.aspmx.l.google.com.:50 alt3.aspmx.l.google.com.:40 
yahoo.com 15 mta6.am0.yahoodns.net.:1 mta7.am0.yahoodns.net.:1 mta5.am0.yahoodns.net.:1 
reddit.com 15 aspmx.l.google.com.:1 aspmx2.googlemail.com.:10 aspmx3.googlemail.com.:10 alt1.aspmx.l.google.com.:5 alt2.aspmx.l.google.com.:5 
Resolved 3 domains in 0.018s. Average retries 1.000. Domains per second: 164.969

==== -txt ====
Server: 8.8.8.8:53, sending delay: 8.333333ms (120 pps), retry delay: 1s
google.com 16 "v=spf1 include:_spf.google.com ~all", 
yahoo.com 16 "v=spf1 redirect=_spf.mail.yahoo.com", 
reddit.com 16 "v=spf1 include:_spf.google.com include:mailgun.org a:mail.reddit.com ip4:174.129.203.189 ip4:52.205.61.79 ip4:54.172.97.247 ~all", 
Resolved 3 domains in 0.018s. Average retries 1.000. Domains per second: 165.404

==== -6 ====
Server: 8.8.8.8:53, sending delay: 8.333333ms (120 pps), retry delay: 1s
google.com 28 2607:f8b0:400e:c04::71
yahoo.com 28 2001:4998:44:204::a7 2001:4998:58:c02::a9 2001:4998:c:a06::2:4008
reddit.com 28 
Resolved 3 domains in 0.019s. Average retries 1.000. Domains per second: 159.596
```

