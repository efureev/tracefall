[![Build Status](https://travis-ci.org/efureev/traceFall.svg?branch=master)](https://travis-ci.org/efureev/traceFall)
[![Version](https://img.shields.io/badge/version-1.0.1-blueviolet.svg)](https://travis-ci.org/efureev/traceFall)
[![Maintainability](https://api.codeclimate.com/v1/badges/c933f06740177611ff5a/maintainability)](https://codeclimate.com/github/efureev/traceFall/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/c933f06740177611ff5a/test_coverage)](https://codeclimate.com/github/efureev/traceFall/test_coverage)
[![Go Report Card](https://goreportcard.com/badge/github.com/efureev/traceFall)](https://goreportcard.com/report/github.com/efureev/traceFall)
[![codecov](https://codecov.io/gh/efureev/traceFall/branch/master/graph/badge.svg)](https://codecov.io/gh/efureev/traceFall)

## Info
Package for sending logs to the storage, for the subsequent withdrawal of the traceViewer service and display there.

Supported storage drivers:  
- [x] Console // invalid realisation
- [x] Postgres // invalid realisation
- [ ] Algolia
- [ ] ElasticSearch 

## Content
- Thread Line: Line of logs. Contains Logs. Thread ID = First root Log ID 
- Log: data node. May contents other Logs as children

## Use

**Create new Log node**
```go
import "github.com/efureev/traceFall"
// ...
log := traceFall.NewLog(`test log`)
```

**Finish log**
```go
log := traceFall.NewLog(`test log`)

// with fail result
log.Fail(err error)

// with success result 
log.Success()

// without result: set finish time of the log
log.FinishTimeEnd()

```
**Finish thred of logs**
```go
log.ThreadFinish()
```

**Add extra data to Log**
```go
log := traceFall.NewLog(`test log`)
log.Data.Set(`url`, `http://google.com`).Set(`service`, service.Name)
```

**Add notes to Log**
```go
log.Notes.Add(`send to redis`, `ok`).Add(`send to rabbit`, `ok`)
//or
log.Notes.AddGroup(`send to redis`, [`ping`,`processing`,`done`])
```

