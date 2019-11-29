'use strict'

const addon = require('bindings')('addon');

let iterations = 23

const unix = () => { 
    const t = addon.unix()
    console.log(`Unix timestamp in nanoseconds => ${t}`)
    iterations -= 1
    if (iterations == 0) {
        clearInterval(interval)
    }
}

const interval = setInterval(unix, 1000)

