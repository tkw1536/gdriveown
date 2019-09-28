package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"

	"google.golang.org/api/drive/v3"
)

// Options for a run of gdrivebackup
type Options struct {
	CredsPath string // Path to service account credentials

	FileID    string // ID of file to re-own
	FromEmail string // Email to transfer file from
	ToEmail   string // Email of new email
}

func readOptions() (opts Options) {

	// read the arguments from command line
	flag.StringVar(&opts.CredsPath, "credspath", "", "Path to service account credentials")
	flag.StringVar(&opts.FileID, "id", "", "Google Drive File ID to transfer file from. ")
	flag.StringVar(&opts.ToEmail, "to", "", "Email address to transfer file to. Should be the new owner. ")
	flag.StringVar(&opts.FromEmail, "from", "", "Email address to use to transfer file from. Should be the old owner. ")
	flag.Parse()

	// check that there aren't any extra arguments
	if flag.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "Too many arguments.")
		flag.Usage()
		os.Exit(2)
	}

	// check that credspath exists
	if _, err := os.Stat(opts.CredsPath); err != nil || opts.CredsPath == "" {
		fmt.Fprintf(os.Stderr, "credspath %q does not exist.\n", opts.CredsPath)
		flag.Usage()
		os.Exit(2)
	}

	// Contend must be provided
	if opts.FileID == "" {
		fmt.Fprintln(os.Stderr, "FileID (--id) must not be empty.")
		flag.Usage()
		os.Exit(2)
	}

	// --from
	if opts.FromEmail == "" {
		fmt.Fprintln(os.Stderr, "FromEmail (--from) must not be empty.")
		flag.Usage()
		os.Exit(2)
	}

	// --to
	if opts.ToEmail == "" {
		fmt.Fprintln(os.Stderr, "ToEmail (--to) must not be empty.")
		flag.Usage()
		os.Exit(2)
	}

	return
}

func newClientWith(opts Options, subject string) (service *drive.Service, err error) {
	bytes, err := ioutil.ReadFile(opts.CredsPath)
	if err != nil {
		return
	}

	var config *jwt.Config
	if config, err = google.JWTConfigFromJSON(bytes, drive.DriveScope); err != nil {
		return
	}
	config.Subject = subject

	client := config.Client(oauth2.NoContext)
	return drive.New(client)
}

func main() {
	// read in the options from command line
	opts := readOptions()

	// make a new client with the new email
	service, err := newClientWith(opts, opts.FromEmail)
	if err != nil {
		panic(err)
	}

	// create the new permission
	_, err = service.Permissions.Create(opts.FileID, &drive.Permission{
		EmailAddress: opts.ToEmail,
		Type:         "user",
		Role:         "owner",
	}).TransferOwnership(true).Do()
	if err != nil {
		panic(err)
	}

	// and done
	fmt.Printf("Transfered %q to %q\n", opts.FileID, opts.ToEmail)
}
