/*
 * This file was last modified at 2024-07-08 13:46 by Victor N. Skurikhin.
 * init_config.go
 * $Id$
 */

package env

import (
	"strconv"
	"time"

	"github.com/vskurikhin/gometrics/internal/util"
)

func initAgentConfig() {

	initAgentCryptoKey()
	initGRPCAddress()
	initServerAddress()
	if env.ReportInterval > 0 {
		cfg.reportInterval = time.Duration(env.ReportInterval) * time.Second
	} else {
		setIfFlagChanged(ReportInterval, func() {
			cfg.reportInterval = *flag.reportInterval
		})
		if cfg.reportInterval == 0 {
			reportInterval, err := time.ParseDuration(jsonAgentConfig.ReportInterval)
			util.IfErrorThenPanic(err)
			cfg.reportInterval = reportInterval
		}
	}
	if env.PollInterval > 0 {
		cfg.pollInterval = time.Duration(env.PollInterval) * time.Second
	} else {
		setIfFlagChanged(PollInterval, func() {
			cfg.pollInterval = *flag.pollInterval
		})
		if cfg.pollInterval == 0 {
			pollInterval, err := time.ParseDuration(jsonAgentConfig.PollInterval)
			util.IfErrorThenPanic(err)
			cfg.pollInterval = pollInterval
		}
	}
	initKey()
}

func initServerConfig() {

	initGRPCAddress()
	initServerAddress()
	initServerCryptoKey()
	if env.DataBaseDSN != "" {
		cfg.dataBaseDSN = &env.DataBaseDSN
	} else {
		setIfFlagChanged(DatabaseDSN, func() {
			cfg.dataBaseDSN = flag.dataBaseDSN
		})
		if cfg.dataBaseDSN == nil || *cfg.dataBaseDSN == "" {
			cfg.dataBaseDSN = &jsonServerConfig.DatabaseDSN
		}
	}
	if env.StoreInterval != "" {
		storeInterval, err := strconv.Atoi(env.StoreInterval)
		if err == nil {
			cfg.storeInterval = time.Duration(storeInterval) * time.Second
		} else {
			cfg.storeInterval = 300 * time.Second
		}
	} else {
		setIfFlagChanged(StoreInterval, func() {
			cfg.storeInterval = *flag.storeInterval
		})
		if cfg.storeInterval == 0 {
			storeInterval, err := time.ParseDuration(jsonServerConfig.StoreInterval)
			util.IfErrorThenPanic(err)
			cfg.storeInterval = storeInterval
		}
	}
	if env.FileStoragePath != "" {
		cfg.fileStoragePath = env.FileStoragePath
	} else {
		cfg.fileStoragePath = jsonServerConfig.StoreFile
		setIfFlagChanged(DatabaseDSN, func() {
			cfg.fileStoragePath = *flag.fileStoragePath
		})
	}
	if env.Restore != "" {
		restore, err := strconv.ParseBool(env.Restore)
		if err == nil {
			cfg.restore = restore
		}
	} else {
		cfg.restore = jsonServerConfig.Restore
		setIfFlagChanged(Restore, func() {
			cfg.restore = *flag.restore
		})
	}
	if env.TrustedSubnet != "" {
		cfg.trustedSubnet = env.TrustedSubnet
	} else {
		cfg.trustedSubnet = jsonServerConfig.TrustedSubnet
		setIfFlagChanged(TrustedSubnet, func() {
			cfg.trustedSubnet = *flag.trustedSubnet
		})
	}
	initKey()
}

func initKey() {
	if env.Key != "" {
		cfg.key = &env.Key
	} else {
		cfg.key = flag.key
	}
}

func initAgentCryptoKey() {
	if len(env.CryptoKey) > 1 {
		cfg.cryptoKey = env.CryptoKey
	} else {
		setIfFlagChanged(CryptoKey, func() {
			cfg.cryptoKey = util.Str(flag.cryptoKey)
		})
		if len(cfg.cryptoKey) == 0 {
			cfg.cryptoKey = jsonAgentConfig.CryptoKey
		}
	}
}

func initServerCryptoKey() {
	if len(env.CryptoKey) > 1 {
		cfg.cryptoKey = env.CryptoKey
	} else {
		setIfFlagChanged(CryptoKey, func() {
			cfg.cryptoKey = *flag.cryptoKey
		})
		if len(cfg.cryptoKey) == 0 {
			cfg.cryptoKey = jsonServerConfig.CryptoKey
		}
	}
}

func initGRPCAddress() {
	if len(env.GRPCAddress) > 1 {
		cfg.grpcAddress = env.parseEnvGRPCAddress()
	} else {
		cfg.grpcAddress = jsonConfig.getGRPCAddress()
		setIfFlagChanged(Address, func() {
			cfg.grpcAddress = util.Str(flag.grpcAddress)
		})
	}
}

func initServerAddress() {
	if len(env.Address) > 1 {
		cfg.serverAddress = env.parseEnvAddress()
	} else {
		cfg.serverAddress = jsonConfig.getAddress()
		setIfFlagChanged(Address, func() {
			cfg.serverAddress = util.Str(flag.serverAddress)
		})
	}
}
