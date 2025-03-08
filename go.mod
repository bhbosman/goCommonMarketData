module github.com/bhbosman/goCommonMarketData

go 1.24.0

require (
	github.com/bhbosman/goCommsDefinitions v0.0.0-20250308000247-4306925b3dfd
	github.com/bhbosman/goMessages v0.0.0-20250307224348-83ddb4c19467
	github.com/bhbosman/goUi v0.0.0-20250308052840-a0e5fd7e5f88
	github.com/bhbosman/gocommon v0.0.0-20250307235859-f370cb0a3bac
	github.com/bhbosman/goerrors v0.0.0-20250307194237-312d070c8e38
	github.com/bhbosman/goprotoextra v0.0.2
	github.com/cskr/pubsub v1.0.2
	github.com/emirpasic/gods v1.18.1
	github.com/gdamore/tcell/v2 v2.8.1
	github.com/golang/mock v1.6.0
	github.com/reactivex/rxgo/v2 v2.5.0
	github.com/rivo/tview v0.0.0-20241227133733-17b7edb88c57
	github.com/stretchr/testify v1.10.0
	go.uber.org/fx v1.23.0
	go.uber.org/zap v1.27.0
	golang.org/x/net v0.37.0
	google.golang.org/protobuf v1.28.0
)

require (
	github.com/bhbosman/gomessageblock v0.0.0-20250307141417-ab783e8e2eba // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gdamore/encoding v1.0.1 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/teivah/onecontext v1.3.0 // indirect
	go.uber.org/dig v1.18.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/term v0.30.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/gdamore/tcell/v2 => github.com/bhbosman/tcell/v2 v2.5.2-0.20220624055704-f9a9454fab5b

replace github.com/golang/mock => github.com/bhbosman/gomock v1.6.1-0.20250308071159-4cf72f668c72

replace github.com/cskr/pubsub => github.com/bhbosman/pubsub v1.0.3-0.20220802200819-029949e8a8af

replace github.com/rivo/tview => github.com/bhbosman/tview v0.0.0-20250308051327-a656c1bc9cfa

//replace github.com/bhbosman/goerrors => ../goerrors

//

//replace github.com/bhbosman/goMessages => ../goMessages

//replace github.com/bhbosman/goCommsDefinitions => ../goCommsDefinitions
