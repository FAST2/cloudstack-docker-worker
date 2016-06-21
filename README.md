# Cloudstack utils
Small util scripts used for managing an cloudstack infrastructure as a worker pool.

## Schedelued components
- *supervisor* is the guy who starts new instanses, it will start an new instance and upload a cloud-configuration file which installs docker and runs the desired image.
- *janitor* is the cleaner. He will periodically (well, a cronjob needs to be created) destroy instanses within the group that is not running any docker containers. A warmup period of 30 minutes is used to not destroy instanses before they get the chance to start docker.
- *cleaner* Cleans old jobs and files associated once the CLEAN_OLD_JOBS days has passed.
- *project_pusher* util used for uploading files to object storage
- *screamer* Shouts (loud) when an instance has been running for a long long time (and is using YOUR money).

## Utils
- *uploader* is a tool which uploads job results (FAST2 specifics maybe?) to swift object storage.
- *project_pusher* is a tool which uploads a new project

# Environment variables used for the above components
```
export SWIFT_API_USER="username"
export SWIFT_API_KEY="apikey"
export SWIFT_AUTH_URL="authurl"
export RBC_API_KEY="apikey"
export RBC_SECRET="secret"
export WPAU_SLACK_HOOK_URL="slack-url-if-desired"
```

# Cronjob
Use with cronjob, supervisor once a day and janitor every 20 min
