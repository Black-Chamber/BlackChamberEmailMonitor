# Black Chamber Email Monitor
Shadow IT and SaaS via Email Log Monitoring

# ALPHA - IN ACTIVE DEVELOPMENT


## Overview
Black Chamber is a tool designed to help IT teams detect unauthorized SaaS usage within an organization. Its current function is to provide a lackluster but free method detecting Shadow IT when agent based and networking monitoring is not enough (BYOD). 

Black Chamber uses community driven detection rules (basic regex) in order to log the detected message id, recipent and suspected SaaS service.

## Contributing 
Please use the development branch

This tool isnt much, but it has the potential help many organizations manage risk with no additional costs. While selecting the methodology of detection, M365 MessageTrace was selected becuase it was available to all organizations regardless of license (vs EDiscovery Components)

**Anyone can help this project immensely by identifying risky SaaS technologies and creating regex detection rules under Services**


## How it works

![image](https://github.com/user-attachments/assets/ca870cad-13b7-4aca-8f03-cd53ce226a62)

Black Chamber connects to the outlook reporting service at regular intervals to grab messagetrace logs. While this does not contain the message content, it provides the sender and recipent information in addition to the subject. This is more than enough to flag recipents who may be using unauthorized SaaS applications

**Example:** A recipent recieving a notification that someone has accepted the invite to join a file sharing site.


## Key Benefits

- **Enhanced Visibility:** Gain a better understanding of SaaS usage in your enviroment.
- **Improved Security:** Detect and mitigate risks associated with unauthorized usage.

## Getting Started

1. **Download Black Chamber:** Follow the installation guide to set up agents and integrations.
2. **Configure Permissions In M365:** Follow the next below to configure an App Registration with the required permissions
3. **Edit The Settings:** Customize the configuration file with the folder path containing the detection rules and your app registration rules
4. **Review Finding:** Access the sqllite database to find and sort through detections

## Running The Application
The application can be ran with "Go run " or built into an executable using "Go Build ./..." on the root folder
It is important you configure the config.yaml and reference it properly be either changing the reference in main.go OR passing the --cfg flag

```

go run cmd/bcem/main.go --cfg=../../config/config.yml

```
Or, if you have already built the executable:
```
./bcem --cfg=../../config/config.yml
```




## M365 Permission Configuration
1. **Create An Application Registration:** follow this guide to create an app registration https://learn.microsoft.com/en-us/graph/auth-register-app-v2
2. **Assign API Permissions:** Under the App Registration, Select API Permission > Add A Permission > APIs my organization uses > Office 365 Exchange Online. **ADD BOTH OF THE FOLLOWING*
   - A: Delegated > "ReportingWebService.Read"
   - B: Application > "ReportingWebService.Read.All"
5. **Assign Role:** Assign the "Global Reader" Role to the Application Registration (TIP: when assigning the role, you may only see a tab for owners, start typing the application registartion name and it will appear)


## Supported Platforms
Black Chamber is compatible with:
- Microsoft 365 (M365) email environments.



