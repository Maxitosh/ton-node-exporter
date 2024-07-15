# ton-node-exporter

A Prometheus metrics exporter for The Open Network (TON) Node.  
TON Node does not provide any built-in metrics for monitoring, so this exporter was created to fill that gap.

## Prerequisites
- Go 1.22 or later
- Docker (optional)
- Access to a TON Lite Server (Node) with a valid key.

## Installation
### From source
To build and install the exporter from source, follow these steps:
```bash
git clone https://github.com/Maxitosh/ton-node-exporter.git
cd ton-node-exporter
go build -o ton-node-exporter ./cmd/ton-node-exporter
```

### Docker
To pull the Docker image, use the following command:
```bash
docker pull ghcr.io/maxitosh/ton-node-exporter:latest
```

## Usage
### Configuration
The exporter is configured using environment variables or a `.env` file. The following variables are available:

- `EXPORTER_PORT`: The port on which the exporter will listen (default: `9100`).
- `LITE_SERVER_ADDR`: The address of the TON Lite Server (Node) to monitor.
- `LITE_SERVER_KEY`: The key to access the TON Lite Server.
- `GLOBAL_CONFIG_URL`: The URL of the TON global lite server config.

Check the .env.template file for an example and relevant environment variables. To create a .env file, copy the template:
```bash
cp .env.template .env
```

## Running
### From source
Assuming you set up the environment variables in the `.env` file, you can run the exporter with the following command:
```bash
./ton-node-exporter
```
If you want to use environment variables instead of a `.env` file, you can run the exporter like this:
```bash
export EXPORTER_PORT=9100
export LITE_SERVER_ADDR=""
export LITE_SERVER_KEY=""
export GLOBAL_CONFIG_URL=""
./ton-node-exporter
```

### Docker
To run the exporter using Docker:
```bash
docker run -d --name ton-node-exporter --env-file .env -p 9100:9100 ghcr.io/maxitosh/ton-node-exporter:latest
```

## Metrics
The exporter exposes the following metrics:

| Metric name                        | Metric type | Description                                                                 | Labels/tags | Status |  
|------------------------------------|-------------|-----------------------------------------------------------------------------|-------------|--------|  
| ton_node_master_chain_block_number | Gauge       | The current master chain block number.                                      | env         | ✅      |
| ton_node_head_lag                  | Gauge       | The lag between the current master chain block on the node and the network. | env         | ✅      |

## Testing
To run tests, execute the following command:
```bash
go test -v ./...
```

## Examples
To access the metrics, do curl request to the exporter:
```bash
curl http://localhost:9100/metrics
```
Response:
```plaintext
...
# HELP ton_node_head_lag Head block lag
# TYPE ton_node_head_lag gauge
ton_node_head_lag 0
# HELP ton_node_master_chain_block_number Master chain block number
# TYPE ton_node_master_chain_block_number gauge
ton_node_master_chain_block_number{env="global"} 3.9033746e+07
ton_node_master_chain_block_number{env="local"} 3.9033746e+07
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## Author
Created by [Maxitosh (Max Kureikin)](https://github.com/Maxitosh).
