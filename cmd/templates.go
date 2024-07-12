package cmd

type template struct {
	Label   string
	Content string
}

var templates []template = []template{
	{
		Label: "Static Files",
		Content: `example.com {
		root * /var/www
		file_server
	}`},
	{
		Label: "Reverse proxy all requests",
		Content: `example.com {
	reverse_proxy localhost:5000
}`,
	},
	{
		Label: "Reverse proxy only requests starting with a given path",
		Content: `example.com {
	root * /var/www
	reverse_proxy /api/* localhost:5000
	file_server
}`,
	},
	{
		Label: "PHP-FPM",
		Content: `example.com {
	root * /srv/public
	encode gzip
	php_fastcgi localhost:9000
	file_server
}`,
	},
	{
		Label: "FrankenPHP",
		Content: `{
    frankenphp
    order php_server before file_server
}

example.com {
	root * /srv/public
    encode zstd br gzip
    php_server
}
`,
	},
	{
		Label: "Add www. subdomain with an HTTP redirect",
		Content: `example.com {
	redir https://www.{host}{uri}
}

www.example.com {
}`,
	},
	{
		Label: "Remove www. subdomain with an HTTP redirect",
		Content: `www.example.com {
	redir https://example.com{uri}
}

example.com {
}`,
	},
	{
		Label: "Remove www. subdomaing for multiple domains at once",
		Content: `www.example-one.com, www.example-two.com {
	redir https://{labels.1}.{labels.0}{uri}
}

example-one.com, example-two.com {
}`,
	},
	{
		Label: "Remove trailing slashes internally (using the rewrite directive)",
		Content: `example.com {
	rewrite /add     /add/
	rewrite /remove/ /remove
}`,
	},
	{
		Label: "Remove trailing slashes externally (using the redir directive)",
		Content: `example.com {
	redir /add     /add/
	redir /remove/ /remove
}`,
	},
	{
		Label: "Serve multiple domains with the same wildcard certificate",
		Content: `*.example.com {
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
	},
	{
		Label: "Single-Page Application (SPA) with no backend",
		Content: `example.com {
	root * /srv
	encode gzip
	try_files {path} /index.html
	file_server
}`,
	},
	{
		Label: "Single-Page Application (SPA) with a backend api",
		Content: `example.com {
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
	},
	{
		Label: "Caddy proxying to another Caddy (Front instance)",
		Content: `foo.example.com, bar.example.com {
	reverse_proxy 10.0.0.1:80
}`,
	},
	{
		Label: "Caddy proxying to another Caddy (Back instance)",
		Content: `{
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
	},
}
