# Election phase
## Voting process

<a><img src="http://www.pebble.vote/images/diagrams/4Pebble_voting_phase.png" alt="Pebble User notification" width="600"></a>

Each voter encrypts their vote with time-lock encryption using the parameters received during the credential generation phase, for instance with a repeated squaring based puzzle, then signs it using the private credential (serial number) generated in the previous phase.Voters can post their encrypted votes to the smart contract anytime before the end of the phase. The difficulty of the time-lock encryption is automatically adapted depending on the remaining time of the election, to both guarantee election fairness and fast tallying.

## Self-tallying

<a><img src="http://www.pebble.vote/images/diagrams/5Pebble_tallying.png" alt="Pebble User notification" width="600"></a>

At the end of the voting phase, voters are invited to take part in the public, decentralized tallying process, where anyone can break the time-lock of encrypted votes on the smart contract one by one, posting the results alongside the corresponding ZKP credential -to prevent talliers to open the same ballot twice- publicly.

To save time and energy, every voter is invited to self-reveal her own vote and anonymously post the time-lock encryption puzzle solution to the smart contract, as if they successfully resolved it. This task can be performed as a background task by the client app.

The tallying process can start as soon as a ballot is posted on the smart contract, ensuring near-instantaneous results at the end of the election, or can be performed overtime. At the end of this phase, all of the decrypted votes are available on the smart contract. Each voter receives the results and anyone can independently calculate the outcome of the election and verify the authenticity of the results.

## Blockchain interface

<a><img src="http://www.pebble.vote/images/diagrams/5Pebble_tallying.png" alt="Pebble User notification" width="600"></a>

Pebble can be used by mainstream voters who would not have a pair of keys on the blockchain.

To maintain anonymity, client apps must use an anonymization network like Tor or I2P.

Since asking every voter to pay network fees would make the UX cumbersome and could threaten voters’ privacy through fund tracing, we suggest a semi-centralized setup. In the common scenario client apps will use semi-centralized servers with provisioned wallets to make blockchain transactions on behalf of voters. These servers are not trusted-third parties since the validity of messages on the smart contract is still determined by the use of client credentials, which means a server cannot spoof messages or ballots. If a server refuses or is unable to make a valid blockchain transaction, the client app will be able to quickly detect it and try another server or prompt the user to make a direct blockchain transaction.