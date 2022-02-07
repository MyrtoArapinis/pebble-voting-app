Pebble voting app
=================

Pebble is a work-in-progress implementation of the [E-cclesia](https://eprint.iacr.org/2020/513.pdf) voting protocol.

Implementation progress:

- [x] Serialization of data structures used by the protocol, including the eligibility list, anonymous credential announcement, encrypted and signed ballots, and timelock encryption solutions
- [x] Timelock encryption using [Pietrzak's VDF](https://eprint.iacr.org/2018/627.pdf)
- [x] Signing and verification using non-anonymous keys used for the eligibility list
- [x] Signing and verification using anonymous credentials
- [ ] Interface to the Tezos blockchain
- [ ] Protocol logic
