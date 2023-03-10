module github.com/bhbosman/goCommonMarketData

go 1.18

require (
	github.com/bhbosman/goCommsDefinitions v0.0.0-20230308211643-24fe88b2378e
	github.com/bhbosman/goMessages v0.0.0-20230302063433-258339efe599
	github.com/bhbosman/goUi v0.0.0-20230302065227-24c3cb06165e
	github.com/bhbosman/gocommon v0.0.0-20230302065226-cf91cedfab33
	github.com/bhbosman/goerrors v0.0.0-20220623084908-4d7bbcd178cf
	github.com/bhbosman/goprotoextra v0.0.1
	github.com/cskr/pubsub v1.0.2
	github.com/emirpasic/gods v1.18.1
	github.com/gdamore/tcell/v2 v2.5.1
	github.com/reactivex/rxgo/v2 v2.5.0
	github.com/rivo/tview v0.0.0-20220709181631-73bf2902b59a
	go.uber.org/fx v1.19.2
	go.uber.org/zap v1.23.0
	golang.org/x/net v0.0.0-20211015210444-4f30a5c0130f
	google.golang.org/protobuf v1.28.0
)

require (
	github.com/bhbosman/gomessageblock v0.0.0-20220617132215-32f430d7de62 // indirect
	github.com/cenkalti/backoff/v4 v4.1.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gdamore/encoding v1.0.0 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/stretchr/objx v0.4.0 // indirect
	github.com/stretchr/testify v1.8.0 // indirect
	github.com/teivah/onecontext v0.0.0-20200513185103-40f981bfd775 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/dig v1.16.1 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/sys v0.0.0-20220318055525-2edf467146b5 // indirect
	golang.org/x/term v0.0.0-20210220032956-6a3ed077a48d // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/gdamore/tcell/v2 => github.com/bhbosman/tcell/v2 v2.5.2-0.20220624055704-f9a9454fab5b

replace github.com/golang/mock => github.com/bhbosman/gomock v1.6.1-0.20230302060806-d02c40b7514e

replace github.com/rivo/tview => ../tview

replace github.com/bhbosman/gocommon => ../gocommon

replace github.com/bhbosman/goUi => ../goUi

replace github.com/bhbosman/goerrors => ../goerrors

replace github.com/bhbosman/goprotoextra => ../goprotoextra

replace github.com/bhbosman/goMessages => ../goMessages

replace github.com/cskr/pubsub => github.com/bhbosman/pubsub v1.0.3-0.20220802200819-029949e8a8af

replace github.com/bhbosman/goCommsDefinitions => ../goCommsDefinitions
