package handler

import "io"

type CacheReader struct {
	reader io.Reader
	cache  []byte
}

func NewCacheReader(r io.Reader) *CacheReader {
	return &CacheReader{
		reader: r,
	}
}

func (c *CacheReader) Read(p []byte) (n int, err error) {
	n, err = c.reader.Read(p)
	if n > 0 {
		c.cache = append(c.cache, p[0:n]...)
	}
	return n, err
}

func (c *CacheReader) Cache() []byte {
	return c.cache
}
