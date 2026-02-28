package boxmain

import (
	"context"
	"runtime/debug"

	box "github.com/sagernet/sing-box"
	"github.com/sagernet/sing-box/include"
	"github.com/sagernet/sing-box/option"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/service"
)

func Create(nekoConfigContent []byte) (*box.Box, context.CancelFunc, error) {
	ctx, cancel := context.WithCancel(context.Background())
	var options option.Options
	ctx = include.Context(service.ContextWithDefaultRegistry(ctx))
	err := options.UnmarshalJSONContext(ctx, nekoConfigContent)
	if err != nil {
		cancel()
		return nil, nil, E.Cause(err, "decode config")
	}
	//
	if disableColor {
		if options.Log == nil {
			options.Log = &option.LogOptions{}
		}
		options.Log.DisableColor = true
	}
	instance, err := box.New(box.Options{
		Context: ctx,
		Options: options,
	})
	if err != nil {
		cancel()
		return nil, nil, E.Cause(err, "create service")
	}

	err = instance.Start()
	if err != nil {
		cancel()
		return nil, nil, E.Cause(err, "start service")
	}
	debug.FreeOSMemory()
	return instance, cancel, nil
}

func SetDisableColor(dc bool) {
	disableColor = dc
}
