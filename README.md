> Note that Dictator is currenlty in alpha stage, we are working to release a prod ready version soon.

# Redis Dictator
Redis Dictator is a simple program running next to Redis Server in a Master-Slaves topology. His role is to promote a Redis node as master through an election process.
There is one Dictator by Redis Server on a given cluster topology, the election process are done with Zookeeper, the result of the election (the master promotion) is persisted into Zookeeper also.

## Redis at BlaBlaCar
We are using Redis at BlaBlaCar since some couples of years, for a large spectrum of usecases:
- to cache data (obvious no?)
- maintain counters (ex: "quota" usages)
- create functional locks
- stored some configuration/feature swicthes
- ...
 
Basically, our use cases involved that we should provide a "quite" high available solution, the Redis replication and persistence solutions allow us to propose a satisfying HA. Provided that a master/slaves topology should be quickly/well reconfigured in case of master failure...

## Motivation
We spent lot of time/energy to test some HA & Cluster solutions around Redis. We put aside the idea of clusterize Redis (in term of turnkey sharding solution), it complexify your topologies and reduce drastically your consistency... We choosed to shard our dataset functionally by creating several master/slaves clusters instead of one magical auto-sharding, auto-scaling, auto-[...]ing black box.

We are not the first to develop tooling to manipulate master/slaves toplologies. The most known is surely [Redis Sentinel](http://redis.io/topics/sentinel) but the configuration file rewriting bother us a little (note that we are in a full containers context at BlaBlaCar).  By the way, we should admit that our main motivation is certainly because developping your own tool is fun and offers a lot of advantages, the solution fits perfectly to your needs, you can chose the language, merge PRs quickly...

## Building from source
Make sure to have a go 1.4 or higher:
    
    $ go version
    $ go version go1.5.3 linux/amd64

Clone the Github repository:

    $ cd /tmp
    $ go get github.com/mattn/gom
    $ git clone https://github.com/mfouilleul/redis-dictator
    $ cd redis-dictator

Build & install with make:

    $ make
    $ sudo make install
    
### Configuration
    
The Dictator configuration is a JSON file located at `/etc/dictator/dictator.conf`

    $ cat /etc/dictator/dictator.conf
    {
        "svc_name" : "default",
        "log_level" : "INFO",
        "http_port": 8000,
        "zk_hosts": ["127.0.0.1"],
        "node" : {
            "name" : "local",
            "host" : "127.0.0.1",
            "port" : 6379,
            "loading_timeout" : 30
        } 
    }

The main section is composed by:
 
- `svc_name`: The Service/Cluster Name (default is `local`)
- `log_level`: The log level `DEBUG`, `INFO`, `WARN`, `FATAL` (default is `INFO`)
- `http_port`: The port of the HTTP listener, used for interact with Dictator (Disable/Enable)
- `zk_hosts`: An array of strings of Zookeeper nodes (default is `localhost:2181`)
- `node`: The Redis node info (detailed bellow)

The node section is composed by:

- `name`: The server name, FQDN or "display" name (default is `local`)
- `host`: The Address of the Redis server (default is `localhost`)
- `port`: The Port of the Redis server (default is `6379`)
- `loading_timeout`: The time in second that Dictator accepts to wait during Redis loads its dataset to memory

## Run

    $ /usr/local/bin/dictator --config=/etc/dictator/dictator.conf

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request
