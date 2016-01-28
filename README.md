> Note that Dictator is currenlty in alpha stage, we are working to release a prod ready version soon.

# Redis Dictator
Redis Dictator is a simple program running next to Redis Server in a Master-Slaves topology. His role is to promote a Redis node as master througt an election process.
There is one Dictator by Redis Server on a given cluster topology, the election process are done with Zookeeper, the result of the election (the master promotion) is persisted into Zookeeper also.

## Redis at BlaBlaCar
We are using Redis at BlaBlaCar since some couples of years, for a large spectrum of usecases:
- to cache data (obvious no?)
- maintain counters (ex: "quota" usages)
- create functional locks
- stored some configuration/feature swicthes
- ...
 
Basically, our use cases involved that we should provide a "quite" high available solution, the Redis replication and persistence solutions allow us to propose an satisfying HA. Provided that a master/slaves topology should be quickly/well reconfigured in case of master failure...

## Motivation
We spent lot of time/energy to test some HA & Cluster solutions arround Redis. We put aside the idea of clusterize Redis (in term of turnkey sharding solution), it complexify your topologies and reduce drastically your consistency...We choosed to shard our dataset functionally by creating several master/slaves clusters instead of one magical auto-sharding, auto-scaling, auto-[...]ing black box.

We should admit that we are not the first to develop tooling to manipulate master/slaves toplologies. We know the official solution [Redis Sentinel](http://redis.io/topics/sentinel) but the configuration file rewriting bother us a little (note that we are in a full containers context at BlaBlaCar). We also tested the redis_elector.py script from [gutefrage](https://github.com/gutefrage/aurora-redis/blob/master/redis_elector.py) but we experienced connectivity troubles with the Zookeeper client, causing election triggering or freezing when we are manipulating (add/remove node) our Zookeeper cluster.

Finally the main motivation is surely because developping your own tool is fun and offers a lot of advantages, the solution fits perfectly to your needs, you can chose the language, merge PRs quickly...

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
        "zk_hosts": ["127.0.0.1"],
        "node" : {
            "name" : "local",
            "host" : "127.0.0.1",
            "port" : 6379
        } 
    }

The main section is composed by:
 
- `svc_name`(required): The Service/Cluster Name
- `log_level`: The log level `DEBUG`, `INFO`, `WARN`, `FATAL` (default is `INFO`)
- `zk_hosts`: An array of strings of Zookeeper nodes (default is `localhost:2181`)
- `node`: The Redis node info (detailed bellow)

The node section is composed by:

- `name`: The server name, FQDN or "display" name (default is `local`)
- `host`: The Address of the Redis server (default is `localhost`)
- `port`: The Port of the Redis server (default is `6379`)

## Run

    $ /usr/local/bin/dictator --config=/etc/dictator/dictator.conf

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request