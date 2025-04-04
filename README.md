> ### ⚠️ Deprecated!!!
>
> This extension was originally created as a _proof of concept_.
> At this time, there are no maintainers available to support this extension.
>
> USE AT YOUR OWN RISK!

<br />

# xk6-amqp

A k6 extension for publishing and consuming messages from queues and exchanges.
This project utilizes [AMQP 1.4.0](https://github.com/Azure/go-amqp), the most common AMQP protocol in use today.

> :Note: This project is compatible with [AMQP 1.0].

## Build

To build a `k6` binary with this extension, first ensure you have the prerequisites:

- [Go toolchain](https://go101.org/article/go-toolchain.html)
- Git

Then:

1. Download [xk6](https://github.com/grafana/xk6):
  ```bash
  $ go install go.k6.io/xk6/cmd/xk6@latest
  ```

2. [Build the k6 binary](https://github.com/grafana/xk6#command-usage):
  ```bash
  $ xk6 build --with github.com/prmuthu/xk6-amqp@latest
  ```

## Development
To make development a little smoother, use the `Makefile` in the root folder. The default target will format your code, run tests, and create a `k6` binary with your local code rather than from GitHub.

```shell
git clone git@github.com:prmuthu/xk6-amqp.git
cd xk6-amqp
make
```

## Example

```javascript
import Amqp from 'k6/x/amqp';

const url = "amqps://localhost:443";
const session = Amqp.start({
  connection_url: url,
  username: username,
  password: password,
})
export default function () {
  
  // console.log("K6 amqp extension enabled, version: " + Amqp.version)
  
  const queueName = 'TEST.QUEUE.1::TEST.QUEUE.1';
  
  Amqp.publish({
    session: session,
    queue_name: queueName,
    body: "Ping from k6 script567"
  })

  let message = Amqp.listen({
    session: session,
    queue_name: queueName,
    auto_ack: true,
  })
  console.log("Message received: " + message);
}


```

Result output:

```plain
$ ./k6 run script.js

          /\      |‾‾| /‾‾/   /‾‾/
     /\  /  \     |  |/  /   /  /
    /  \/    \    |     (   /   ‾‾\
   /          \   |  |\  \ |  (‾)  |
  / __________ \  |__| \__\ \_____/ .io

  Connected to AMQP server
     execution: local
        script: /Users/muthukumar/Documents/demo-k6-operator/xk6-amqp/examples/test.js
        output: -

     scenarios: (100.00%) 1 scenario, 1 max VUs, 10m30s max duration (incl. graceful stop):
              * default: 1 iterations for each of 1 VUs (maxDuration: 10m0s, gracefulStop: 30s)

INFO[0003] Message received: Ping from k6 script567      source=console

     data_received........: 0 B 0 B/s
     data_sent............: 0 B 0 B/s
     iteration_duration...: avg=1.01s min=1.01s med=1.01s max=1.01s p(90)=1.01s p(95)=1.01s
     iterations...........: 1   0.986663/s
     vus..................: 1   min=0      max=1
     vus_max..............: 1   min=0      max=1


running (00m01.0s), 0/1 VUs, 1 complete and 0 interrupted iterations
default ✓ [======================================] 1 VUs  00m01.0s/10m0s  1/1 iters, 1 per VU

```

Inspect examples folder for more details.

# Testing Locally

This repository includes a [docker-compose.yml](./docker-compose.yml) file that starts RabbitMQ with Management Plugin for testing the extension locally.

> :warning: This environment is intended for testing only and should not be used for production purposes.

1. Start the docker compose environment.
   ```bash
   docker compose up -d
   ```
   Output should appear similar to the following:
   ```shell
   ✔ Network xk6-amqp_default       Created               ...    0.0s
   ✔ Container xk6-amqp-rabbitmq-1  Started               ...    0.2s
   ```
2. Use your [custom k6 binary](#build) to run a k6 test script connecting to your RabbitMQ server started in the previous step.
   ```bash
   ./k6 run examples/test.js
   ```
3. Use the RabbitMQ admin console by accessing [http://localhost:15672/](http://localhost:15672/), then login using `guest` for both the Username and Password.
   This will allow you to monitor activity within your messaging server.
