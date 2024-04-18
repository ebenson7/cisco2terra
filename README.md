# Cisco2Terra

A tool built to be used for building out Meraki infrastructure, transforming ASA configuration to Terraform HCL file(s) using the Hashicorp hclwrite Golang library. Currently a prototype with very basic functionality.

There are a lot of things to implement in future feature sets:
- Integrate with ServiceNow/Jira/BitBucket so that once a configuration file gets generated, a Change Request ticket is opened, assigned, pushed, and deployed.
- Versioning Control of current Meraki infrastructure
- Automatically scraping ASA infrastructure to generate configuration files.
- Pushing metrics to monitoring service such as Prometheus, Thanos, Datadog, etc. for milestones regarding progress made of all infrastructure moved over.
- Much more that I'm probably not thinking of

The current limitation is it only takes one file at a time. Once basic functionality is established, the next step would be automating infrastructure changes.

### I think the ideal situation is to not use this tool as a method of transferring full infrastructure over and not passing traffic through the ASA, but to mirror configuration and slowly push traffic towards a Cloud-based SD-WAN network.

