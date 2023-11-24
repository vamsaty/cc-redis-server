## File Structure

### executor
Defines the supported commands. Executes the commands in the datastore and generates a response.

### server
Contains the code for the server. Starts a listener (at 6379 port) and connection handler (concurrent).

### tokenizer
Contains the code for the tokenizer. This is used by the server to parse the commands sent by the client.

### types
Defines the RESP datatypes and helper function to convert RESP types to string & vice-versa.
