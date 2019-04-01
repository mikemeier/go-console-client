package socket

import (
	"console/both/message"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"net"
	"time"
)

var connection net.Conn
var IsConnected = make(chan bool)
var Messages = make(chan message.Response)
var Retries = 0

const rootPEM = `
-----BEGIN CERTIFICATE-----
MIIFaDCCBFCgAwIBAgISA53c4z7Ekw/2Cwy27dtg6vaKMA0GCSqGSIb3DQEBCwUA
MEoxCzAJBgNVBAYTAlVTMRYwFAYDVQQKEw1MZXQncyBFbmNyeXB0MSMwIQYDVQQD
ExpMZXQncyBFbmNyeXB0IEF1dGhvcml0eSBYMzAeFw0xOTAzMjgwODQyMDdaFw0x
OTA2MjYwODQyMDdaMBcxFTATBgNVBAMTDHRoZXRyYWRlci5jaDCCASIwDQYJKoZI
hvcNAQEBBQADggEPADCCAQoCggEBAMR1OF4ykV8faFCYiB5/dwb5zRAV/7anf7IZ
qiuNc/yM9mu0LL5mHLSeeUEDC6/1bpYK9P4try20O8Tn3UHuw7fsQvdpMsjru7vL
I43uOxUb8XVzIX7IdPoREgtjJH+3TsOWLSIsptrDzu1ou+Uqkyyip6WuD9mcrwiI
EHkpjgHNQfVZPu6C0crzUheVGCWrq6UsBZoFllvtwU4sSqE6Esw5nXxRAahKhJTM
Yavj8gu98wna5uBV7LRhAwjyEiAHu35bNNGlNxiKKpHyxP25Zsox8yAANN4VGHOw
xbgEUGQq4Pvy5q0ogx4BzAcD9RrXSrZCD19L7L3qt/I9bbpkROsCAwEAAaOCAnkw
ggJ1MA4GA1UdDwEB/wQEAwIFoDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUH
AwIwDAYDVR0TAQH/BAIwADAdBgNVHQ4EFgQUjguNA5WvvSIqtL5vLPIJZV63s6Aw
HwYDVR0jBBgwFoAUqEpqYwR93brm0Tm3pkVl7/Oo7KEwbwYIKwYBBQUHAQEEYzBh
MC4GCCsGAQUFBzABhiJodHRwOi8vb2NzcC5pbnQteDMubGV0c2VuY3J5cHQub3Jn
MC8GCCsGAQUFBzAChiNodHRwOi8vY2VydC5pbnQteDMubGV0c2VuY3J5cHQub3Jn
LzAtBgNVHREEJjAkghRjb25zb2xlLnRoZXRyYWRlci5jaIIMdGhldHJhZGVyLmNo
MEwGA1UdIARFMEMwCAYGZ4EMAQIBMDcGCysGAQQBgt8TAQEBMCgwJgYIKwYBBQUH
AgEWGmh0dHA6Ly9jcHMubGV0c2VuY3J5cHQub3JnMIIBBgYKKwYBBAHWeQIEAgSB
9wSB9ADyAHcA4mlLribo6UAJ6IYbtjuD1D7n/nSI+6SPKJMBnd3x2/4AAAFpw62W
zgAABAMASDBGAiEA4QWDQVl2l1Su5QUmctlXdv7h1nVP4HH1IzIx0AbwJ9kCIQDS
RQF8V+wY1uvLgkUmDsuql19K/0/jDm2AG2uUQc1YXwB3AGPy283oO8wszwtyhCdX
azOkjWF3j711pjixx2hUS9iNAAABacOtl0cAAAQDAEgwRgIhAIn+8X8Rh5yQS8W9
of26n/hDagA4PyzLNP6dLJgn83dMAiEAiJEtE9wT9A84A1NCCCLFCfqTDPa+xnFg
iUgxaqV+IfYwDQYJKoZIhvcNAQELBQADggEBAIMNH3lVp5XBlb7+ehb5z5N2r+L6
8TrcvVZTbhFYqd6LK16ejdwN/QbCv3n7wXk4aWSK5RnIkO2pJLcpStMgwi2lt3NX
B7SlnWL2jtGkiNVg25Yr7pUJIdVx97wkl9qL+Rf58arWNitUL0oanywa7+nELrSC
Q5ujpZbq9WD9D7GP2Faafy5uC5loFsVrkrTJJBMFSpk5jgY3dbNqUT52I7jmaxYv
p5FmWZIFvmk93t4FoxqAm7AhicydT6J8B+HmSelZzRGxCRTjYQhs1xqSUmuntU1J
PnoeeR2MElUAAMabksi6jdDCLvXX1VJkHXCMhkEb9aBeCyX5RHGQ/jEco/I=
-----END CERTIFICATE-----`

func Heartbeat(host string, useTls bool) {
	defer func() {
		time.AfterFunc(time.Second, func() {
			Heartbeat(host, useTls)
		})
	}()

	if connection == nil {
		var c net.Conn
		var connErr error
		if useTls {
			roots := x509.NewCertPool()
			ok := roots.AppendCertsFromPEM([]byte(rootPEM))
			if !ok {
				panic("failed to parse root certificate")
			}
			c, connErr = tls.Dial("tcp", host, &tls.Config{RootCAs: roots})
		} else {
			c, connErr = net.Dial("tcp", host)
		}
		if connErr == nil {
			setConnected(c)
		} else {
			setConnected(nil)
		}
		return
	}

	if Send("ping-check-heartbeat") {
		setConnected(connection)
	} else {
		setConnected(nil)
	}
}

func ListenForData() {
	defer func() {
		time.AfterFunc(time.Second, func() {
			ListenForData()
		})
	}()

	if connection == nil {
		return
	}

	for {
		msgString, err := message.ReadFromSocket(connection)
		if err != nil {
			break
		}

		var msg message.Response
		json.Unmarshal([]byte(msgString), &msg)

		Messages <- msg
	}
}

func Send(cmd string) bool {
	if connection == nil {
		return false
	}

	return message.Send(connection, []byte(cmd))
}

func Disconnect() {
	if connection != nil {
		connection.Close()
	}
	close(IsConnected)
}

func setConnected(c net.Conn) {
	if c != nil {
		connection = c
		Retries = 0
		IsConnected <- true
	} else {
		connection = nil
		Retries++
		IsConnected <- false
	}
}
