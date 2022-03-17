# Using Pebble

*alpha version*

## Starting the mock server

Running `./server` from your terminal will start the mock server, with it listening on port 8090. You can specify an alternative port by providing it as the first (and only) argument.

## Configuring the command line app

It is recommended to rename the  compiled `jar` to a simpler name e.g. `pebble.jar`.

Create a text file named `server.txt` within the directory where the `jar` is located. The content should be the address of the server the app should use, e.g. `http://localhost:8090`.

## Using the command line app

Run the app with

    java -jar pebble.jar

Without providing any arguments the app will display the list of available commands.

To generate a (non-anonymous) key pair, type

    java -jar pebble.jar pubkey

The app will display your public key.

To create an election, type

    java -jar pebble.jar election create

You will be prompted to fill the details of the election, and at the end you will be presented with its ID.

To interact with a particular election you need to first set the current election ID. To do that, type

    java -jar pebble.jar election id ID

replacing `ID` with the election's ID.

Consult the list of available commands for actions that can be taken.
