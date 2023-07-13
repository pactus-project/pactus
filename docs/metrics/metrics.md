# Metrics

Before starting to get metrics of the Pactus software, there are a few important steps that need to be followed.
Please follow the instructions below:

1. after you `Init` the Pactus in your working directory you need to update `config.toml` file before starting the node.

by default metrics in the Pactus is not enable!

        1.1. open config.toml file which its in your working directory.
        1.2. go to network section.
        1.3.find enable_metrics = false
        1.4. change it to enable_metrics = true
        1.5. done!, now you can save the file.

2. start the Pactus node.


3. now metrics are available, for example if you want to get metrics in the Pactus testnet you can find it at below url

```
http://localhost:8080/metrics/prometheus
```

4. if you want to use metrics in monitoring tools like `Prometheus` you can use the `docker-compose` file.

just be careful about the addresses and ports that they're existing in `docker-compose` and `prometheus` yml.

if you run the Pactus node on docker container you should change targets of `pactus-testnet` job in `prometheus.yml`


