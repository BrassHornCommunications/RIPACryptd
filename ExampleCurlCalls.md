# Register

Registering is a single call that only requires your GPG public key _(new lines escaped please)_

You will receive back a `user_id` and a `btc_addr`. The `user_id` identifies you as a user, the `btc_addr` can be used for donations or _(in future versions)_ to pre-pay for larger amounts of storage

```
KEY=`gpg --export | base64 -w0`
curl -v -H "Accept: application/json" -X POST -d "{\"public_key\":\"-----BEGIN PGP PUBLIC KEY BLOCK-----\\n\\n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\\n=4UbB\\n-----END PGP PUBLIC KEY BLOCK-----\\n\"}" https://ripacrypt.download/1/register/
```

# Authentication
Most calls require some form of authentication, since the application relies heavily on GPG to provide encryption for the secrets you upload we also use GPG for authentication.

Unfortunately GoLangs GPG signing implementation is broken so instead we use an encrypted challenge/response system.

## Get a Challenge
Challenges are nonces that are encrypted with your public key, you will need to decrypt the challenge and send the nonce back as _plain text_ with the next request. Successful challenges are deleted so you need to request a new challenge before each authenticated request.

```bash 
curl -v -H "Accept: application/json" -X POST -d "{\"fingerprint\":\"447055F36F6D60C1\",\"user_id\":9}" https://ripacrypt.download/1/challenge/
```

# Managing Secrets
## Store something in a Crypt
Secrets are the thing you want to 'lose' if you fail to checkin, if using the RIPACrypt client they will be encrypted with your GPG public key but you can store whatever you want

Everyone gets 2Kb of storage for free (enough for a ~1450 character password encrypted with a 4096 bit key piped through `base64 -w0`

```bash
SECRET=`pwgen 1450 1 | gpg --encrypt --recipient hello@brasshorncommunications.uk | base64 -w0`
curl -v -H "Accept: application/json" -X POST -d "{\"challenge\":\"9b759040321a408a5c7768b4511287a6\", \"challenge_id\":18,\"user_id\":9, \"crypt_content\":\"$SECRET\"}" https://ripacrypt.download/1/crypt/new/
```

## Retrieve a Crypt
GET requests do not require any authentication *(which is why we recommend that the 'secret' be encrypted)* and will return a chunk of JSON containing meta-data and the secret.

```bash
curl https://ripacrypt.download/1/crypt/3317f65eed845697115d962a572d11e3/
```

HEAD requests also do not require any authentication and will return a HTTP status code of `200` if all is OK or `410` if the crypt has been destroyed. No body is returned so this is useful for monitoring the status of a crypt where the secret is quite large.

```bash
curl -I https://ripacrypt.download/1/crypt/3317f65eed845697115d962a572d11e3/
```


## Delete the Secret

## Check-in
Check-ins are what keep your crypts 'alive'. If you have set your crypt to self-destruct if three hours *(Check in duration: 3600, MissCount: 3)* pass without a check-in then you will need to checkin at least every 2 hours and 50 minutes *(the self-destruct watchdog runs every ~5 minutes)* otherwise the secret will be destroyed and the crypt will be locked.

```bash
curl -v -H "Accept: application/json" -X POST -d "{\"challenge\":\"5f185b74b96edab86b0c29f26760c433\",\"user_id\":9, \"challenge_id\":22}" https://ripacrypt.download/1/crypt/a78d8016131d6c73539d464fef2be8b8/
```
