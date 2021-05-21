package conf

type Download struct {
	Host string
	Path string
}

type Proxy struct {
	Enable bool
	Socket string
}

type Config struct {
	Addr     string
	Dbfile   string
	Download *Download
	Proxy    *Proxy
}
