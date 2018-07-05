## sysInfoWebServer
YMMV.

I'm bothered that many of the standard daemons at home (OSX/Win7/Linux-ARM) might overheat in the Japanese summer. Most shareware out there is operationally expensive. The hope is that this go-package is lightweight enought to allow me to track statistics over time - eventually looking for trending/cooling opportunities (read: playtime) at home.

### Why a Dockerfile?

As it turns out - this is entirely useless in its current form as the OSX Docker Daemon is a linux VM which means every call to the temperature function returns 1... Yikes
