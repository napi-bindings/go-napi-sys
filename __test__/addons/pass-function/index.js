'use strict'

const addon = require('bindings')('addon');

addon(function(msg){
    console.log(msg);
});