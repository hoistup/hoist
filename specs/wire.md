# Wire

Wire is the format for transporting data "on the wire" in Hoist.
It is targeted at running functions on microservices, and the data associated with those functions.

## Format

Wire uses the following format to transport data:

| Protocol ID | , | Payload 1 Size | , | Payload 2 Size | , | ... | : | Payload 1 | Payload 2 | ... |
|-------------|---|----------------|---|----------------|---|-----|---|-----------|-----------|-----|

Where:
  - `Payload ID` is a number corresponding to the format of the following data
  - `Payload # Size` is the size in bytes of `Payload #`
  - `Payload #` contains the data for that payload stored in the format specified by the `Payload ID`

Currently the only valid `Payload ID` is `1`, which indicates JSON encoding.
In the future, we may support other types of encoding, or even an encoding + encryption.

Usually `Payload 1` indicates the service and function to call, and `Payload 2` contains the parameters for that function. This separation of concerns allows `Payload 2` to not be unmarshalled until it reaches the function. In a sense, `Payload 1` contains the "routing information".

## Example

```
1,26,20:{"svc":"echo","fn":"echo"}{"msg":"What's up?"}
```

Where:
  - `1` indicates we are using JSON encoding
  - `26` indicates `Payload 1` has 26 bytes
  - `20` indicates `Payload 2` has 20 bytes
  - `{"svc":"echo","fn":"echo"}` is `Payload 1` encoded in JSON
  - `{"msg":"What's up?"}` is `Payload 2` encoded in JSON
