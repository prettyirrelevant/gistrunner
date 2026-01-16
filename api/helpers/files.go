package helpers

import (
	"archive/tar"
	"bytes"
	"io"
)

func writeCodeToTarWriter(tw *tar.Writer, code, extension string) error {
	header := &tar.Header{
		Name: "code" + extension,
		Mode: 0o640,
		Size: int64(len(code)),
	}
	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	_, err := tw.Write([]byte(code))
	return err
}

func CreateTarArchiveForSourceCode(code, extension string) (io.Reader, error) {
	buf := new(bytes.Buffer)
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
