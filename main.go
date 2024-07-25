package main

import (
	"html/template"
	"net"
	"net/http"
)

type DNSRecord struct {
	Domain    string   `json:"domain"`
	ARecord   []net.IP `json:"a_record"`
	CNAME     string   `json:"cname_record"`
	MXRecord  []string `json:"mx_record"`
	NSRecord  []string `json:"ns_record"`
	TXTRecord []string `json:"txt_record"`
}

var templates *template.Template

func main() {
	templates = template.Must(template.ParseGlob("templates/*.html"))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/lookup", dnsLookupHandler)
	http.ListenAndServe(":8000", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}

func dnsLookupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	} else if r.Method != http.MethodPost {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	domain := r.FormValue("domain")

	RecordTypeA, _ := net.LookupIP(domain)
	RecordTypeCNAME, _ := net.LookupCNAME(domain)
	RecordTypeMX, _ := net.LookupMX(domain)
	RecordTypeNS, _ := net.LookupNS(domain)
	RecordTypeTXT, _ := net.LookupTXT(domain)

	NSHosts := []string{}
	for _, ns := range RecordTypeNS {
		NSHosts = append(NSHosts, ns.Host)
	}

	MXHosts := []string{}
	for _, mx := range RecordTypeMX {
		MXHosts = append(MXHosts, mx.Host)
	}

	dnsRecord := &DNSRecord{
		Domain:    domain,
		ARecord:   RecordTypeA,
		CNAME:     RecordTypeCNAME,
		MXRecord:  MXHosts,
		NSRecord:  NSHosts,
		TXTRecord: RecordTypeTXT,
	}

	templates.ExecuteTemplate(w, "index.html", map[string]any{
		"DNSRecord": dnsRecord,
	})
}
