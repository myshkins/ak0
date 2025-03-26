### Overview
This is the code for my personal website. It's a bit overengineered because I wanted to explore/learn some new things in building this. I will probably continue to add more as well :)

### Architecture
Everything on my site (homepage and blog) are static pages. Everything lives as an embedded file system in my webapp server (written in Go), and is served from there. Originally the frontend was built with react and vite, but that started to feel bloated so now it's just vanilla html and css.

The server is intrumented with opentelemtry to collect metrics. Those metrics are sent to an opentelemetry collector which then exports them to a prometheus backend. Grafana is used for visualizing the metrics (see blog post about it [here](https://ak0.io/blog/How%20I%20Monitor%20This%20Website).

I'm using docker compose to deploy everything. I've been learning a bit of Nix lately, so I am using that to build a docker image that includes the static files for my frontend and a go binary for my backend. Since that image has no real dependencies, it's built on top of a distroless debian image.

