'use strict'

import {SummaryView} from './summary'
import React from 'react'
import ReactDOM from 'react-dom'

let outputEl = document.querySelector("#output")
let el = (<SummaryView />)
ReactDOM.render(el, outputEl, function () {
 this.load(outputEl.getAttribute('data-url'))
})
