package cmd

var placeholdersData = map[string][]string{
	"Static Files": {"hostname"},
}

var templates map[string]string = map[string]string{

	"Static Files": `%s {
		root * /var/www
		file_server
	}`,

	"Reverse proxy all requests": `example.com {
	reverse_proxy localhost:5000
}`,

	"Reverse proxy only requests starting with a given path": `example.com {
	root * /var/www
	reverse_proxy /api/* localhost:5000
	file_server
}`,

	"PHP-FPM": `example.com {
	root * /srv/public
	encode gzip
	php_fastcgi localhost:9000
	file_server
}`,

	"FrankenPHP": `{
    frankenphp
    order php_server before file_server
}

example.com {
	root * /srv/public
    encode zstd br gzip
    php_server
}
`,

	"Add www. subdomain with an HTTP redirect": `example.com {
	redir https://www.{host}{uri}
}

www.example.com {
}`,

	"Remove www. subdomain with an HTTP redirect": `www.example.com {
	redir https://example.com{uri}
}

example.com {
}`,

	"Remove www. subdomaing for multiple domains at once": `www.example-one.com, www.example-two.com {
	redir https://{labels.1}.{labels.0}{uri}
}

example-one.com, example-two.com {
}`,

	"Remove trailing slashes internally (using the rewrite directive)": `example.com {
	rewrite /add     /add/
	rewrite /remove/ /remove
}`,

	"Remove trailing slashes externally (using the redir directive)": `example.com {
	redir /add     /add/
	redir /remove/ /remove
}`,

	"Serve multiple domains with the same wildcard certificate": `*.example.com {
	tls {
		dns <provider_name> [<params...>]
	}

	@foo host foo.example.com
	handle @foo {
		respond "Foo!"
	}

	@bar host bar.example.com
	handle @bar {
		respond "Bar!"
	}

	# Fallback for otherwise unhandled domains
	handle {
		abort
	}
}`,

	"Single-Page Application (SPA) with no backend": `example.com {
	root * /srv
	encode gzip
	try_files {path} /index.html
	file_server
}`,

	"Single-Page Application (SPA) with a backend api": `example.com {
	encode gzip

	handle /api/* {
		reverse_proxy backend:8000
	}

	handle {
		root * /srv
		try_files {path} /index.html
		file_server
	}
}`,

	"Caddy proxying to another Caddy (Front instance)": `foo.example.com, bar.example.com {
	reverse_proxy 10.0.0.1:80
}`,

	"Caddy proxying to another Caddy (Back instance)": `{
	servers {
		trusted_proxies static private_ranges
	}
}

http://foo.example.com {
	reverse_proxy foo-app:8080
}

http://bar.example.com {
	reverse_proxy bar-app:9000
}`,
}
