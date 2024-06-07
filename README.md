# Card Cache

## What is this?
<p>Card Cache is a simple library written in Go for distributed caching of key-value pairs in memory. My aim is to make this similar to something like Redis w/ a simple API.
This is a work in progress and I'm still working on it. This is mainly just for learning purposes, but I am hoping I could use this in my Svoker project, maybe for the betting system.</p>

## How to use?

1. Install & Run
```
git clone git@github.com:tom773/cardcache.git
cd cardcache
make run
```
2. Usage. At the moment, netcat is the best way to interact with the server.
```
nc localhost 42069
```
3. Commands
```
SET <key> <value>
e.g SET foo bar

GET <key>
e.g GET foo
Outputs: bar

DEL <key>
e.g DEL foo
```

### Features
- [x] Basic GET, SET, DEL commands
- [x] Basic TCP server
- [ ] Pub/Sub via channels / websockets
- [ ] Types (int, string, JSON)
- [ ] Sharding
- [ ] Persistence
- [ ] Clustering
- [ ] Benchmarking
- [ ] Authentication

