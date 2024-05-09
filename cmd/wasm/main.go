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

		var config earthly.EarthlyConfig

		err := json.Unmarshal(
			[]byte(args[0].String()),
			&config,
		)

		if err != nil {
			return fmt.Sprintf("Encountered error while trying to unmarshal earthly config: %s", err)
		}

		//buffer := config.Generate()

		return 0
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
