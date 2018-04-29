# Overview
This is an *experiment* to create a [Kubernetes KMS provider](https://kubernetes.io/docs/tasks/administer-cluster/kms-provider/).  The goal is to provide an implementation of the K8 KMS specification using AWS KMS.

The CLI also comes with a very simple Client CLI to test the Client/Server/AWS interaction.

This should *NOT* be used for production, it is just an attempt to learn several new technologies and better understand this interaction.

# Installation

## Local Installation

**Installation:**
`go get github.com/jrnt30/k8-kms-enc-provider`

**Testing:**
```
# server
k8-kms-enc-provider server --key-id <ARN TO YOUR KEY> --region <AWS REGION OF KEY>

# client encrypt
k8-kms-enc-provider client encrypt --plain-text=test1234

# client decrypt
k8-kms-enc-provider client decrypt --cipher-text=<OUTPUT FROM ENCRYPT>

# client roundtrip
k8-kms-enc-provider client encrypt --plain-text=test1234 | xargs k8-kms-enc-provider client decrypt --cipher-text
```

# Cluster Installation

**NOTE:** The KMS API is in Alpha in K8 1.10 and is sure to change.  During the testing of this I noted several differences with the cluster I had running on 1.9, so please consult the [Official K8 KMS Documentation](https://kubernetes.io/docs/tasks/administer-cluster/kms-provider/#encrypting-your-data-with-the-kms-provider)

My process was to:

- Add a [KMS Encryption Configuration](examples/encryption.conf) to the master node
- Create a [static Pod specification](examples/kms-server.yaml) for the KMS server, copy it to the master's static pod manifests folder
- Adjust the API Server Specification:
  - Add the `- --experimental-encryption-provider-config=/etc/kubernetes/kms/encryption.conf`
  - Add an additional mount to the API server for the shared socket
  - Restart the API server

Deployment notes:

- Using the static pods was easiest to ensure the KMS sidecar was bootstrapped, however it makes debugging slightly difficult
- Initially tried adding this as a sidecar to the apiserver