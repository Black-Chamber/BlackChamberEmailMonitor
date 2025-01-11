# Black Chamber

## Overview
Black Chamber is tool designed to help IT teams detect unauthorized SaaS usage within an organization's. It current function is to provide a lackluster but free method detecting Shadow IT in organizations where networking monitoring is not enough (BYOD). 

Black Chamber uses community managed detection rules (regex) in order to log the detected message id, recipent and suspected service.

## How it works

Black Chamber connects to the outlook reporting service at regular intervals to grab messagetrace information. While this does not contain the message content, it provides the sender and recipent infromation in addition to the subject. This is more than enough to determine of mark recipents who may be using unauthorized SaaS applications

**Example:** A recipent recieving a notification that someone has accepted the invite to join a file sharing site.


## Key Benefits

- **Enhanced Visibility:** Gain a better understanding of SaaS usage in your enviroment.
- **Improved Security:** Detect and mitigate risks associated with unauthorized usage.

## Getting Started

1. **Download Black Chamber:** Follow the installation guide to set up agents and integrations.
2. **Edit The Settings:** Customize discovery and reporting settings based on your organization's requirements.
4. **Review Finding:** Access the sqllite database to find and sort through detections

## Supported Platforms
Black Chamber is compatible with:
- Microsoft 365 (M365) email environments.


