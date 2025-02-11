### Overview
This is the code for my personal website. It's a bit overengineered because I wanted to explore/learn some new things in building this. I will probably continue to add more as well :)

### Architecture
The frontend is a static site built with vite and react. It's fairly old, and I could build something much nicer, so updating that is on the todo list.

I'm using docker compose to deploy everything. I've been learning a bit of Nix lately, so I am using that to build a docker image that includes the static files for my frontend and a go binary for my backend. Since that image has no real dependencies, it's built on top of a distroless debian image.

The backend is intrumented with opentelemtry to collect metrics. Those metrics are sent to an opentelemetry collector which then exports them to a prometheus backend. Grafana is used for visualizing the metrics.
