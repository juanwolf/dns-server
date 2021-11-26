# DNS Server

A simple authoritative DNS server using the amazing DNS library from @miekg: github.com/miekg/dns

## Goal

The goal of this exercise was just to dig on how Cloudflare implemented CName flatenning as it saved my ass in some projects. (Especially to set CNAMEs at the root level)

In here, you have a test zone record that contains CNAMEs pointing to gandi mail servers. If you query this DNS server, you should get an A record instead of a CNAME. Incredible isn't it?

As we can see in a few lines, we can get a similar behaviour. Obviously some stuff has been hardcoded to gain time (I only had an afternoon and rusty Golang memories to spend on this server). But was pretty fun anyway so :grin: .

## How To Run

Just don't but if you insist:

```
go run .
```

## How to test

```
dig @127.0.0.1 -p 5553 pop.mydomain.com
```

:warning: DO NOT RUN IN PRODUCTION OR ANYWHERE ELSE

## LICENSE

MIT
