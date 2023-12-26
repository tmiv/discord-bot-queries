# Discord Bot Queries

Various Discord Bot Queries hosted as a service

## Env vars
- DISCORD_TOKEN - Bot token for discord
- SKIP_OIDC - Disable OIDC Check by setting to TRUE
- CORS_ORIGINS - CORS allowed origins
- SECURITY_AUDIENCE - Security Audience to confirm OIDC JWT check
- SECURITY_ISSUER - Security Issuer to confirm OIDC JWT check
- SECURITY_ALLOW - Comma separated list of email addresses to confirm OIDC JWT check 

## Endpoints
### /v1/VerifyMembership
Check a user's membership to a channel
#### Parameters
- user - UserID of the user
- channel - ChannelID of the channel
#### Returns
```json
    {"member":false,"user_id":"1184213922513948886","channel_id":"1185732803762069584"}
```

