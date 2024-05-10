package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/gotloaf/earthly"
)

type EarthlyWASM struct {
	generateCallback	js.Func

	/*
	Channel & function that indicates when the WASM program should shut down.
	If we shut down too early, the runtime doesn't stay alive for when it needs to be called.
	If we don't trigger this, the shutdown is not graceful.
	*/
	done				chan struct{}
	shutdownCallback	js.Func
}

func New() *EarthlyWASM {
	return &EarthlyWASM {
		done: make(chan struct{}),
	}
}

func (app *EarthlyWASM) Initialize() {
	app.generateCallback = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 1 {
			return "earthlyGenerate takes one (1) argument for the earth generation config"
		}

		var config *earthly.EarthlyConfig
		config_bytes := []byte(args[0].String())

		err := json.Unmarshal(
			config_bytes,
			&config,
		)

		if err != nil {
			fmt.Println("Encountered error while trying to unmarshal earthly config: ", err)
			return 1
		}

		buffer := config.Generate()
		data_bytes := buffer.Bytes()
		js_buffer := js.Global().Get("Uint8Array").New(len(data_bytes))
		js.CopyBytesToJS(js_buffer, data_bytes)

		return js_buffer
	})
	js.Global().Set("earthlyGenerate", app.generateCallback)

	app.shutdownCallback = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// Send signal to channel to allow Go to gracefully terminate
		app.done <- struct{}{}
		return nil
	})
	js.Global().Set("earthlyShutdown", app.shutdownCallback)
}

func (app *EarthlyWASM) Teardown() {
	// Release all resources
	app.shutdownCallback.Release()
}

func main() {
	fmt.Println("Initialized earthly go package")
	earthly := New()
	earthly.Initialize()

	// Now that initialization is complete, notify JS that we are mounted and ready to use
	js.Global().Get("_earthlyResolve").Invoke()

	// Wait for shutdown
	<-earthly.done

	// Teardown and release resources
	earthly.Teardown()
	fmt.Println("earthly shutting down")
}
