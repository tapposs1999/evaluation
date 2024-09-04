## ⚡️ Getting Started

Clone the project

```bash
git clone https://github.com/huntabyte/tig-stack.git
```

Change the environment variables define in `.env` that are used to setup and deploy the stack
```bash
├── telegraf/
├── .env         <---
├── docker-compose.yml
├── entrypoint.sh
└── ...
```

Start the services
```bash
docker-compose up -d
```
## Docker Images Used (Official & Verified)

[**InfluxDB**](https://hub.docker.com/_/influxdb) / `2.1.1`

run: go mod init my-module-name

run go mod tidy