package boxapi

import (
	"context"
	"io"
	"runtime/debug"

	box "github.com/sagernet/sing-box"
	"github.com/sagernet/sing-box/include"
	"github.com/sagernet/sing-box/option"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/service"
)

var disableColor bool

func Create(nekoConfigContent []byte, externalWriter io.Writer) (*box.Box, context.CancelFunc, error) {
	ctx, cancel := context.WithCancel(context.Background())
	var options option.Options
	ctx = include.Context(service.ContextWithDefaultRegistry(ctx))
	err := options.UnmarshalJSONContext(ctx, nekoConfigContent)
	if err != nil {
		cancel()
		return nil, nil, E.Cause(err, "decode config")
	}
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

	// 核心替换逻辑：如果外部提供了 Android GUI 的 Writer，直接注入内核
	if externalWriter != nil {
		instance.SetLogWritter(externalWriter)
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
