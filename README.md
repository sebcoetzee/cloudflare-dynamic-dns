# Cloudflare Dynamic DNS

Cloudflare Dynamic DNS is a program that allows the use of Cloudflare's DNS servers as a Dynamic DNS service. The program checks the public IP address of the host machine every 5 minutes and updates the DNS record on Cloudflare with the IP address.

## Usage

To run the program, execute the binary for your platform. This example uses the binary compiled for Mac OS X.

```
./bin/cloudflare-dynamic-dns-darwin-amd64 --api_token="<YOUR_CLOUDFLARE_API_TOKEN>" --zone_name="example.com" --subdomain="myhome"
```

The `<YOUR_CLOUDFLARE_API_TOKEN>` is your Cloudflare token. For Cloudflare Dynamic DNS to work, the API Token needs at least the following permissions for the specified `zone_name`:

- Zone - Zone - Read
- Zone - DNS - Edit
