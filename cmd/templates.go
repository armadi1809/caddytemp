package cmd

var placeholdersData = map[string][]string{
	"Static Files":               {"hostname", "root path"},
	"Reverse proxy all requests": {"hostname", "reverse proxy hostname"},
	"Reverse proxy only requests starting with a given path": {"hostname", "root path", "matching path", "reverse proxy hostname"},
	"PHP-FPM":    {"hostname", "root path", "fastcgi server address"},
	"FrankenPHP": {"hostname", "root path"},
	"Add www. subdomain with an HTTP redirect":                         {"hostname"},
	"Remove www. subdomain with an HTTP redirect":                      {"hostname"},
	"Remove www. subdomaing for multiple domains at once":              {"hostname 1", "hostname 2"},
	"Remove trailing slashes internally (using the rewrite directive)": {"hostname"},
	"Remove trailing slashes externally (using the redir directive)":   {"hostname"},
	"Serve multiple subdomains with the same wildcard certificate":     {"hostname", "subdomain 1", "subdomain 2"},
	"Single-Page Application (SPA) with no backend":                    {"hostname", "root path"},
	"Single-Page Application (SPA) with a backend api":                 {"hostname", "backend base address", "root path"},
}

var templates map[string]string = map[string]string{

	"Static Files": `{hostname} {
		root * {root path}
		file_server
}`,

	"Reverse proxy all requests": `{hostname} {
	reverse_proxy {reverse proxy hostname}
}`,

	"Reverse proxy only requests starting with a given path, and serve static files otherwise": `{hostname} {
	root * {root path}
	reverse_proxy {matching path}/* {reverse proxy hostname}
	file_server
}`,

	"PHP-FPM": `{hostname} {
	root * {root path}
	encode gzip
	php_fastcgi {fastcgi server address}
	file_server
}`,

	"FrankenPHP": `{
    frankenphp
    order php_server before file_server
}

{hostnmae} {
	root * {root path}
    encode zstd br gzip
    php_server
}`,

	"Add www. subdomain with an HTTP redirect": `{hostname} {
	redir https://www.{host}{uri}
}

www.{hostname} {
}`,

	"Remove www. subdomain with an HTTP redirect": `www.{hostname} {
	redir https://{hostname}{uri}
}

example.com {
}`,

	"Remove www. subdomaing for multiple domains at once": `www.{hostname 1}, www.{hostname 2} {
	redir https://{labels.1}.{labels.0}{uri}
}

{hostname 1}, {hostname 2} {
}`,

	"Remove trailing slashes internally (using the rewrite directive)": `{hostname} {
	rewrite /add     /add/
	rewrite /remove/ /remove
}`,

	"Remove trailing slashes externally (using the redir directive)": `{hostname} {
	redir /add     /add/
	redir /remove/ /remove
}`,

	"Serve multiple subdomains with the same wildcard certificate": `*.{hostname} {
	tls {
		dns <provider_name> [<params...>]
	}

	@{subdomain 1} host {subdomain 1}.example.com
	handle @{subdomain 1} {
		respond "{subdomain 1}!"
	}

	@{subdomain 2} host {subdomain 2}.example.com
	handle @{subdomain 2} {
		respond "{subdomain 2}!"
	}

	# Fallback for otherwise unhandled domains
	handle {
		abort
	}
}`,

	"Single-Page Application (SPA) with no backend": `{hostname} {
	root * {root path}
	encode gzip
	try_files {path} /index.html
	file_server
}`,

	"Single-Page Application (SPA) with a backend api": `{hostname} {
	encode gzip

	handle {backend base address}* {
		reverse_proxy backend:8000
	}

	handle {
		root * {root path}
		try_files {path} /index.html
		file_server
	}
}`,

	// 	"Caddy proxying to another Caddy (Front instance)": `foo.{hostname}, bar.{hostname} {
	// 	reverse_proxy {reverse proxy address}
	// }`,

	// 	"Caddy proxying to another Caddy (Back instance)": `{
	// 	servers {
	// 		trusted_proxies static private_ranges
	// 	}
	// }

	// http://foo.example.com {
	// 	reverse_proxy foo-app:8080
	// }

	//	http://bar.example.com {
	//		reverse_proxy bar-app:9000
	//	}`,
}
