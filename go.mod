module github.com/smartcontractkit/external-initiator

go 1.15

require (
	github.com/Conflux-Chain/go-conflux-sdk v1.0.4
	github.com/Depado/ginprom v1.2.1-0.20200115153638-53bbba851bd8
	github.com/avast/retry-go v2.6.0+incompatible
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869
	github.com/centrifuge/go-substrate-rpc-client v2.0.0+incompatible
	github.com/centrifuge/go-substrate-rpc-client/v3 v3.0.0
	github.com/ethereum/go-ethereum v1.10.2
	github.com/facebookgo/clock v0.0.0-20150410010913-600d898af40a
	github.com/gin-gonic/gin v1.6.3
	github.com/golang/mock v1.6.0
	github.com/google/uuid v1.1.5
	github.com/gorilla/websocket v1.4.2
	github.com/iotexproject/iotex-proto v0.4.3
	github.com/jinzhu/gorm v1.9.16
	github.com/magiconair/properties v1.8.5
	github.com/mattn/go-sqlite3 v2.0.1+incompatible // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.11.0
	github.com/shopspring/decimal v1.2.0
	github.com/smartcontractkit/chainlink v0.9.5-0.20201214122441-66aaea171293
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.8.0
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.11
	github.com/terra-money/core v0.5.0-rc0
	github.com/tidwall/gjson v1.6.7
	go.uber.org/atomic v1.7.0
	go.uber.org/zap v1.17.0
	google.golang.org/grpc v1.39.0
	gopkg.in/gormigrate.v1 v1.6.0
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
