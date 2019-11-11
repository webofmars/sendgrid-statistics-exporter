# SendGrid Prometheus Exporter

## install

```
docker pull riccardopomato/sendgrid-statistics-exporter
```

## environment

* LISTEN_ADDR
* METRICS_ENDPOINT
* SENDGRID_API_KEY

## Metrics
These are the metrics collected using this exporter.
There are two main groups one is collect the metrics aggregated by day and another one aggregated by month.
```
DAILY aggregated:
sendgrid_dailyblocks{name="",type=""}
sendgrid_dailybounce_drops{name="",type=""}
sendgrid_dailybounces{name="",type=""}
sendgrid_dailyclicks{name="",type=""}
sendgrid_dailydeferred{name="",type=""}
sendgrid_dailydelivered{name="",type=""}
sendgrid_dailyinvalid_emails{name="",type=""}
sendgrid_dailyopens{name="",type=""}
sendgrid_dailyprocessed{name="",type=""}
sendgrid_dailyrequests{name="",type=""}
sendgrid_dailyspam_report_drops{name="",type=""}
sendgrid_dailyspam_reports{name="",type=""}
sendgrid_dailyunique_clicks{name="",type=""}
sendgrid_dailyunique_opens{name="",type=""}
sendgrid_dailyunsubscribe_drops{name="",type=""}
sendgrid_dailyunsubscribes{name="",type=""}
```

```
MONTHLY aggregated:
sendgrid_monthlyblocks{name="",type=""}
sendgrid_monthlybounce_drops{name="",type=""}
sendgrid_monthlybounces{name="",type=""}
sendgrid_monthlyclicks{name="",type=""}
sendgrid_monthlydeferred{name="",type=""}
sendgrid_monthlydelivered{name="",type=""}
sendgrid_monthlyinvalid_emails{name="",type=""}
sendgrid_monthlyopens{name="",type=""}
sendgrid_monthlyprocessed{name="",type=""}
sendgrid_monthlyrequests{name="",type=""}
sendgrid_monthlyspam_report_drops{name="",type=""}
sendgrid_monthlyspam_reports{name="",type=""}
sendgrid_monthlyunique_clicks{name="",type=""}
sendgrid_monthlyunique_opens{name="",type=""}
sendgrid_monthlyunsubscribe_drops{name="",type=""}
sendgrid_monthlyunsubscribes{name="",type=""}
```

```
sendgrid_up
```