package boxapi

import (
	"context"
	"runtime/debug"

	box "github.com/sagernet/sing-box"
	"github.com/sagernet/sing-box/include"
	"github.com/sagernet/sing-box/log"
	"github.com/sagernet/sing-box/option"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/service"
)

var disableColor bool

// Create builds and starts a box instance from a config. When logger is
// non-nil it is used as the box's log factory (e.g. so a GUI can redirect
// kernel logs to its own writer), otherwise a default factory is created.
func Create(nekoConfigContent []byte, logger log.Factory) (*box.Box, context.CancelFunc, error) {
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
		Logger:  logger,
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
