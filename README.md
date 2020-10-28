# kimg
# kimg

Simple program to extract OCI or docker images from a well-formed kubernetes manifest file

```sh
#example
curl -sL 'https://raw.githubusercontent.com/kubernetes-csi/csi-driver-smb/master/deploy/v0.4.0/csi-smb-controller.yaml' | kimg
```