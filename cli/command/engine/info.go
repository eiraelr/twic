package engine

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/juliengk/go-cert/pkix"
	"github.com/juliengk/go-utils/user"
	"github.com/juliengk/stack/client"
	"github.com/kassisol/twic/pkg/adf"
	"github.com/kassisol/twic/pkg/date"
	"github.com/spf13/cobra"
)

func newInfoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "info",
		Short:   "Information about Docker engine certificate",
		Long:    infoDescription,
		Run:     runInfo,
	}

	return cmd
}

func runInfo(cmd *cobra.Command, args []string) {
	user, err := user.New()
	if err != nil {
		log.Fatal(err)
	}

	if !user.IsRoot() {
		log.Fatal("You must be root to run engine subcommand")
	}

	if len(args) > 0 {
		cmd.Usage()
		os.Exit(-1)
	}

	config := adf.New("engine")

	if err := config.Init(); err != nil {
		log.Fatal(err)
	}

	cf, err := config.CertFilesName()
	if err != nil {
		log.Fatal(err)
	}

	certificate, err := pkix.NewCertificateFromPEMFile(cf.Crt)
	if err != nil {
		log.Fatal(err)
	}

	crldp := certificate.Crt.CRLDistributionPoints[0]

	url, err := client.ParseUrl(crldp)
	if err != nil {
		log.Fatal(err)
	}

	tsaurl := fmt.Sprintf("%s://%s", url.Scheme, url.Host)
	if url.Port != 443 {
		tsaurl = fmt.Sprintf("%s://%s:%d", url.Scheme, url.Host, url.Port)
	}

	cn := certificate.Crt.Subject.CommonName
	expire := date.ExpireDateString(certificate.Crt.NotAfter)

	fmt.Println("TSA URL:", tsaurl)
	fmt.Println("CN:", cn)
	fmt.Println("Expire:", expire)
}

var infoDescription = `
Information about Docker engine certificate

`
