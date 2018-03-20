# spin-helm-template

A REST API inspired [helm-template](https://github.com/technosophos/helm-template). The goal of this API is to accept some parameters about a Helm chart, render it, and return the rendered template in Spinnaker Artifact format.

## NOT FOR PRODUCTION USAGE

This is just a toy. It's not done. It's just an experiment to see how we might use Helm + Spinnaker for manifest management + deployment.


## Installation (Your Mileage May Vary)

```bash
$ go install github.com/ethanfrogers/spin-helm-template
$ spin-helm-template // server will be running on port 3005
```

## Usage

*Note: this only serves as an example and uses the default chart values. no custom values can be passed in*

This application depends on `HELM_HOME` being set, just like the Helm client. Running `helm init --client-only` and setting `HELM_HOME` to the path that it says will work.

The `/template` input accepts a `POST` with some parameters about the chart you want to template. It will return a JSON response with artifacts that can be used by Spinnaker.

Let's say I wanted to deploy the `stable/jenkins@0.13.2` chart. I would call this API like so (or using a Webhook stage in Spinnaker)

```bash
curl -X POST \
  http://localhost:3005/template \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: application/json' \
  -d '{
	"chart": "stable/jenkins",
	"version": "0.13.2",
	"releaseName": "my-jenkins",
	"namespace": "default"
}'
```

The output of which would be:

```bash
{
  "artifacts": [
    {
      "type": "embedded/base64",
      "reference": "MS4gR2V0IHlvdXIgJ2FkbWluJyB1c2VyIHBhc3N3b3JkIGJ5IHJ1bm5pbmc6CiAgcHJpbnRmICQoa3ViZWN0bCBnZXQgc2VjcmV0IC0tbmFtZXNwYWNlIGRlZmF1bHQgbXktamVua2lucy1qZW5raW5zIC1vIGpzb25wYXRoPSJ7LmRhdGEuamVua2lucy1hZG1pbi1wYXNzd29yZH0iIHwgYmFzZTY0IC0tZGVjb2RlKTtlY2hvCjIuIEdldCB0aGUgSmVua2lucyBVUkwgdG8gdmlzaXQgYnkgcnVubmluZyB0aGVzZSBjb21tYW5kcyBpbiB0aGUgc2FtZSBzaGVsbDoKICBOT1RFOiBJdCBtYXkgdGFrZSBhIGZldyBtaW51dGVzIGZvciB0aGUgTG9hZEJhbGFuY2VyIElQIHRvIGJlIGF2YWlsYWJsZS4KICAgICAgICBZb3UgY2FuIHdhdGNoIHRoZSBzdGF0dXMgb2YgYnkgcnVubmluZyAna3ViZWN0bCBnZXQgc3ZjIC0tbmFtZXNwYWNlIGRlZmF1bHQgLXcgbXktamVua2lucy1qZW5raW5zJwogIGV4cG9ydCBTRVJWSUNFX0lQPSQoa3ViZWN0bCBnZXQgc3ZjIC0tbmFtZXNwYWNlIGRlZmF1bHQgbXktamVua2lucy1qZW5raW5zIC0tdGVtcGxhdGUgInt7IHJhbmdlIChpbmRleCAuc3RhdHVzLmxvYWRCYWxhbmNlci5pbmdyZXNzIDApIH19e3sgLiB9fXt7IGVuZCB9fSIpCiAgZWNobyBodHRwOi8vJFNFUlZJQ0VfSVA6ODA4MC9sb2dpbgoKMy4gTG9naW4gd2l0aCB0aGUgcGFzc3dvcmQgZnJvbSBzdGVwIDEgYW5kIHRoZSB1c2VybmFtZTogYWRtaW4KCkZvciBtb3JlIGluZm9ybWF0aW9uIG9uIHJ1bm5pbmcgSmVua2lucyBvbiBLdWJlcm5ldGVzLCB2aXNpdDoKaHR0cHM6Ly9jbG91ZC5nb29nbGUuY29tL3NvbHV0aW9ucy9qZW5raW5zLW9uLWNvbnRhaW5lci1lbmdpbmUK",
      "name": "jenkins/templates/NOTES.txt"
    },
    {
      "type": "embedded/base64",
      "reference": "YXBpVmVyc2lvbjogdjEKa2luZDogQ29uZmlnTWFwCm1ldGFkYXRhOgogIG5hbWU6IG15LWplbmtpbnMtamVua2lucy10ZXN0cwpkYXRhOgogIHJ1bi5zaDogfC0KICAgIEB0ZXN0ICJUZXN0aW5nIEplbmtpbnMgVUkgaXMgYWNjZXNzaWJsZSIgewogICAgICBjdXJsIC0tcmV0cnkgMjQgLS1yZXRyeS1kZWxheSAxMCBteS1qZW5raW5zLWplbmtpbnM6ODA4MC9sb2dpbgogICAgfQo=",
      "name": "jenkins/templates/test-config.yaml"
    },
    ...
]}
```
