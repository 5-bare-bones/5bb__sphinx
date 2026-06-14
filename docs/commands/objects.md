## Concepts

## Description

Sphinx defines the following objects"

| Object   |    Description                        |
|----------|---------------------------------------|
| Vault    | the file holding the other objects    |
| |- Entry | default login entry                   |
| |- File  | a file object in the vault            |
| |- Card  | a credit card in the vault            |
| |- TOPT  | a two-factor authentication code / one time password.  |

## Entries 

The entry holds the following elements:

- username
- password
- URL
- expiration date (iso-format)
- notes

## Files

The file holds the following elements:

- filename
- content
- size
- creation date (iso-format)
- modification date (updated on/iso-format)

## Card 

The card holds the following elements:

- cardholder name
- card number
- expiration date (iso-format)
- security code (CVV)
- notes 

## TOTP (Time-based One-Time Password)

The TOTP holds the following elements:

- secret
- name
- issuer
