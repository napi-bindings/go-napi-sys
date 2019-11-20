'use strict'

const addon = require('bindings')('addon');

console.log(addon.echo('Please give me back ...'));