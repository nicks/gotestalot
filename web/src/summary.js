'use strict'

import React from 'react'
import {List, Map} from 'immutable'

class SummaryView extends React.Component {
  constructor(props) {
    super(props)

    this.state = {
      pkgOrder: List(),
      pkgSpecs: Map(),
      done: false,
    }
  }

  load(path) {
    let url = `ws://${window.location.host}${path}`
    let socket = new WebSocket(url)
    socket.addEventListener('message', this.onmessage.bind(this))
    socket.addEventListener('close', this.onclose.bind(this))
  }

  onclose(event) {
    this.setState((oldState) => {
      return {
        pkgOrder: oldState.pkgOrder,
        pkgSpecs: oldState.pkgSpecs,
        done: true
      }
    })
  }

  onmessage(event) {
    let data = JSON.parse(event.data)
    this.setState((oldState) => {
      let pkg = data.Package
      let pkgOrder = oldState.pkgOrder
      let pkgSpecs = oldState.pkgSpecs
      if (!pkg) {
        let action = data.Action
        if (action == 'error' || action == 'done') {
          return {pkgOrder, pkgSpecs, done: true}
        }
        return
      }

      if (!oldState.pkgSpecs.has(pkg)) {
        pkgOrder = pkgOrder.push(pkg)
        pkgSpecs = pkgSpecs.set(pkg, List([data]))
      } else {
        pkgSpecs = pkgSpecs.update(pkg, (val) => val.push(data))
      }
      return {pkgOrder, pkgSpecs, done: oldState.done}
    })
  }

  render() {
    let pass = 0
    let fail = 0
    let finished = this.state.done

    let elems = this.state.pkgOrder.map((key) => {
      let list = this.state.pkgSpecs.get(key)
      list.forEach((data) => {
        if (data.Action == 'pass') {
          pass++
        }
        if (data.Action == 'fail') {
          fail++
        }
      })

      return <PackageView key={key} pkg={key} list={list} />
    })

    return (<div>
      <TestCountView key="count" finished={finished} pass={pass} fail={fail} />
      {elems}
      <SummaryEndView key="end" finished={finished} fail={fail} />
    </div>)
  }
}

class TestCountView extends React.Component {
  render() {
    let total = this.props.pass + this.props.fail
    let className = "testCountView"
    if (this.props.fail) {
      className += " text-danger"
    }
    let text = ' (and counting)'
    if (this.props.finished) {
      text = ' (done)'
    }
    return (<div className={className}>
            Tests run: {total}{text}<br />
            Passed: {this.props.pass}<br />
            Failed: {this.props.fail}
            </div>)
  }
}

class SummaryEndView extends React.Component {
  render() {
    let text = 'Done'
    let className = "summaryEndView"
    if (this.props.fail) {
      className += " text-danger"
      text += ' (failed)'
    }

    if (!this.props.finished) {
      className += " d-none"
    }

    return (<div className={className}>{text}</div>)
  }
}

class PackageView extends React.Component {
  render() {
    let pkg = this.props.pkg
    let skipped = this.props.list.some((data) => data.Action == 'skip' && !data.Test)
    if (skipped) {
      return (<div className="packageView packageView--skipped">
        <div className="packageSummary">SKIPPED: {pkg}</div>
      </div>)
    }

    let passed = this.props.list.some((data) => data.Action == 'pass' && !data.Test)
    if (passed) {
      return (<div className="packageView packageView--passed text-success">
        <div className="packageSummary">PASSED: {pkg}</div>
     </div>)
    }

    let failed = this.props.list.some((data) => data.Action == 'fail' && !data.Test)
    if (failed) {
      return (<div className="packageView packageView--failed text-danger">
        <div className="packageSummary">FAILED: {pkg}</div>
      </div>)
    }

    return (<div className="packageView packageView--pending">
      <div className="packageSummary">PENDING: {pkg}</div>
    </div>)
  }
}

module.exports = {
  SummaryView,
}
