// splice.go
package proxy

import "io"

func Splice(r io.Reader, w io.WriteCloser, exit chan<- bool) {
	buf := make([]byte, 1024)
	for {
		rn, err := r.Read(buf)
		if rn > 0 {
			_, err := w.Write(buf[:rn])
			if err != nil {
				w.Close()
				break
			}
		}
		if err != nil {
			w.Close()
			break
		}
	}
	exit <- true
}
