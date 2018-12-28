# kmstool

This is a simple little util to encrypt/decrypt a file using a symmetric KMS
key, that can read/write to a mounted file path, or a GCS bucket object.

## Usage

``` shell
kmstool [decrypt | encrypt] \
    --ciphertext PATH \
    --plaintext PATH \
    --key KEY
    
where PATH is either a local file, or a GCS object specified as
gs://bucket/name, and KEY is a fully qualified identifier for a unique KMS key
in the form projects/MYPROJECTID/locations/LOCATION/keyRings/KEYRING/KEYNAME
```

## Why?

This could be done with ```gcloud``` shell commands, or using
```gcr.io/cloud-builders/gcloud```, but those are trickier to make work reliably
as a systemd unit in *cloud-config*. At Neudesic We use this tool to prepare 
[COS](https://cloud.google.com/container-optimized-os/) VMs with secrets that
have been encrypted with a KMS key.

E.g. for [Atum](https://github.com/NeudesicGCP/atum), the *user-data* entry
looks something like this:-

``` yaml
#cloud-config

# Create Atum user account with fixed UID
users:
- name: atum
  uid: 2000
  groups: docker
  lock_passwd: true
  homedir: /var/lib/atum

write_files:
- path: /etc/systemd/system/terraform-credentials.service
  permissions: 0644
  owner: root
  content: |
    [Unit]
    Description=Extract terraform service account credentials
    Wants=gcr-online.target
    After=gcr-online.target

    [Service]
    User=atum
    ExecStartPre=/usr/bin/docker-credential-gcr configure-docker
    ExecStartPre=/usr/bin/docker pull gcr.io/neudesicgcp/kmstool:latest
    ExecStart=/usr/bin/docker run \
        --rm \
        --name kmstool \
        --volume /var/lib/atum:/var/lib/atum \
        --user 2000 \
        --log-driver gcplogs \
        gcr.io/neudesicgcp/kmstool:latest \
        decrypt \
        --key "${google_kms_crypto_key.terraform_robot.self_link}" \
        --ciphertext "${local.terraform_robot_creds_gcs}" \
        --plaintext "/var/lib/atum/terraform-credentials.json"
    ExecStop=-/usr/bin/docker stop kmstool
    ExecStopPost=-/usr/bin/docker rm kmstool
    Restart=on-abnormal
...
runcmd:
- systemctl daemon-reload
- systemctl start atum.service
```

The *terraform-credentials.service* is invoked as a dependency of *atum.service*
to read an encrypted file from GCS, and decrypt it to a local file, using a
specified KMS key. When the container is finished executing, the file
`/var/lib/atum/terraform-credentials.json` will be available for use by other
containers.

## Pre-built container

The latest build from this repo can be pulled from gcr.io/neudesicgcp/kmstool

```shell
docker run --rm --name kmstool gcr.io/neudesicgcp/kmstool:latest ...
```

## Contributing

Contributions are welcome; open an issue, fork the repo and submit a PR!
