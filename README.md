dpaycol
==============

dpayol is a command line tool for collecting statistics which arise during a payroll run. dpacol can be used with following parameters:

        dpaycol -ak 70 -am 04 -lt 50 -ut 50 -jn KIDIJOB -jkid 0001
        dpaycol -ak 70 -am 04 -lt 50 -ut 50 -jn KIDIJOB -jkid 0001 -e -rc 0

Description of the used parameters:

- ak: Abrechnugskreis
- am: Abrechnungsmonat
- lt: Lauftermin
- ut: Untertermin
- jn: Jobname
- jkid: ID der Jobkette
- e: Ende Flag
- rc: Return Code des ausgeführten Jobs

The output is written to json encoded file named stats-dpay.json.

Build instructions:

- docker build -t dataport.de/dpaycol --no-cache --build-arg http_proxy=http://proxy:80 --build-arg https_proxy=http://proxy:80 .
- docker create -it --name dpaycol dataport.de/dpaycol
- docker cp dpaycol:/go/src/app/dpaycol .
