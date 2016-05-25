# Mongobucks

A slack bot clone of redit's tip bot. It also comes with a web page for viewing balances and such.

## Usage

To invite mongobucks bot into your slack channel, enter:
`/invite mongobucks`

To use within the channel:

```
mongobucks: balance
  # Respond to use with current balance
mongobucks: give @<username> 10 for being a rockstar
  # Transfer 10 mongobucks from sender to <username>
```

## Description

Every user gets 100 mongobucks to start out with. If someone does something nice, feel free to give them mongobucks. 

Mongobucks have no real value.

## Building

Secrets are managed via an environment file.

```
$ echo """
GOOGLE_OAUTH_CLIENT_ID=secret
GOOGLE_OAUTH_SECRET=secret
SLACK_API_TOKEN=secret
SLACK_USERNAME=mongobucks
MONGO_URI=secret""" > .ENV
```

You'll need a pair of google oauth credentials (they can be generated on the google developers console). A slack API key, which can be generated on slack's bot page, and a remote MONGO_URI that can be used.

To buld and run:

```
make
```

## Contact

stuart@mongodb.com