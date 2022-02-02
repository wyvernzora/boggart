# Boggart
**Boggart** is a simple AWS identity broker that allows microservices to exchange their 
SSH identities for AWS temporary credentials.

## Why Boggart?
Boggart is here to solve problems related to managing AWS account credentials across 
multiple microservices and/or hosts and/or containers.

 - Instead of each microservice needing its own AWS user, there is only one user for Boggart
 - Only one long-lived AWS credential to store and secure; microservice use short-lived credentials
 - Microservices can only access AWS roles with minimum privilege, reducing impact of compromise

## Usage
```shell
$ docker run -d -v /etc/boggart/:/etc/boggart/ -p 2222:2222 boggart:latest
```

Use SSH to retrieve credentials:
```shell
$ ssh -p 2222 boggart '{"roleArn": "arn:aws:iam::123456789012:role/TestRole", "sessionName": "TestSession"}'
```

Where Boggart responds with the following on success:
```json
{
  "success": true,
  "credentials": {
    "accessKeyId": "<access-key-id>",
    "secretAccessKey": "<secret-access-key>",
    "sessionToken": "<session-token>",
    "expiresAt": "<expiration-date-iso8601>"
  }
}
```

On error, the response may look like the following:
```json
{
  "success": false,
  "error": "not allowed to assume arn:aws:iam::123456789012:role/TestRole"
}
```

## Configuration
Boggart looks for configuration file at `/etc/boggart/config.yml`.
See `config.example.yml` for more details on config options.

## Host Key
Even though Boggart can be run without specifying a host key, not having one causes Boggart to
re-generate one on every restart. This causes issues with host key pinning (aka `known_hosts`).

Generate host key using the following:
```shell
$ ssh-keygen -t ed25519 -N '' -f /etc/boggart/host_key 
```

Then configure Boggart to use it
```yaml
HostKey: /etc/boggart/host_key
```

## License
MIT