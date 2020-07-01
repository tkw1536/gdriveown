# gdrivebackup

This is a quick tool to transfer a google drive file owner from one user to another. 
Needs a service account with domain-wide delegation enabled. 

Usage:

```
Usage of ./gdriveown:
  -credspath string
        Path to service account credentials
  -from string
        Email address to use to transfer file from. Should be the old owner. 
  -id string
        Google Drive File ID to transfer file from. 
  -to string
        Email address to transfer file to. Should be the new owner. 
```

## License

GPL3, see [LICENSE](LICENSE)