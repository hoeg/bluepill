# bluepill

Mutating admission controller for kubernetes that ensures restricted ingress.

![alt text](https://github.com/hoeg/bluepill/blob/main/pills.png?raw=true)

## Configuration

The hook services utilizes a mix of environment variables and files for configuration.

### BLUEPILL_HTTP_PORT

Port that the webhook will listen on, defaults to `8443`.

### TLS

A mutating webhook can only be called on a TLS connection.
Therefore a utility cli has been added to generate a self signed certificate for deploying bluepill.

Run `go run cmd/certificate_generator/main.go bluepill default` to get a secret containing certificate and private key.

### BLUEPILL_HTTP_CERTIFICATE_FILE

Points to the mounted certificate file.

### BLUEPILL_HTTP_KEY_FILE

Points to the mounted private key file.

### BLUEPILL_ENFORCEMENT_WHITELIST_FILE

Points to the mounted whitelist file.

### Whitelist format

```
name1=ip1
name2=ip2
...
nameN=ipN
```

### BLUEPILL_ENFORCEMENT_ENFORCE

Indicates if the we should mutate or just log.

## Deployment

The `deploy/` folder contains all resources to get up and running except the secret containing the self signed certificate.

`./setup-kind` prepares a cluster where the functionality can be tested by running:

```
kubectl apply -f deploy/admission-webhook.yaml && kubectl apply -f deploy/test/ingress-blue.yaml
```

and then inspecting the ingress resource afterwards.