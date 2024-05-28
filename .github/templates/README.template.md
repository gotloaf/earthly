
<h1 align="center">
    <picture>
        <source media="(prefers-color-scheme: light)" srcset="/.github/artefacts/logo-black.svg">
        <source media="(prefers-color-scheme: dark)" srcset="/.github/artefacts/logo-white.svg">
        <img alt="earthly logo" src="/.github/artefacts/logo-black.svg">
    </picture>
</h1>

Earthly is an API for generating images of Earth. The core rendering functionality is accessible via HTTP GET, making it available either via the request library of your choice, or by generating URLs to be directly included into Markdown or other similar documents.

A hosted version is available at [`https://earthly.gotloaf.dev`](https://earthly.gotloaf.dev/?size=512&longitude=135&latitude=30&roll=-15). Swagger format API documentation is [also available](https://earthly.gotloaf.dev/swagger/index.html), compliant with OpenAPI 2.0.


## Usage

TODO

## Alternative setups

This repository contains a variety of ways to use and host earthly.

### Hosted server

To run earthly as a REST HTTP server, first you must build the docs, and then the server. You can see an example of how to do these steps in the [scripts](scripts/unix) folder.

### WASM/Cloudflare Worker

Earthly can be built for WASM, allowing it to be used on serverless platforms like Cloudflare Workers. An example of how to build for WASM is in the [scripts](scripts/unix) folder.

> [!WARNING]
> Cloudflare Workers, alongside many edge runtimes, has both a bundle size and execution time limit. To accommodate for this, the repository contains a [tinygo-based build script](scripts/windows/build_wasm_tinygo.ps1) that uses wasm-opt and wasm-strip to minimize bundle size, and an [example worker file](build/worker.mjs) to demonstrate reasonable limits.

### Command-line interface (CLI)

Earthly can be run to create one-off images using the command line. To use it, run:
```bash
go run ./cmd/cli -longitude -80 -latitude 10
```
By default, it will output PNG data to stdout, but you can also use `-output` to make it output directly to a file. You can also run `go run ./cmd/cli -help` to get information on what flags are available.

## Acknowledgements

Images of the Earth's topography are courtesy of [NASA's Visible Earth project](https://visibleearth.nasa.gov/).
