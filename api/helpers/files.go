package helpers

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"sync"
)

var bufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func writeCodeToTarWriter(tw *tar.Writer, code, extension string) error {
	header := &tar.Header{
		Name: "code" + extension,
		Mode: 0640,
		Size: int64(len(code)),
	}
	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	_, err := tw.Write([]byte(code))
	return err
}

func CreateTarArchiveForSourceCode(code, extension string) (io.Reader, error) {
	buf, ok := bufPool.Get().(*bytes.Buffer)
	if !ok {
		return nil, fmt.Errorf("error reflecting buffer pool to buffer")
	}

	buf.Reset()
	defer bufPool.Put(buf)

	tarWriter := tar.NewWriter(buf)
	defer tarWriter.Close()

	if err := writeCodeToTarWriter(tarWriter, code, extension); err != nil {
		return nil, err
	}

	if err := tarWriter.Close(); err != nil {
		return nil, err
	}

	return bytes.NewReader(buf.Bytes()), nil
}
