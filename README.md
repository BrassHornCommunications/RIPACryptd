# RIPACryptd
The server side of RIPACrypt - Remotely store secrets that are destroyed if the specified time has passed without a check-in *(a form of deadman switch)*

This is another thought experiment in defeating the Regulation of Investigatory
Powers Act 2000 Section 49 (compelled decryption) legislation using technology.

## Process
1. Create a RIPACrypt account (on the computer in which the disks are protected by full disk encryption)
2. Store the pass phrase for the full disk encryption in a crypt with a 3 hour timeout *(3x 1 hour check-ins)*
3. Configure `cron` to checkin every hour 
4. Time passes...
5. Computer is seized *(Harrassment for running a Tor relay etc etc)*
6. After 3 hours without a successful checkin RIPACryptd deletes the secret from the crypt, locks the crypt and sets `IsDestroyed` to true
7. When issued with a s.49 notice you simply point to https://ripacrypt.download/view/YOURCRYPTID

## Questions
- You no longer possess the key to decrypt the FDE, are you failing to comply? *(RIPA 2000 s.53)*
- You couldn't prevent the destruction of the key because you were no longer in possession of the computer authorised to prevent it's deletion, have you obstructed a Police officer or perverted the course of justice?
- Have the Police 'stolen' your *digital* property by "permanently depriving" you of it? *(Theft Act 1968)*

## Development
- [x] Register
- [x] Create a new crypt
- [x] Checkin with crypt
- [x] View status of a crypt (API)
- [x] View status of a crypt (HTML)
- [x] Generate a challenge
- [x] Generate a new Bitcoin address
- [ ] Delete a crypt
- [ ] Specify a notification method if a checkin period is missed


## Disclaimers
**We are not a lawyers, solictiors or in any way well versed in the law. This software may result in you being found guilty of Failure to comply with a notice under the Regulation of Investigatory Powers Act 2000 and sentenced to up to 5 years in prison (or worse)**

> THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE. 

# API
RIPACrypt provides an open API so that other clients can work with the service.
It can even be used with simple CURL calls *(see ExampleCurlCalls.md)*

### Endpoints
#### `POST` /1/register/
To register an account all that is needed is a GPG public key. This is used as authentication *(and for encrypting secrets)*

**Send:**

```json
{
 "public_key":"armoured public key (escaped line endings)"
}
```

**Receive:**

```json
{
 "status_code":200,
 "success":true,
 "status_message":"Your account has been successfully created",
 "version":1,
 "btc_addr":"",
 "user_id":10
}
```

#### `POST` /1/challenge
Challenges are nonces that are encrypted with your public key, you will need 
to decrypt the challenge and send the nonce back as _plain text_ with the next
request. Successful challenges are deleted so you need to request a new 
challenge before each authenticated request.

**Send:**

```json
{
 "fingerprint": "16 char GPG fingerprint",
 "user_id":your user id
}
```

**Receive:**

```json
{
 "challenge":"wcFMAwAx48WIB6+GARAAFBRUQf+IwA <snip>",
 "challenge_id":22,
 "user_id":9,
 "status_code":200, 
 "success":true,
 "status_message":"challenge Successfully generated!",
 "version":1
}
```

#### `POST` /1/crypt/new/
Submit a secret to be stored specifying the CheckIn duration *(time between checkins)* in seconds and the miss count 

**Send:**

```json
{
 "challenge":"PLAIN TEXT MD5 FROM CHALLENGE",
 "challenge_id":NUMERIC ID OF CHALLENGE,
 "user_id":YOUR ID,
 "crypt_content":"YOUR SECRET (PREFERABLY ENCRYPTED)",
 "description":"Description of your crypt (preferably tying it the computer/secret you are protecting",
 "checkin_duration": 86400,
 "miss_count": 3
}
```

**Receive:**

```json
{
 "status_code":201,
 "success":true,
 "status_message":"Crypt successfully created!",
 "version":1,
 "crypt":
 {
  "crypt_id": "29d48f65cf5899ba6119db5fc37648e5",
  "ciphertext":"W2+6uaow2WwHSitrOw1KmY9yqUykv41cHmdPtMF9pwP4Y9Q/eqOeqOv86+aFFoELgs4=",
  "crypt_timestamp":1459468801,
  "crypt_description":"Description of your crypt (preferably tying it the computer/secret you are protecting",
  "is_crypt_destroyed":false,
  "last_checkin":1459468801,
  "check_in_duration":86400,
  "miss_count":3
 }
}
```

#### `HEAD` /1/crypt/CRYPTID

**Send:**

N/A

**Receive:**

*Note these are HTTP headers*

Live crypt:

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 01 Apr 2016 0:0:10 GMT
```

Destroyed Crypt:

```
HTTP/1.1 410 Gone
Content-Type: application/json
Date: Mon, 01 Apr 2016 0:0:10 GMT
```


#### `GET` /1/crypt/CRYPTID
**Send:**

N/A

**Receive:**

*HTTP headers will match the JSON `status_code`*

Live Crypt:

```json
{
 "status_code":200,
 "success":true,
 "status_message":"Crypt Retrieval Successful!",
 "version":1,
 "crypt":
 {
  "crypt_id":"29d48f65cf5899ba6119db5fc37648e5",
  "ciphertext":"hQIMAwAx48WIB6+ <snip>",
  "crypt_timestamp":1459468801,
  "crypt_description":"Description of your crypt (preferably tying it the computer/secret you are protecting",
  "is_crypt_destroyed":false,
  "last_checkin":1459468801,
  "check_in_duration":86400,
  "miss_count":3
 }
}
```

Destroyed crypt:

```json
{
 "status_code":410,
 "success":true,
 "status_message":"Crypt metadata retrived but crypt was wiped",
 "version":1,
 "crypt":
 {
  "crypt_id":"9b759040321a408a5c7768b4511287a6",
  "ciphertext":"",
  "crypt_timestamp":1459468801,
  "crypt_description":"Description of your crypt (preferably tying it the computer/secret you are protecting",
  "is_crypt_destroyed":true,
  "last_checkin":1459468801,
  "check_in_duration":3600,
  "miss_count":20
 }
}
```

#### `DELETE` /1/crypt/CRYPTID
Delete allows you to destroy a crypt immediately, unfortunately the `last_checkin` 
will still be valid so an observer can extrapolate that the crypt was destroyed 
by an authorised user rather than by the system if they query the crypt before
the requisite amount of time has passed.

**Send:**

*Unspecified*

**Receive:**

*Unspecified*
