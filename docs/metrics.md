# Metrics

Pactus node offers [Prometheus](https://prometheus.io/) metrics to monitor and analyze various network-related and resource statistics.

# Usage

To activate this feature, inside the `config.toml`,  set the `enable_metrics` parameter to true.
Also, ensure that the HTTP module is enabled. You can enable HTTP module under the `[http]` section of the `config.toml` file. 
Once enabled, the metrics can be accessed at [http://localhost:80/metrics/prometheus](http://localhost:80/metrics/prometheus) (this url going to be use by prometheus).

> NOTE: if you are running Pactus with docker image, make sure to expose `:80` port.

After these changes, restart the Pactus node; you should now be able to view the metrics.

## Prometheus Configuration

Prometheus is an open-source monitoring and alerting tool that facilitates the collection and processing of metrics. A common method of running Prometheus is via Docker containers. To use Prometheus with Docker, follow these steps:

1- Ensure [Docker](https://www.docker.com/) is installed on your system.

2- Pull the Prometheus Docker image:

```text
docker pull prom/prometheus
```

3- Create a configuration file named `prometheus.yml` to define the Prometheus configuration. You can refer to the Prometheus [documentation](https://prometheus.io/docs/prometheus/latest/configuration/configuration/) for more guidance. As an example, here's a simple configuration:

```yaml
scrape_configs:
  - job_name: "prometheus"
    scrape_interval: 1m
    static_configs:
      - targets: [ "127.0.0.1:9090" ]

  - job_name: "pactus-metrics"
    metrics_path: /metrics/prometheus
    static_configs:
      - targets: [ "node_IP:80" ]
```
> NOTE: you should relace node_IP with your server IP or 127.0.0.1 if you are running Pactus locally.

4- Start Prometheus as a Docker container:

```text
docker run -p 9090:9090 -v /path/to/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus
```
Replace `/path/to/prometheus.yml` with the actual path to your configuration file.

5- Prometheus should now be up and running. Access the Prometheus web interface by visiting [http://localhost:9090/](http://localhost:9090/) in your web browser.
