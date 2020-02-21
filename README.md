# Chatbot

Chatbot, a very simple chatbot build with Go.
Chatbot is a HTTP server that listen to a endpoint (http://localthost:4000/chatbot?sentence=ping) and answer question (at least try).

## Build

You need to make sure to have installed on your machine

- Docker with Kubernetes configured
- Go (version 1.13)
- Make
- Tilt

### Make option

type `make help` to see the list of options

### Run on your machine

You can run chatbot on your machine, but first your need

- Build chatbot by running `make`
- Start the app by running `./dist/chatbot-v0.0.1-darwin`
- Visit http://localhost:3000/chatbot?sentence=ping and start to ask question
- Metric are exposed here: http://localhost:5000/metrics
- Pprof is exposed here: http://localhost:3000/debug

### Run on Kubernetes (local)

For deployment in a local Kubernetes, we will be using tilt (tilt.dev)

- Install tilt `brew install windmilleng/tap/tilt`
- Run chat bot `tilt up`
- Visit http://localhost:3000/chatbot?sentence=ping and start to ask question
- Metric are exposed here: http://localhost:5000/metrics
- Pprof is exposed here: http://localhost:3000/debug
