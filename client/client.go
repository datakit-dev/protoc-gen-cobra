package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/datakit-dev/protoc-gen-cobra/iocodec"
	"github.com/datakit-dev/protoc-gen-cobra/naming"
)

type (
	FlagBinder func(*pflag.FlagSet, naming.Namer)
	PreDialer  func(context.Context, *[]grpc.DialOption) error
)

type Config struct {
	GetContextFunc func(context.Context) (context.Context, error)
	ClientConnFunc func() (*grpc.ClientConn, error)
	PreDecoder     func(context.Context) func(any) error
	ServerAddr     string
	RequestFile    string
	RequestFormat  string
	ResponseFormat string
	Timeout        time.Duration
	UseEnvVars     bool
	EnvVarPrefix   string

	CommandNamer naming.Namer
	FlagNamer    naming.Namer
	EnvVarNamer  naming.Namer

	TLS                bool
	ServerName         string
	InsecureSkipVerify bool
	CACertFile         string
	CertFile           string
	KeyFile            string

	flagBinders []FlagBinder
	preDialers  []PreDialer
	inDecoders  map[string]iocodec.DecoderMaker
	outEncoders map[string]iocodec.EncoderMaker
}

var DefaultConfig = &Config{
	ServerAddr:     "localhost:8080",
	RequestFormat:  "json",
	ResponseFormat: "json",

	CommandNamer: naming.LowerKebab,
	FlagNamer:    naming.LowerKebab,
	EnvVarNamer:  naming.UpperSnake,

	inDecoders: map[string]iocodec.DecoderMaker{
		// "json": iocodec.JSONDecoderMaker(),
		"xml": iocodec.XMLDecoderMaker(),
	},
	outEncoders: map[string]iocodec.EncoderMaker{
		// "json":       iocodec.JSONEncoderMaker(false),
		"prettyjson": iocodec.JSONEncoderMaker(true),
		"xml":        iocodec.XMLEncoderMaker(false),
		"prettyxml":  iocodec.XMLEncoderMaker(true),
	},
}

func NewConfig(options ...Option) *Config {
	c := *DefaultConfig
	for _, opt := range options {
		opt(&c)
	}
	if c.CommandNamer == nil {
		panic("CommandNamer not specified")
	}
	if c.FlagNamer == nil {
		panic("FlagNamer not specified")
	}
	if c.EnvVarNamer == nil {
		panic("EnvVarNamer not specified")
	}
	return &c
}

func RegisterFlagBinder(binder FlagBinder) {
	DefaultConfig.flagBinders = append(DefaultConfig.flagBinders, binder)
}

func RegisterPreDialer(dialer PreDialer) {
	DefaultConfig.preDialers = append(DefaultConfig.preDialers, dialer)
}

func RegisterInputDecoder(format string, maker iocodec.DecoderMaker) {
	DefaultConfig.inDecoders[format] = maker
}

func RegisterOutputEncoder(format string, maker iocodec.EncoderMaker) {
	DefaultConfig.outEncoders[format] = maker
}

func (c *Config) BindFlags(fs *pflag.FlagSet) {
	for _, binder := range c.flagBinders {
		binder(fs, c.FlagNamer)
	}
}

func (c *Config) decoderFormats() []string {
	f := make([]string, len(c.inDecoders))
	i := 0
	for k := range c.inDecoders {
		f[i] = k
		i++
	}
	sort.Strings(f)
	return f
}

func (c *Config) encoderFormats() []string {
	f := make([]string, len(c.outEncoders))
	i := 0
	for k := range c.outEncoders {
		f[i] = k
		i++
	}
	sort.Strings(f)
	return f
}

func RoundTrip(ctx context.Context, cfg *Config, fn func(grpc.ClientConnInterface, iocodec.Decoder, iocodec.Encoder) error) error {
	var err error
	var in iocodec.Decoder
	if cfg.RequestFile == "" && cfg.PreDecoder != nil {
		in = cfg.PreDecoder(ctx)
	} else if in, err = cfg.makeDecoder(); err != nil {
		return err
	}

	var out iocodec.Encoder
	if out, err = cfg.makeEncoder(); err != nil {
		return err
	}

	var cc *grpc.ClientConn
	if cfg.ClientConnFunc != nil {
		cc, err = cfg.ClientConnFunc()
		if err != nil {
			return err
		}
	} else {
		opts := []grpc.DialOption{}
		if err := cfg.dialOpts(ctx, &opts); err != nil {
			return err
		}

		cc, err = grpc.NewClient(cfg.ServerAddr, opts...)
		if err != nil {
			return err
		}
	}
	defer cc.Close()

	return fn(cc, in, out)
}

func (c *Config) makeDecoder() (iocodec.Decoder, error) {
	if c.RequestFile == "" {
		return iocodec.NoOp, nil
	}
	if c.RequestFile != "-" {
		f, err := os.Open(c.RequestFile)
		if err != nil {
			return nil, fmt.Errorf("request file: %v", err)
		}
		var m iocodec.DecoderMaker
		if ext := strings.TrimLeft(filepath.Ext(c.RequestFile), "."); ext != "" {
			m = c.inDecoders[ext]
		}
		if m == nil {
			var ok bool
			if m, ok = c.inDecoders[c.RequestFormat]; !ok {
				return nil, fmt.Errorf("unknown request format: %s", c.RequestFormat)
			}
		}
		return func(v interface{}) error {
			defer f.Close()
			return m(f)(v)
		}, nil
	}

	if c.RequestFormat == "" {
		return iocodec.NoOp, nil
	}
	if m, ok := c.inDecoders[c.RequestFormat]; !ok {
		return nil, fmt.Errorf("unknown request format: %s", c.RequestFormat)
	} else {
		return m(os.Stdin), nil
	}
}

func (c *Config) makeEncoder() (iocodec.Encoder, error) {
	if c.ResponseFormat == "" {
		return iocodec.NoOp, nil
	}
	if m, ok := c.outEncoders[c.ResponseFormat]; !ok {
		return nil, fmt.Errorf("unknown response format: %s", c.ResponseFormat)
	} else {
		return m(os.Stdout), nil
	}
}

func (c *Config) dialOpts(ctx context.Context, opts *[]grpc.DialOption) error {
	if c.TLS {
		tlsConfig := &tls.Config{InsecureSkipVerify: c.InsecureSkipVerify}
		if c.CACertFile != "" {
			caCert, err := os.ReadFile(c.CACertFile)
			if err != nil {
				return fmt.Errorf("ca cert: %v", err)
			}
			certPool := x509.NewCertPool()
			certPool.AppendCertsFromPEM(caCert)
			tlsConfig.RootCAs = certPool
		}
		if c.CertFile != "" {
			if c.KeyFile == "" {
				return fmt.Errorf("key file not specified")
			}
			pair, err := tls.LoadX509KeyPair(c.CertFile, c.KeyFile)
			if err != nil {
				return fmt.Errorf("cert/key: %v", err)
			}
			tlsConfig.Certificates = []tls.Certificate{pair}
		}
		if c.ServerName != "" {
			tlsConfig.ServerName = c.ServerName
		} else {
			addr, _, _ := net.SplitHostPort(c.ServerAddr)
			tlsConfig.ServerName = addr
		}
		cred := credentials.NewTLS(tlsConfig)
		*opts = append(*opts, grpc.WithTransportCredentials(cred))
	} else {
		*opts = append(*opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	for _, dialer := range c.preDialers {
		if err := dialer(ctx, opts); err != nil {
			return err
		}
	}

	return nil
}
