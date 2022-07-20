package main

const sextetBitMask byte = 1<<6 - 1 // 0b00111111
const capitalAOffset = 65
const lowercaseAOffset = capitalAOffset | 1<<5 // (97)
const zeroCharOffset = 48
