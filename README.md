# go-dns-resolver
go based dns resolver for bulk lookups

based on https://github.com/majek/goplayground -> resolve

SOA queries
$ echo -en "google.com\n" | $GOPATH/bin/go-dns-resolver -soa -server="8.8.8.8:53"
