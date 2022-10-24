![Pebble logo](doc/Logo-Black.svg)

# Pebble

*The first fully decentralized, secure, and transparent e-voting solution.*

---

## What is Pebble?

Pebble is an open source decentralized, secure, and transparent online voting solution that combines time lock encryption, zero-knowledge proofs, and blockchain technology to enable self-tallying elections.

By empowering every user to act as their own election authority, Pebble eliminates fees and vulnerabilities induced by intermediaries and trusted-third parties.

## Current version

Pebble is currently released as an open early-alpha.

This is a work in progress project. Feel free to review, contribute, and get in touch.

For more details, check out the [documentation section](doc/) and the [Pebble website](https://www.pebble.vote/).

## Implementation progress

- [x] Serialization of data structures used by the protocol, including the eligibility list, anonymous credential announcement, encrypted and signed ballots, and timelock encryption solutions
- [x] Timelock encryption using [Pietrzak's VDF](https://eprint.iacr.org/2018/627.pdf)
- [x] Signing and verification using non-anonymous keys used for the eligibility list (Ed25519, Tezos)
- [x] Signing and verification using anonymous credentials
- [x] Protocol logic
- [x] Mock server
- [x] Election setup service
- [ ] Interface to the Tezos blockchain (In Progress)

## Background

Pebble is a work-in-progress implementation of the [E-cclesia](https://eprint.iacr.org/2020/513.pdf) decentralized and self-tallying voting protocol.

The development of Pebble is made possible thanks to the support of the [Tezos Foundation](https://tezos.foundation/).
