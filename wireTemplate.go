package main

const wireTemplate = `
//go:build wireinject
// +build wireinject

package main

import (
    pkg "github.com/ahrtolia/goboot/pkg"
    "github.com/ahrtolia/goboot/pkg/config"
    "github.com/ahrtolia/goboot/pkg/gin"
    "github.com/ahrtolia/goboot/pkg/gorm"
    "github.com/ahrtolia/goboot/pkg/logger"
    "github.com/google/wire"
)

var (
    configSet = wire.NewSet(
        config.ProviderSet,
        config.NacosProvider,
    )

    loggerSet = wire.NewSet(
        logger.ProviderSet,
    )

    httpSet = wire.NewSet(
        gin.ProviderSet,
    )

    dbSet = wire.NewSet(
        gorm.ProviderSet,
    )

    appSet = wire.NewSet(
        wire.Struct(new(pkg.App), "*"),
    )

    globalSet = wire.NewSet(
        configSet,
        loggerSet,
        httpSet,
        dbSet,
        appSet,
    )
)

func CreateApp(configFile string) (*pkg.App, error) {
    wire.Build(
        globalSet,
    )
    return nil, nil // wire 会替换
}
`
