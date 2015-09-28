## Go HTTPS Server and Client with Mutual Authentication
A skeleton HTTPS client and server that supports: no authentication (HTTP), server-side authentication (standard HTTPS), mutual authentication.

Used to test [symbios](https://github.com/dnascimento/symbios)

Docker file included

## Usage
### No Authentication
```
http-client  <endpoint>
http-server  <port>
```

### Server Side Authentication
```
http-client  <endpoint> <ca.pem>
http-server  <port> <key.pem> <cert.pem>
```

### Mutual Authentication
```
http-client  <endpoint> <key.pem> <cert.pem> <ca.pem> 
http-server  <port> <key.pem> <cert.pem> <ca.pem>

```

### Contributors
[Dário Nascimento](mailto:dfrnascimento@gmail.com)
