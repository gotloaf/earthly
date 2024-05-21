
import "./wasm_exec.js"
import mod from "./earthly.wasm";

async function run(config) {
    const earth = await import(
        config.size * config.zoom > 1024.0 ?
        "./equirectangular/earth_2x.png" :
        "./equirectangular/earth_1x.png"
    );

    console.log(earth);
    console.log(config);

    const go = new Go();

    const mountedPromise = new Promise((resolve, reject) => {
		globalThis._earthlyResolve = resolve;
	});

    const instance = new WebAssembly.Instance(mod, {
        ...go.importObject,
    });

    const goExitHandle = go.run(instance);

    await mountedPromise;

    const output = earthlyGenerate(JSON.stringify(config), new Uint8Array(earth.default));

    earthlyShutdown();

    await goExitHandle;

    if (!(output instanceof Uint8Array) || output.length == 0) {
        return Response.json({
            "error": "An internal error occurred."
        });
    }

    return new Response(
        output
    );
}

/**
 * @param {Request} request
 */
function requestToConfiguration(request) {
    const { searchParams } = new URL(request.url);
    const errors = [];

    let size = parseInt(searchParams.get('size') || "512");
    if (isNaN(size)) {
        size = 512;
        errors.push("parameter `size` could not be parsed");
    }
    if (size < 16 || size > 2048) {
        size = Math.max(16, Math.min(2048, size));
        errors.push("parameter `size` was outside range [16-2048]");
    }

    let latitude = parseFloat(searchParams.get('latitude') || "0.0");
    if (isNaN(latitude)) {
        latitude = 0.0;
        errors.push("parameter `latitude` could not be parsed");
    }

    let longitude = parseFloat(searchParams.get('longitude') || "0.0");
    if (isNaN(longitude)) {
        longitude = 0.0;
        errors.push("parameter `longitude` could not be parsed");
    }

    let roll = parseFloat(searchParams.get('roll') || "0.0");
    if (isNaN(roll)) {
        roll = 0.0;
        errors.push("parameter `roll` could not be parsed");
    }

    let zoom = parseFloat(searchParams.get('zoom') || "1.0");
    if (isNaN(zoom)) {
        zoom = 1.0;
        errors.push("parameter `zoom` could not be parsed");
    }

    return {
		size: size,
		background: [0, 0, 0, 0],
		latitude: latitude,
		longitude: longitude,
		roll: roll,
		halo: true,
		radius: zoom,
	};
}

export default {
    async fetch(request, env, ctx) {
        const config = requestToConfiguration(request);
        return await run(config);
    }
}
