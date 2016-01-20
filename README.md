# Docker task runner for cloudstack

## Components
1. *supervisor* is the guy who starts new instanses, it will start an new instance and upload a cloud-configuration file which installs docker and runs the desired image.
2. *janitor* is the cleaner. He will periodically (well, a cronjob needs to be created) destroy instanses within the group that is not running any docker containers. A warmup period of 30 minutes is used to not destroy instanses before they get the chance to start docker.
3. *uploader* is a tool which uploads job results (FAST2 specifics maybe?) to swift object storage.


# Variables (used for supervisor and janitor)
export SWIFT_API_USER="username"

export SWIFT_API_KEY="apikey"

export SWIFT_AUTH_URL="authurl"

export RBC_API_KEY="apikey"

export RBC_SECRET="secret"

export WPAU_SLACK_HOOK_URL="slack-url-if-desired"


# Cronjob
Use with cronjob, supervisor once a day and janitor every 20 min
