package client

import (
	"strings"
	"time"

	"github.com/datakit-dev/protoc-gen-cobra/iocodec"
	"github.com/datakit-dev/protoc-gen-cobra/naming"
	"github.com/spf13/pflag"
)

type Option func(*Config)

func WithServerAddrFlag() Option {
	return func(c *Config) {
		c.flagBinders = append(c.flagBinders, func(fs *pflag.FlagSet, n naming.Namer) {
			fs.StringVarP(&c.ServerAddr, n("ServerAddr"), "s", c.ServerAddr, "server address in the form host:port")
		})
	}
}

func WithRequestFileFlag() Option {
	return func(c *Config) {
		c.flagBinders = append(c.flagBinders, func(fs *pflag.FlagSet, n naming.Namer) {
			fs.StringVarP(&c.RequestFile, n("RequestFile"), "f", c.RequestFile, "client request file; use \"-\" for stdin")
		})
	}
}

func WithRequestFormatFlag() Option {
	return func(c *Config) {
		c.flagBinders = append(c.flagBinders, func(fs *pflag.FlagSet, n naming.Namer) {
			fs.StringVarP(&c.RequestFormat, n("RequestFormat"), "i", c.RequestFormat, "request format ("+strings.Join(c.decoderFormats(), ", ")+")")
		})
	}
}

func WithResponseFormatFlag() Option {
	return func(c *Config) {
		c.flagBinders = append(c.flagBinders, func(fs *pflag.FlagSet, n naming.Namer) {
			fs.StringVarP(&c.ResponseFormat, n("ResponseFormat"), "o", c.ResponseFormat, "response format ("+strings.Join(c.encoderFormats(), ", ")+")")
		})
	}
}

func WithTimeoutFlag() Option {
	return func(c *Config) {
		c.flagBinders = append(c.flagBinders, func(fs *pflag.FlagSet, n naming.Namer) {
			fs.DurationVar(&c.Timeout, n("Timeout"), c.Timeout, "RPC timeout")
		})
	}
}

func WithTLSFlag() Option {
	return func(c *Config) {
		c.flagBinders = append(c.flagBinders, func(fs *pflag.FlagSet, n naming.Namer) {
			fs.BoolVar(&c.TLS, n("TLS"), c.TLS, "enable TLS")
		})
	}
}

func WithServerNameFlag() Option {
	return func(c *Config) {
		c.flagBinders = append(c.flagBinders, func(fs *pflag.FlagSet, n naming.Namer) {
			fs.StringVar(&c.ServerName, n("TLS ServerName"), c.ServerName, "TLS server name override")
		})
	}
}

func WithInsecureSkipVerifyFlag() Option {
	return func(c *Config) {
		c.flagBinders = append(c.flagBinders, func(fs *pflag.FlagSet, n naming.Namer) {
			fs.BoolVar(&c.InsecureSkipVerify, n("TLS InsecureSkipVerify"), c.InsecureSkipVerify, "INSECURE: skip TLS checks")
		})
	}
}

func WithCACertFileFlag() Option {
	return func(c *Config) {
		c.flagBinders = append(c.flagBinders, func(fs *pflag.FlagSet, n naming.Namer) {
			fs.StringVar(&c.CACertFile, n("TLS CACertFile"), c.CACertFile, "CA certificate file")
		})
	}
}

func WithCertFileFlag() Option {
	return func(c *Config) {
		c.flagBinders = append(c.flagBinders, func(fs *pflag.FlagSet, n naming.Namer) {
			fs.StringVar(&c.CertFile, n("TLS CertFile"), c.CertFile, "client certificate file")
		})
	}
}

func WithKeyFileFlag() Option {
	return func(c *Config) {
		c.flagBinders = append(c.flagBinders, func(fs *pflag.FlagSet, n naming.Namer) {
			fs.StringVar(&c.KeyFile, n("TLS KeyFile"), c.KeyFile, "client key file")
		})
	}
}

func WithServerAddr(addr string) Option {
	return func(c *Config) {
		c.ServerAddr = addr
	}
}

func WithRequestFormat(format string) Option {
	return func(c *Config) {
		c.RequestFormat = format
	}
}

func WithResponseFormat(format string) Option {
	return func(c *Config) {
		c.ResponseFormat = format
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

func WithEnvVars(prefix string) Option {
	return func(c *Config) {
		c.UseEnvVars = true
		c.EnvVarPrefix = prefix
	}
}

func WithCommandNamer(namer naming.Namer) Option {
	return func(c *Config) {
		c.CommandNamer = namer
	}
}

func WithFlagNamer(namer naming.Namer) Option {
	return func(c *Config) {
		c.FlagNamer = namer
	}
}

func WithEnvVarNamer(namer naming.Namer) Option {
	return func(c *Config) {
		c.EnvVarNamer = namer
	}
}

func WithTLSCACertFile(certFile string) Option {
	return func(c *Config) {
		c.TLS = true
		c.CACertFile = certFile
	}
}

func WithTLSCertFile(certFile, keyFile string) Option {
	return func(c *Config) {
		c.TLS = true
		c.CertFile = certFile
		c.KeyFile = keyFile
	}
}

func WithTLSServerName(serverName string) Option {
	return func(c *Config) {
		c.TLS = true
		c.ServerName = serverName
	}
}

func WithFlagBinder(binder FlagBinder) Option {
	return func(c *Config) {
		c.flagBinders = append(c.flagBinders, binder)
	}
}

func WithPreDialer(dialer PreDialer) Option {
	return func(c *Config) {
		c.preDialers = append(c.preDialers, dialer)
	}
}

func WithInputDecoder(format string, maker iocodec.DecoderMaker) Option {
	return func(c *Config) {
		d := make(map[string]iocodec.DecoderMaker, len(c.inDecoders)+1)
		for k, v := range c.inDecoders {
			d[k] = v
		}
		d[format] = maker
		c.inDecoders = d
	}
}

func WithOutputEncoder(format string, maker iocodec.EncoderMaker) Option {
	return func(c *Config) {
		e := make(map[string]iocodec.EncoderMaker, len(c.outEncoders)+1)
		for k, v := range c.outEncoders {
			e[k] = v
		}
		e[format] = maker
		c.outEncoders = e
	}
}
