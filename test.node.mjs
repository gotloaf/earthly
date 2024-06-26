"use strict";

import crypto from "crypto";
import fs from "fs";
import fs_promises from "fs/promises"
import { createRequire } from "module";
import { TextDecoder, TextEncoder } from "util";

globalThis.require = createRequire(import.meta.url);
globalThis.fs = fs;
globalThis.TextEncoder = TextEncoder;
globalThis.TextDecoder = TextDecoder;
globalThis.crypto ??= crypto;

globalThis.require("./build/wasm_exec_tinygo_edge");

(async () => {
	// Load assets ahead of time
	const earthJPGBuffer = await fs_promises.readFile("./build/equirectangular/earth_1x.jpg");
	const earthJPG = new Uint8Array(earthJPGBuffer);

	const go = new Go();

	process.on("exit", (code) => {
		// Emit errors if a clean exit did not result in the Go runtime being gracefully shut down
		if (code === 0 && !go.exited) {
			go._pendingEvent = { id: 0 };
			go._resume();
		}
	})

	const wasmData = await fs_promises.readFile("./build/earthly_tinygo.wasm");
	const instantiated = await WebAssembly.instantiate(wasmData, go.importObject);

	// In cmd/wasm/main.go we expect that a signal handler `_earthlyResolve` exists that allows us to know our program has mounted.
	// To do this we're going to make a Promise that attaches its resolve to the global scope, and then we can wait on that promise.
	const mountedPromise = new Promise((resolve, reject) => {
		globalThis._earthlyResolve = resolve;
	});
	// We should also await on the exit handle once we have done and shut down everything, but we don't want to do that yet.
	const goExitHandle = go.run(instantiated.instance);

	// Go should be in the process of starting. Wait for Go to signal to us that the functions are ready to use.
	await mountedPromise;

	const output = earthlyGenerate(JSON.stringify({
		size: 512,
		background: [0, 0, 0, 0],
		longitude: -138,
		latitude: 36,
		roll: -30,
		halo: true,
		radius: 1.0,
	}), earthJPG);

	console.log(output);

	if (!(output instanceof Uint8Array) || output.length == 0) {
		console.log(`Test failed, output should have been Uint8Array of greater than 0 size but it was: ${output}`);
		process.exit(-1);
	}

	await fs_promises.writeFile("output.test.png", output);

	// Shut down earthly, which should allow the Go runtime to gracefully terminate.
	earthlyShutdown();
	// Now wait for that termination.
	await goExitHandle;

	// If we've reached this point, the lifecycle has completed with no issues.
})();
