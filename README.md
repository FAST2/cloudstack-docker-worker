# Cloudstack utils
These are small utility scripts used for managing cloudstack project as a worker pool.

## Scheduled components
These components should be run with crontab or other scheduler.

- *supervisor* starts new instances, will start a new instance and upload a cloud-configuration file which installs docker and runs the desired image.
- *janitor*  cleaner, will periodically (well, a cronjob needs to be created) destroy instanses within the group that is not running any docker containers. A warmup period of 30 minutes is applied as to not destroy instances before they've had time to start docker.
- *cleaner* Cleans old jobs and files associated once the CLEAN_OLD_JOBS days has passed.
- *project_pusher* utility used for uploading files to object storage
- *screamer* Shouts (loud) when an instance has been running for a long long time (and is using YOUR money).

## Utils
- *uploader* is a tool which uploads job results (FAST2 specifics maybe?) to swift object storage.
- *project_pusher* is a tool which uploads a new project

# Environment variables used for the above components
```
export OS_APPLICATION_CREDENTIAL_ID="application credential ID"
export OS_APPLICATION_CREDENTIAL_SECRET="application credential PASSWORD"
export OS_AUTH_URL="authentication url"
export RBC_API_KEY="apikey"
export RBC_SECRET="secret"
export WPAU_SLACK_HOOK_URL="slack-url-if-desired"
```
