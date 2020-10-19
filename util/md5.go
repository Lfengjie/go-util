package util

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
)

// calc http body md5, the size is http header content-length, not be the real size
func MD5(file io.Reader, length int64) (reader io.Reader, md5Str string, err error) {
	md5Cal := md5.New()
	const BatchSize int64 = 32 * 1024
	buf := make([]byte, BatchSize)
	var size int64

	Min := func(x, y int64) int64 {
		if x < y {
			return x
		}
		return y
	}

	for size < length {
		readSize := Min(length-size, BatchSize)
		nr, er := file.Read(buf[:readSize])
		if nr > 0 {
			nw, ew := md5Cal.Write(buf[:nr])
			if nw > 0 {
				size += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er == io.EOF {
			break
		}
		if er != nil {
			err = er
			break
		}
	}
	if err != nil {
		return nil, "", err
	}
	if length > size {
		return nil, "", errors.New("IncompleteBody")
	}

	md5Str = hex.EncodeToString(md5Cal.Sum(nil))

	// seek fd to header
	file.(io.ReadSeeker).Seek(0, io.SeekStart)
	return file, md5Str, nil
}

// calc string md5
func Md5Str(object string) string {
	sum := md5.Sum([]byte(object))
	return hex.EncodeToString(sum[:])
}

// calc the real size http body md5.
func MD5Hex(file io.Reader) (size int64, md5Str string, err error) {
	md5Cal := md5.New()
	buf := make([]byte, 32*1024)
	for {
		nr, er := file.Read(buf)
		if nr > 0 {
			nw, ew := md5Cal.Write(buf[:nr])
			if nw > 0 {
				size += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er == io.EOF {
			break
		}
		if er != nil {
			err = er
			break
		}
	}

	md5Str = hex.EncodeToString(md5Cal.Sum(nil))

	// seek fd to header
	file.(io.ReadSeeker).Seek(0, io.SeekStart)

	return
}
